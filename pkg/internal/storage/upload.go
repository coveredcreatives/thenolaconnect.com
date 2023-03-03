package internal

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"time"

	"cloud.google.com/go/storage"
	alog "github.com/apex/log"
	"github.com/coveredcreatives/thenolaconnect.com/pkg/devtools"
	"github.com/coveredcreatives/thenolaconnect.com/pkg/model"
	"gorm.io/gorm"
)

func UploadImageToStorage(storagec *storage.Client, file image.Image, filename string, bucketname string) (*storage.ObjectAttrs, error) {
	storageobj := storagec.Bucket(bucketname).Object(filename)
	storageobjw := storageobj.NewWriter(context.Background())
	err := png.Encode(storageobjw, file)
	if err != nil {
		return nil, err
	}
	if err := storageobjw.Close(); err != nil {
		alog.WithError(err).Error("failed to close object writer")
		return nil, err
	}
	// make publicly accessible
	acl := storageobj.ACL()
	if err := acl.Set(context.Background(), storage.AllUsers, storage.RoleReader); err != nil {
		alog.WithError(err).Error("failed to set access control level to public")
		return nil, err
	}
	return storageobj.Attrs(context.Background())
}

func UploadFileToStorage(storagec *storage.Client, file multipart.File, header *multipart.FileHeader, bucketname string) (*storage.ObjectAttrs, error) {
	storageobj := storagec.Bucket(bucketname).Object(header.Filename)
	storageobjw := storageobj.NewWriter(context.Background())
	fileb := make([]byte, int(header.Size))
	_, err := file.Read(fileb)
	if err != nil {
		alog.WithError(err).Error("failed to read bytes from file")
		return nil, err
	}
	_, err = fmt.Fprint(storageobjw, string(fileb))
	if err != nil {
		alog.WithError(err).Error("failed to write bytes to storage object")
		return nil, err
	}
	if err := storageobjw.Close(); err != nil {
		alog.WithError(err).Error("failed to close object writer")
		return nil, err
	}
	// Make publicly accessible
	acl := storageobj.ACL()
	if err := acl.Set(context.Background(), storage.AllUsers, storage.RoleReader); err != nil {
		alog.WithError(err).Error("failed to set access control level to public")
		return nil, err
	}
	return storageobj.Attrs(context.Background())
}

func StoreCloudStorageObject(db *gorm.DB, storagec *storage.Client, pgw model.PDFGenerationWorker, order model.Order, b []byte, google_env_config devtools.GoogleApplicationConfig, start time.Time) (err error) {
	alog.Trace("StoreCloudStorageObject").Stop(&err)
	object_name := fmt.Sprint("order_communications/", start.Year(), "/", start.Month(), "/", start.Day(), "/", order.MediaFilename())

	storage_object := storagec.
		Bucket(google_env_config.EnvGoogleStorageBucketName).
		Object(object_name)
	wc := storage_object.
		NewWriter(context.Background())
	wc.ChunkSize = 2 * (1024 * 1024) // 2 MiB
	bytes_written, err := io.Copy(wc, bytes.NewReader(b))
	if err != io.EOF && err != nil {
		return
	}

	alog.
		WithField("object_name", object_name).
		WithField("bytes_written", bytes_written).
		Info("stream object to cloud storage")

	err = wc.Close()
	if err != nil {
		return
	}

	alog.WithField("object_name", object_name).WithField("bucket", google_env_config.EnvGoogleStorageBucketName).Info("uploaded")

	acl := storage_object.ACL()
	err = acl.Set(context.Background(), storage.AllUsers, storage.RoleReader)
	if err != nil {
		return
	}
	completed_at := time.Now()

	public_url := fmt.Sprint("https://storage.googleapis.com/", google_env_config.EnvGoogleStorageBucketName, "/", object_name)
	alog.WithField("public_url", public_url).Info("located")

	tx := db.
		Model(&model.PDFGenerationWorker{}).
		Where(&model.PDFGenerationWorker{PDFGenerationWorkerId: pgw.PDFGenerationWorkerId}).
		Update("completed_at", completed_at)

	err = tx.Error
	if err != nil {
		return
	}
	alog.WithField("pdf_generation_worker_id", pgw.PDFGenerationWorkerId).
		WithField("completed_at", completed_at).Info("update pdf_generation_worker completion timestamp")

	tx = db.
		Model(&model.Order{}).
		Where(&model.Order{PDFGenerationWorkerId: pgw.PDFGenerationWorkerId}).
		Updates(&model.Order{
			OrderFileStorageURL: public_url,
			IsPDFGenerated:      true,
		})
	alog.WithField("rows_affected", tx.RowsAffected).Info("update order with file storage url")

	err = tx.Error
	if err != nil {
		return
	}
	return
}
