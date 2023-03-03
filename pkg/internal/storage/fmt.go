package internal

import (
	"fmt"

	"cloud.google.com/go/storage"
)

func FmtPublicUrl(objAttrs *storage.ObjectAttrs) string {
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", objAttrs.Bucket, objAttrs.Name)
}
