// internal/aws/s3/service.go
package s3

import (
	"context"
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

// ListBuckets はS3バケットのリストを取得し、ドメインオブジェクトに変換します。
func (s *BucketService) ListBuckets(ctx context.Context, resourceName string) ([]Bucket, error) {
	awsBuckets, err := s.repo.ListBuckets(ctx)
	if err != nil {
		return nil, err
	}

	var buckets []Bucket
	for _, b := range awsBuckets {
		if resourceName == "" || (b.Name != nil && *b.Name == resourceName) {
			tags, err := s.repo.GetBucketTagging(ctx, *b.Name)
			if err != nil {
				// handle error, maybe log it
				continue
			}
			bucket := Bucket{
				Name: *b.Name,
				Tags: tags,
			}
			buckets = append(buckets, bucket)
		}
	}
	return buckets, nil
} 
