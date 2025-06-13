// internal/aws/s3/repository.go
package s3

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// S3RepositoryInterface はS3リソースへのアクセスを抽象化します。
type S3RepositoryInterface interface {
	ListBuckets(ctx context.Context) ([]types.Bucket, error)
}

// S3Repository はS3RepositoryInterfaceを実装します。
type S3Repository struct {
	client *s3.Client
}

// NewS3Repository は新しいS3Repositoryを生成します。
func NewS3Repository(client *s3.Client) *S3Repository {
	return &S3Repository{
		client: client,
	}
}

// ListBuckets はAWSからS3バケットのリストを取得します。
func (r *S3Repository) ListBuckets(ctx context.Context) ([]types.Bucket, error) {
	input := &s3.ListBucketsInput{}
	result, err := r.client.ListBuckets(ctx, input)
	if err != nil {
		return nil, err
	}
	return result.Buckets, nil
} 
