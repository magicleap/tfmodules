package gcs

import (
	"fmt"
	"io"

	// "path"
	// "time"
	"context"

	"cloud.google.com/go/storage"
)

var (
	contentType = "application/x-gzip"
)

// Create a binray file from upload tar.gz
// Also renames the file to match convention and record the metadata version
func (b *Bucket) Write(namespace, name, system, version string, source string, body io.Reader) error {

	path := getFilePath(namespace, name, system, version)

	// timeout in case of blocking
	ctx := context.Background()
	o := b.Handle.Object(path)
	writer := o.NewWriter(ctx)

	if _, err := io.Copy(writer, body); err != nil {
		return fmt.Errorf("unable to write data to bucket : %w", err)
	}

	// the defer writer.Close prevents the update of the metadata.
	// it triggers an object not found if not closed
	// So we need to close the writer before updating the metadata
	if err := writer.Close(); err != nil {
		return fmt.Errorf("unable to close object writer : %w", err)
	}

	// we store version and source as metadata in GCS object
	attrs := storage.ObjectAttrsToUpdate{
		ContentType: contentType,
		Metadata: map[string]string{
			moduleVersionMetadata: version,
			moduleSourceMetadata:  source,
		},
	}

	if _, err := o.Update(context.Background(), attrs); err != nil {
		return fmt.Errorf("unable to update metadata : %w", err)
	}

	return nil
}
