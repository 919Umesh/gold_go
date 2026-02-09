package appwrite

import (
	"github.com/appwrite/sdk-for-go/databases"
)

func WithListDocumentsQueries(queries []string) databases.ListDocumentsOption {
	return func(o *databases.ListDocumentsOptions) {
		o.Queries = queries
	}
}


func WithUpdateDocumentData(data interface{}) databases.UpdateDocumentOption {
	return func(o *databases.UpdateDocumentOptions) {
		o.Data = data
	}
}
