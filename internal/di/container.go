// internal/di/container.go
package di

import (
	"context"
	"fmt"

	"github.com/Haussmann000/tfimport/internal/aws"
	"github.com/Haussmann000/tfimport/internal/aws/ec2"
	"github.com/Haussmann000/tfimport/internal/aws/s3"
	"github.com/Haussmann000/tfimport/internal/hcl"
	"github.com/Haussmann000/tfimport/internal/writer"
)

// RunOptions はコマンドラインから渡されるオプションを保持します。
type RunOptions struct {
	ResourceTypes []string
	ResourceName  string
}

// App はアプリケーションの主要なロジックをカプセル化します。
type App struct {
	ec2Svc       ec2.EC2ServiceInterface
	s3Svc        s3.S3ServiceInterface
	hclGenerator *hcl.HCLGenerator
	fileWriter   *writer.FileWriter
}

// NewApp は新しいAppを生成します。
func NewApp(
	ec2Svc ec2.EC2ServiceInterface,
	s3Svc s3.S3ServiceInterface,
	hclGenerator *hcl.HCLGenerator,
	fileWriter *writer.FileWriter,
) *App {
	return &App{
		ec2Svc:       ec2Svc,
		s3Svc:        s3Svc,
		hclGenerator: hclGenerator,
		fileWriter:   fileWriter,
	}
}

// Run はアプリケーションのメインの処理を実行します。
func (a *App) Run(ctx context.Context, options RunOptions) error {
	for _, resourceType := range options.ResourceTypes {
		switch resourceType {
		case "vpc":
			if err := a.processVpcs(ctx, options.ResourceName); err != nil {
				return err
			}
		case "s3":
			if err := a.processS3Buckets(ctx, options.ResourceName); err != nil {
				return err
			}
		default:
			fmt.Printf("Unsupported resource type: %s\n", resourceType)
		}
	}

	fmt.Println("Terraform files generated successfully.")
	return nil
}

func (a *App) processVpcs(ctx context.Context, resourceName string) error {
	// VPC情報を取得
	vpcs, err := a.ec2Svc.ListVpcs(ctx, resourceName)
	if err != nil {
		return fmt.Errorf("failed to list vpcs: %w", err)
	}
	if len(vpcs) == 0 {
		fmt.Println("No VPCs found.")
		return nil
	}

	// HCLを生成
	resourceFile, importFile, err := a.hclGenerator.GenerateVpcBlocks(vpcs)
	if err != nil {
		return fmt.Errorf("failed to generate hcl for vpc: %w", err)
	}

	// ファイルに書き出し
	err = a.fileWriter.WriteFile("vpc_generated.tf", resourceFile)
	if err != nil {
		return fmt.Errorf("failed to write resource file for vpc: %w", err)
	}
	err = a.fileWriter.WriteFile("vpc_import.tf", importFile)
	if err != nil {
		return fmt.Errorf("failed to write import file for vpc: %w", err)
	}
	return nil
}

func (a *App) processS3Buckets(ctx context.Context, resourceName string) error {
	// S3バケット情報を取得
	buckets, err := a.s3Svc.ListBuckets(ctx, resourceName)
	if err != nil {
		return fmt.Errorf("failed to list s3 buckets: %w", err)
	}
	if len(buckets) == 0 {
		fmt.Println("No S3 buckets found.")
		return nil
	}

	// HCLを生成
	resourceFile, importFile, err := a.hclGenerator.GenerateS3BucketBlocks(buckets)
	if err != nil {
		return fmt.Errorf("failed to generate hcl for s3: %w", err)
	}

	// ファイルに書き出し
	err = a.fileWriter.WriteFile("s3_generated.tf", resourceFile)
	if err != nil {
		return fmt.Errorf("failed to write resource file for s3: %w", err)
	}
	err = a.fileWriter.WriteFile("s3_import.tf", importFile)
	if err != nil {
		return fmt.Errorf("failed to write import file for s3: %w", err)
	}
	return nil
}

// BuildApp は依存関係を解決してAppを構築します。
func BuildApp(ctx context.Context) (*App, error) {
	// AWS Config
	awsCfg, err := aws.NewConfig(ctx)
	if err != nil {
		return nil, err
	}

	// Clients
	ec2Client := aws.NewEC2Client(awsCfg)
	s3Client := aws.NewS3Client(awsCfg)

	// Repositories
	ec2Repo := ec2.NewEC2Repository(ec2Client)
	s3Repo := s3.NewS3Repository(s3Client)

	// Services
	ec2Svc := ec2.NewEC2Service(ec2Repo)
	s3Svc := s3.NewS3Service(s3Repo)

	// HCL Generator
	hclGenerator := hcl.NewHCLGenerator()

	// File Writer
	fileWriter := writer.NewFileWriter()

	// App
	app := NewApp(ec2Svc, s3Svc, hclGenerator, fileWriter)

	return app, nil
} 
