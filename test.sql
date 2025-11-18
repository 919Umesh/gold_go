CREATE OR REPLACE FUNCTION generate_payout_by_brand(
  p_order_data jsonb,
  p_updated_by uuid
)
RETURNS jsonb
LANGUAGE plpgsql
SECURITY DEFINER
AS $$
DECLARE
  brand_record jsonb;
  inserted_payouts jsonb := '[]'::jsonb;
  new_payout_record jsonb;
  aggregated_data jsonb;
  existing_payout_count integer;
  product_ids_array uuid[];
BEGIN
  -- Check if input data is not empty
  IF p_order_data IS NULL OR jsonb_array_length(p_order_data) = 0 THEN
    RAISE NOTICE 'No order data provided';
    RETURN '[]'::jsonb;
  END IF;

  -- Aggregate orders by brand_id
  SELECT jsonb_agg(
    jsonb_build_object(
      'brand_id', brand_id,
      'invoice_number', invoice_number,
      'total_amount', total_amount,
      'platform_fees', total_platform_fees, 
      'total_payout', total_payout,
      'ordered_product_ids', product_ids,
      'order_date', order_date,
      'brand_name', brand_name
    )
  )
  INTO aggregated_data
  FROM (
    SELECT 
      brand_id,
      MIN(order_number) as invoice_number,
      SUM(amount) as total_amount,
      SUM(platform_fees) as total_platform_fees,
      SUM(total_payout) as total_payout,
      jsonb_agg(ordered_product_id) as product_ids,
      MIN(order_date) as order_date,
      MIN(brand_name) as brand_name
    FROM jsonb_to_recordset(p_order_data) AS x(
      brand_id uuid,
      order_number text,
      amount numeric,
      platform_fees numeric,
      total_payout numeric,
      ordered_product_id uuid,
      order_date timestamptz,
      brand_name text
    )
    GROUP BY brand_id
  ) AS aggregated_data;

  -- If no aggregated data found, return empty array
  IF aggregated_data IS NULL THEN
    RETURN '[]'::jsonb;
  END IF;

  -- Loop through each aggregated brand record and insert into payouts
  FOR brand_record IN SELECT * FROM jsonb_array_elements(aggregated_data)
  LOOP
    -- Check if payout already exists for this brand and invoice number
    SELECT COUNT(*) INTO existing_payout_count
    FROM payouts 
    WHERE brand_id = (brand_record->>'brand_id')::uuid
      AND invoice_number = brand_record->>'invoice_number';

    IF existing_payout_count = 0 THEN
      -- Convert JSON array to PostgreSQL UUID array
      SELECT ARRAY(
        SELECT jsonb_array_elements_text(brand_record->'ordered_product_ids')::uuid
      ) INTO product_ids_array;

      -- Insert new payout record
      INSERT INTO payouts (
        created_at,
        payout_date,
        updated_by,
        brand_id,
        total_amount,
        invoice_number,
        platform_fees,
        total_payout,
        ordered_product_ids,
        status
      ) VALUES (
        NOW(),
        (brand_record->>'order_date')::timestamptz, 
        p_updated_by,
        (brand_record->>'brand_id')::uuid, 
        (brand_record->>'total_amount')::numeric, 
        brand_record->>'invoice_number', 
        (brand_record->>'platform_fees')::numeric, 
        (brand_record->>'total_payout')::numeric, 
        product_ids_array,  -- Use the converted array
        'PENDING'
      )
      RETURNING jsonb_build_object(
        'created_at', created_at,
        'payout_date', payout_date,
        'updated_by', updated_by,
        'brand_id', brand_id,
        'total_amount', total_amount,
        'invoice_number', invoice_number,
        'platform_fees', platform_fees,
        'total_payout', total_payout,
        'ordered_product_ids', ordered_product_ids,
        'status', status
      ) INTO new_payout_record;

      inserted_payouts := inserted_payouts || jsonb_build_array(new_payout_record);
    ELSE
      RAISE NOTICE 'Payout already exists for brand % and invoice %', 
        brand_record->>'brand_id', 
        brand_record->>'invoice_number';
    END IF;
  END LOOP;

  RETURN inserted_payouts;
END;
$$;