// internal/aws/client.go
package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// NewConfig はAWSの設定をロードして返します。
func NewConfig(ctx context.Context) (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return aws.Config{}, fmt.Errorf("failed to load aws config: %w", err)
	}
	return cfg, nil
}

// NewEC2Client はEC2サービスクライアントを生成します。
func NewEC2Client(cfg aws.Config) *ec2.Client {
	return ec2.NewFromConfig(cfg)
}

// NewS3Client はS3サービスクライアントを生成します。
func NewS3Client(cfg aws.Config) *s3.Client {
	return s3.NewFromConfig(cfg)
}

// NewS3Client ... (今後他のクライアントもここに追加) 
