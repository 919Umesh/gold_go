package appwrite

import (
	"github.com/appwrite/sdk-for-go/databases"
)

func WithListDocumentsQueries(queries []string) databases.ListDocumentsOption {
	return func(o *databases.ListDocumentsOptions) {
		o.Queries = queries
	}
}

func WithCreateCollectionPermissions(permissions []string) databases.CreateCollectionOption {
	return func(o *databases.CreateCollectionOptions) {
		o.Permissions = permissions
	}
}

func WithCreateCollectionDocumentSecurity(enabled bool) databases.CreateCollectionOption {
	return func(o *databases.CreateCollectionOptions) {
		o.DocumentSecurity = enabled
	}
}
