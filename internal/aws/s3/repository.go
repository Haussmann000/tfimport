// internal/aws/s3/repository.go
package s3

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3ClientInterface interface {
	ListBuckets(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error)
	GetBucketTagging(ctx context.Context, params *s3.GetBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.GetBucketTaggingOutput, error)
}

type S3RepositoryInterface interface {
	ListBuckets(ctx context.Context) ([]types.Bucket, error)
	GetBucketTagging(ctx context.Context, bucketName string) (map[string]string, error)
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

// GetBucketTagging は指定されたバケットのタグを取得します。
func (r *S3Repository) GetBucketTagging(ctx context.Context, bucketName string) (map[string]string, error) {
	output, err := r.client.GetBucketTagging(ctx, &s3.GetBucketTaggingInput{
		Bucket: &bucketName,
	})
	if err != nil {
		return nil, err
	}
	tags := make(map[string]string)
	for _, tag := range output.TagSet {
		tags[*tag.Key] = *tag.Value
	}
	return tags, nil
} 
