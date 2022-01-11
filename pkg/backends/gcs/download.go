package gcs

import (
	"io"
	"fmt"
	"context"
)

// Create a binray file from upload tar.gz
// Also renames the file to match convention and record the metadata version
func (b *Bucket) Read(namespace, name, system, version string, response io.Writer) error {

	path := getFilePath(namespace, name, system, version)
	ctx := context.Background()

	reader, err := b.Handle.Object(path).NewReader(ctx)
	if err != nil {
		return fmt.Errorf("unable to open readser on object : %w", err)
	}
	defer reader.Close()

	if _, err := io.Copy(response, reader); err != nil {
		return fmt.Errorf("unable to read data from bucket : %w", err)
	}

	return nil

}
