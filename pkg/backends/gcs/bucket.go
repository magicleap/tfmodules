package gcs

import (
	"context"
	"fmt"
	"os"
	"path"

	"cloud.google.com/go/storage"
	"github.com/sirupsen/logrus"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

// implement BackendInterface
type Bucket struct {
	Client *storage.Client
	Handle *storage.BucketHandle
}

const (
	moduleVersionMetadata = "x-module-version"
	moduleSourceMetadata  = "x-module-source"
)

// Bucket creates a client to an existing bucket
// It may need GOOGLE_APPLICATION_CREDENTIALS envar if workload identity is not set
func New() (*Bucket, error) {

	logrus.Debug("Creating bucketHandler")
	var err error
	bucketName := os.Getenv("GOOGLE_BUCKET")
	bucket := &Bucket{}
	ctx := context.Background()

	bucket.Client, err = storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not create bucket client : %w", err)
	}

	// Get bucket
	bucket.Handle = bucket.Client.Bucket(bucketName)
	// check that we access the bucket
	_, err = bucket.Handle.Attrs(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting Attrs : %w", err)
	}

	return bucket, nil
}

func (b *Bucket) Close() error {
	return b.Client.Close()
}

func getFilePath(namespace, name, system, version string) string {
	return path.Join(namespace, name, system, fmt.Sprintf("%s-%s-%s.tar.gz", name, system, version))
}
