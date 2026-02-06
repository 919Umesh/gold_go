package appwrite

import (
	"github.com/appwrite/sdk-for-go/databases"
)

// WithListDocumentsQueries creates a ListDocumentsOption with the given queries.
func WithListDocumentsQueries(queries []string) databases.ListDocumentsOption {
	return func(o *databases.ListDocumentsOptions) {
		o.Queries = queries
	}
}

// WithUpdateDocumentData creates an UpdateDocumentOption with the given data.
func WithUpdateDocumentData(data interface{}) databases.UpdateDocumentOption {
	return func(o *databases.UpdateDocumentOptions) {
		o.Data = data
	}
}
