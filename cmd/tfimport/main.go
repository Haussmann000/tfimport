// cmd/tfimport/main.go
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/Haussmann000/tfimport/internal/aws"
	"github.com/Haussmann000/tfimport/internal/aws/ec2"
	"github.com/Haussmann000/tfimport/internal/aws/ecs"
	"github.com/Haussmann000/tfimport/internal/aws/elbv2"
	"github.com/Haussmann000/tfimport/internal/aws/iam"
	"github.com/Haussmann000/tfimport/internal/aws/s3"
	"github.com/Haussmann000/tfimport/internal/hcl"
	"github.com/Haussmann000/tfimport/internal/writer"
	"golang.org/x/sync/errgroup"
)

type App struct {
	s3Service  s3.Service
	vpcService ec2.Service
	ecsService ecs.Service
	elbService elbv2.Service
	iamService iam.Service
	writer     *writer.FileWriter
	generator  *hcl.HCLGenerator
}

func NewApp(
	s3s s3.Service,
	vs ec2.Service,
	ecss ecs.Service,
	elbs elbv2.Service,
	iams iam.Service,
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

func (a *App) Run(ctx context.Context, resourceTypes string, resourceName string, clusterName string) error {
	types := strings.Split(resourceTypes, ",")
	for _, resourceType := range types {
		switch resourceType {
		case "s3":
			err := a.processS3(ctx, resourceName)
			if err != nil {
				return err
			}
		case "vpc":
			err := a.processVpc(ctx, resourceName)
			if err != nil {
				return err
			}
		case "ecs":
			err := a.processEcs(ctx, clusterName, resourceName)
			if err != nil {
				return err
			}
		case "elbv2":
			err := a.processElb(ctx, resourceName)
			if err != nil {
				return err
			}
		case "iam":
			err := a.processIam(ctx, resourceName)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported resource type: %s", resourceType)
		}
	}
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
		cls, err := a.ecsService.GetClusters(ctx, clusterName, serviceName)
		if err != nil {
			return err
		}
		clusters = cls
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
	lbs, err := a.elbService.ListLoadBalancers(ctx, resourceName)
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

func main() {
	var resourceTypes, resourceName, clusterName string
	flag.StringVar(&resourceTypes, "resource-types", "", "aws resource type. s3, vpc, ecs, elbv2, iam")
	flag.StringVar(&resourceName, "resource-name", "", "aws resource name")
	flag.StringVar(&clusterName, "cluster-name", "", "ecs cluster name")
	flag.Parse()

	if resourceTypes == "" {
		log.Fatal("resource-types is required")
	}

	ctx := context.Background()
	cfg, err := aws.NewConfig(ctx)
	if err != nil {
		log.Fatalf("failed to load aws config: %v", err)
	}

	// Manual DI
	s3Client := aws.NewS3Client(cfg)
	s3Repo := s3.NewS3Repository(s3Client)
	s3Service := s3.NewBucketService(s3Repo)

	ec2Client := aws.NewEC2Client(cfg)
	ec2Repo := ec2.NewEC2Repository(ec2Client)
	vpcService := ec2.NewVPCService(ec2Repo)

	ecsClient := aws.NewECSClient(cfg)
	ecsRepo := ecs.NewECSRepository(ecsClient)
	ecsService := ecs.NewECSService(ecsRepo)

	elbClient := aws.NewELBV2Client(cfg)
	elbRepo := elbv2.NewELBV2Repository(elbClient)
	elbService := elbv2.NewELBV2Service(elbRepo)

	iamClient := iam.NewIAMClient(cfg)
	iamRepo := iam.NewIAMRepository(iamClient)
	iamService := iam.NewIAMService(iamRepo)

	hclGenerator := hcl.NewHCLGenerator()
	fileWriter := writer.NewFileWriter()

	app := NewApp(
		s3Service,
		vpcService,
		ecsService,
		elbService,
		iamService,
		fileWriter,
		hclGenerator,
	)

	err = app.Run(ctx, resourceTypes, resourceName, clusterName)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
