package gcs

import (
	"context"
	"fmt"
	"path"

	"github.com/magicleap/tfmodules/pkg/backends/types"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

// ListMetadata returns the available metadata for a given module
func (b *Bucket) ListMetadata(namespace, name, system string) (metadata []types.Metadata, err error) {

	path := path.Join(namespace, name, system)
	query := &storage.Query{Prefix: path}

	// we parse each object of the bucket to retrieve the version metadata
	it := b.Handle.Objects(context.Background(), query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			return metadata, fmt.Errorf("failed listing bucket objects : %w", err)
		}

		// check if the object has a metadata `x-module-version`.
		// No need to parse an object that has no version
		if version, ok := attrs.Metadata[moduleVersionMetadata]; ok {

			m := types.Metadata{Version: version}

			// get source if available
			if source, ok := attrs.Metadata[moduleSourceMetadata]; ok {
				m.Source = source
			}

			metadata = append(metadata, m)
		}
	}

	return metadata, nil
}
