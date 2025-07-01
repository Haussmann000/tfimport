// internal/aws/s3/service.go
package s3

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"golang.org/x/sync/errgroup"
)

// Bucket はHCL生成に必要なS3バケットの情報を保持します。
type Bucket struct {
	Name string
	Tags map[string]string
}

// Service はS3関連のビジネスロジックを定義します。
type Service interface {
	ListBuckets(ctx context.Context, resourceName string) ([]Bucket, error)
}

// BucketService はServiceを実装します。
type BucketService struct {
	repo S3RepositoryInterface
}

// NewBucketService は新しいBucketServiceを生成します。
func NewBucketService(repo S3RepositoryInterface) *BucketService {
	return &BucketService{
		repo: repo,
	}
}

// ListBuckets は指定されたバケット名に一致するS3バケットのリストを取得します。
func (s *BucketService) ListBuckets(ctx context.Context, bucketName string) ([]Bucket, error) {
	awsBuckets, err := s.repo.ListBuckets(ctx)
	if err != nil {
		return nil, err
	}

	var targetBuckets []types.Bucket
	for _, b := range awsBuckets {
		if b.Name != nil && *b.Name == bucketName {
			targetBuckets = append(targetBuckets, b)
			break // Found the bucket, no need to iterate further.
		}
	}

	if len(targetBuckets) == 0 {
		return []Bucket{}, nil
	}

	var eg errgroup.Group
	bucketsWithTags := make(chan Bucket, len(targetBuckets))

	for _, b := range targetBuckets {
		bucket := b
		eg.Go(func() error {
			tags, err := s.repo.GetBucketTagging(ctx, *bucket.Name)
			if err != nil {
				// Tagging might not exist, treat as empty map.
				return nil
			}
			bucketsWithTags <- Bucket{
				Name: *bucket.Name,
				Tags: tags,
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}
	close(bucketsWithTags)

	var result []Bucket
	for b := range bucketsWithTags {
		result = append(result, b)
	}

	return result, nil
} 
