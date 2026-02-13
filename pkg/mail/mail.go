package mail

import "errors"

func SendEmail(email, emailType string, topUpData interface{}) error {
	if email == "" || emailType == "" || topUpData == nil {
		return errors.New("invalid input: email, emailType and topUpData are required")
	}

	return nil
}