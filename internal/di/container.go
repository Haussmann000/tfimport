// internal/di/container.go
package di

import (
	"context"
	"fmt"
	"strings"

	"github.com/Haussmann000/tfimport/internal/aws/ec2"
	"github.com/Haussmann000/tfimport/internal/aws/ecs"
	"github.com/Haussmann000/tfimport/internal/aws/elbv2"
	"github.com/Haussmann000/tfimport/internal/aws/iam"
	"github.com/Haussmann000/tfimport/internal/aws/s3"
	"github.com/Haussmann000/tfimport/internal/hcl"
	"github.com/Haussmann000/tfimport/internal/writer"
	"golang.org/x/sync/errgroup"
)

// RunOptions はコマンドラインから渡されるオプションを保持します。
type RunOptions struct {
	ResourceTypes []string
	ResourceName  string
	ClusterName   string
	ServiceName   string
}

// App はアプリケーションの主要なロジックをカプセル化します。
type App struct {
	s3Service  *s3.BucketService
	vpcService *ec2.VPCService
	ecsService *ecs.ECSService
	elbService *elbv2.ELBV2Service
	iamService *iam.IAMService
	writer     *writer.FileWriter
	generator  *hcl.HCLGenerator
}

// NewApp はAppのコンストラクタです。
func NewApp(
	s3s *s3.BucketService,
	vs *ec2.VPCService,
	ecss *ecs.ECSService,
	elbs *elbv2.ELBV2Service,
	iams *iam.IAMService,
	w *writer.FileWriter,
	g *hcl.HCLGenerator,
) *App {
	return &App{
		s3Service:  s3s,
		vpcService: vs,
		ecsService: ecss,
		elbService: elbs,
		iamService: iams,
		writer:     w,
		generator:  g,
	}
}

// Run はアプリケーションのメインの処理を実行します。
func (a *App) Run(ctx context.Context, options RunOptions) error {
	for _, resourceType := range options.ResourceTypes {
		switch resourceType {
		case "vpc":
			if err := a.processVpc(ctx, options.ResourceName); err != nil {
				return err
			}
		case "s3":
			if err := a.processS3(ctx, options.ResourceName); err != nil {
				return err
			}
		case "ecs":
			if err := a.processEcs(ctx, options.ClusterName, options.ServiceName); err != nil {
				return err
			}
		case "elbv2":
			if err := a.processElb(ctx, options.ResourceName); err != nil {
				return err
			}
		case "iam":
			if err := a.processIam(ctx, options.ResourceName); err != nil {
				return err
			}
		default:
			fmt.Printf("Unsupported resource type: %s\n", resourceType)
		}
	}

	fmt.Println("Terraform files generated successfully.")
	return nil
}

func (a *App) processS3(ctx context.Context, resourceName string) error {
	buckets, err := a.s3Service.ListBuckets(ctx, resourceName)
	if err != nil {
		return err
	}
	hclFile, importFile, err := a.generator.GenerateS3BucketBlocks(buckets)
	if err != nil {
		return err
	}
	err = a.writer.WriteFile("s3_generated.tf", hclFile)
	if err != nil {
		return err
	}
	return a.writer.WriteFile("s3_import.tf", importFile)
}

func (a *App) processVpc(ctx context.Context, resourceName string) error {
	vpcs, err := a.vpcService.ListVpcs(ctx, resourceName)
	if err != nil {
		return err
	}
	hclFile, importFile, err := a.generator.GenerateVpcBlocks(vpcs)
	if err != nil {
		return err
	}
	err = a.writer.WriteFile("vpc_generated.tf", hclFile)
	if err != nil {
		return err
	}
	return a.writer.WriteFile("vpc_import.tf", importFile)
}

func (a *App) processEcs(ctx context.Context, clusterName, serviceName string) error {
	var clusters []ecs.Cluster
	if clusterName != "" {
		cluster, err := a.ecsService.GetClusters(ctx, clusterName, serviceName)
		if err != nil {
			return err
		}
		if cluster != nil {
			clusters = cluster
		}
	}

	hclFile, importFile, err := a.generator.GenerateEcsBlocks(clusters)
	if err != nil {
		return err
	}
	err = a.writer.WriteFile("ecs_generated.tf", hclFile)
	if err != nil {
		return err
	}
	return a.writer.WriteFile("ecs_import.tf", importFile)
}

func (a *App) processElb(ctx context.Context, resourceName string) error {
	var lbs []*elbv2.LoadBalancer
	var err error
	if resourceName == "" {
		lbs, err = a.elbService.ListLoadBalancers(ctx)
	} else {
		var lb *elbv2.LoadBalancer
		lb, err = a.elbService.GetLoadBalancer(ctx, resourceName)
		if lb != nil {
			lbs = append(lbs, lb)
		}
	}
	if err != nil {
		return err
	}
	hclFile, importFile, err := a.generator.GenerateElbBlocks(lbs)
	if err != nil {
		return err
	}
	err = a.writer.WriteFile("elb_generated.tf", hclFile)
	if err != nil {
		return err
	}
	return a.writer.WriteFile("elb_import.tf", importFile)
}

func (a *App) processIam(ctx context.Context, nameContains string) error {
	var policies []iam.Policy
	var roles []iam.Role
	var eg errgroup.Group

	eg.Go(func() error {
		var err error
		policies, err = a.iamService.ListPolicies(ctx, nameContains)
		return err
	})

	eg.Go(func() error {
		var err error
		roles, err = a.iamService.ListRoles(ctx, nameContains)
		return err
	})

	if err := eg.Wait(); err != nil {
		return err
	}

	hclFile, importFile, err := a.generator.GenerateIamBlocks(policies, roles)
	if err != nil {
		return err
	}

	err = a.writer.WriteFile("iam_generated.tf", hclFile)
	if err != nil {
		return err
	}
	return a.writer.WriteFile("iam_import.tf", importFile)
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
	ecsClient := aws.NewECSClient(awsCfg)
	elbv2Client := aws.NewELBV2Client(awsCfg)

	// Repositories
	ec2Repo := ec2.NewEC2Repository(ec2Client)
	s3Repo := s3.NewS3Repository(s3Client)
	ecsRepo := ecs.NewECSRepository(ecsClient)
	elbv2Repo := elbv2.NewELBV2Repository(elbv2Client)

	// Services
	ec2Svc := ec2.NewEC2Service(ec2Repo)
	s3Svc := s3.NewS3Service(s3Repo)
	ecsSvc := ecs.NewECSService(ecsRepo)
	elbv2Svc := elbv2.NewELBV2Service(elbv2Repo)

	// HCL Generator
	hclGenerator := hcl.NewHCLGenerator()

	// File Writer
	fileWriter := writer.NewFileWriter()

	// App
	app := NewApp(s3Svc, ec2Svc, ecsSvc, elbv2Svc, hclGenerator, fileWriter)

	return app, nil
} 
