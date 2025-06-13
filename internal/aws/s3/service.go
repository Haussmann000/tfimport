// internal/aws/s3/service.go
package s3

import (
	"context"
	"strings"
)

// Bucket はHCL生成に必要なS3バケットの情報を保持します。
type Bucket struct {
	Name string
	// Tags は今後実装
}

// S3ServiceInterface はS3関連のビジネスロジックを定義します。
type S3ServiceInterface interface {
	ListBuckets(ctx context.Context, resourceName string) ([]Bucket, error)
}

// S3Service はS3ServiceInterfaceを実装します。
type S3Service struct {
	repo S3RepositoryInterface
}

// NewS3Service は新しいS3Serviceを生成します。
func NewS3Service(repo S3RepositoryInterface) *S3Service {
	return &S3Service{
		repo: repo,
	}
}

// ListBuckets はS3バケットのリストを取得し、ドメインオブジェクトに変換します。
func (s *S3Service) ListBuckets(ctx context.Context, resourceName string) ([]Bucket, error) {
	awsBuckets, err := s.repo.ListBuckets(ctx)
	if err != nil {
		return nil, err
	}

	var buckets []Bucket
	for _, b := range awsBuckets {
		if resourceName == "" || strings.HasPrefix(*b.Name, resourceName) {
			buckets = append(buckets, Bucket{
				Name: *b.Name,
			})
		}
	}
	return buckets, nil
} 
