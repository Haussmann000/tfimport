// internal/di/container.go
package di

import (
	"context"
	"fmt"
	"strings"

	"github.com/Haussmann000/tfimport/internal/aws"
	"github.com/Haussmann000/tfimport/internal/aws/ec2"
	"github.com/Haussmann000/tfimport/internal/aws/ecs"
	"github.com/Haussmann000/tfimport/internal/aws/elbv2"
	"github.com/Haussmann000/tfimport/internal/aws/iam"
	"github.com/Haussmann000/tfimport/internal/aws/rds"
	"github.com/Haussmann000/tfimport/internal/aws/s3"
	"github.com/Haussmann000/tfimport/internal/hcl"
	"github.com/Haussmann000/tfimport/internal/writer"
	"golang.org/x/sync/errgroup"
)

// RunOptions はコマンドラインから渡されるオプションを保持します。
type RunOptions struct {
	ResourceTypes   []string
	ResourceName    string
	ClusterName     string
	ServiceName     string
	SecurityGroupID string
	DBClusterIdentifier string
	DBInstanceIdentifier string
}

// App はアプリケーションの主要なロジックをカプセル化します。
type App struct {
	s3Service  *s3.BucketService
	ec2Service *ec2.EC2Service
	ecsService *ecs.ECSService
	elbService *elbv2.ELBV2Service
	iamService *iam.IAMService
	rdsService *rds.RDSService
	writer     *writer.FileWriter
	generator  *hcl.HCLGenerator
}

// NewApp はAppのコンストラクタです。
func NewApp(
	s3s *s3.BucketService,
	es *ec2.EC2Service,
	ecss *ecs.ECSService,
	elbs *elbv2.ELBV2Service,
	iams *iam.IAMService,
	rdss *rds.RDSService,
	w *writer.FileWriter,
	g *hcl.HCLGenerator,
) *App {
	return &App{
		s3Service:  s3s,
		ec2Service: es,
		ecsService: ecss,
		elbService: elbs,
		iamService: iams,
		rdsService: rdss,
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
		case "security_group":
			if err := a.processSecurityGroup(ctx, options.SecurityGroupID); err != nil {
				return err
			}
		case "rds":
			if err := a.processRds(ctx, options); err != nil {
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
	vpcs, err := a.ec2Service.ListVpcs(ctx, resourceName)
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
		lbs, err = a.elbService.ListLoadBalancers(ctx, resourceName)
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

func (a *App) processSecurityGroup(ctx context.Context, sgIDsStr string) error {
	if sgIDsStr == "" {
		return nil
	}

	sgIDs := strings.Split(sgIDsStr, ",")
	var validSgIDs []string
	for _, id := range sgIDs {
		trimmedID := strings.TrimSpace(id)
		if trimmedID != "" {
			validSgIDs = append(validSgIDs, trimmedID)
		}
	}

	if len(validSgIDs) == 0 {
		return nil
	}

	sgs, err := a.ec2Service.ListSecurityGroups(ctx, validSgIDs)
	if err != nil {
		return err
	}

	if len(sgs) == 0 {
		return nil
	}

	hclFile, importFile, err := a.generator.GenerateSecurityGroupBlocks(sgs)
	if err != nil {
		return err
	}
	err = a.writer.WriteFile("security_group_generated.tf", hclFile)
	if err != nil {
		return err
	}
	return a.writer.WriteFile("security_group_import.tf", importFile)
}

func (a *App) processRds(ctx context.Context, options RunOptions) error {
	var clusters []rds.DBCluster
	var instances []rds.DBInstance
	var pgs []rds.DBParameterGroup
	var err error

	// Case 1: Specific cluster identifier is provided.
	if options.DBClusterIdentifier != "" {
		clusters, err = a.rdsService.ListDBClusters(ctx, options.DBClusterIdentifier)
		if err != nil {
			return err
		}
		if len(clusters) == 0 {
			fmt.Printf("No cluster found with identifier: %s\n", options.DBClusterIdentifier)
			return nil
		}

		pgNames := make(map[string]struct{})
		instanceIdentifiers := make(map[string]struct{})

		for _, c := range clusters {
			if c.DBClusterParameterGroup != "" {
				pgNames[c.DBClusterParameterGroup] = struct{}{}
			}
			for _, memberID := range c.MemberIdentifiers {
				instanceIdentifiers[memberID] = struct{}{}
			}
		}

		for id := range instanceIdentifiers {
			memberInstances, err := a.rdsService.ListDBInstances(ctx, id)
			if err != nil {
				return err
			}
			instances = append(instances, memberInstances...)
		}

		for _, inst := range instances {
			for _, pgName := range inst.DBParameterGroups {
				pgNames[pgName] = struct{}{}
			}
		}

		for name := range pgNames {
			paramGroups, err := a.rdsService.ListDBParameterGroups(ctx, name)
			if err != nil {
				return err
			}
			pgs = append(pgs, paramGroups...)
		}

	// Case 2: Specific instance identifier is provided (but not cluster).
	} else if options.DBInstanceIdentifier != "" {
		instances, err = a.rdsService.ListDBInstances(ctx, options.DBInstanceIdentifier)
		if err != nil {
			return err
		}

		pgNames := make(map[string]struct{})
		for _, inst := range instances {
			for _, pgName := range inst.DBParameterGroups {
				pgNames[pgName] = struct{}{}
			}
		}
		for name := range pgNames {
			paramGroups, err := a.rdsService.ListDBParameterGroups(ctx, name)
			if err != nil {
				return err
			}
			pgs = append(pgs, paramGroups...)
		}
	
	// Case 3: No specific identifier, fetch all.
	} else {
		var eg errgroup.Group

		eg.Go(func() error {
			var err error
			clusters, err = a.rdsService.ListDBClusters(ctx, "")
			return err
		})

		eg.Go(func() error {
			var err error
			instances, err = a.rdsService.ListDBInstances(ctx, "")
			return err
		})

		eg.Go(func() error {
			var err error
			pgs, err = a.rdsService.ListDBParameterGroups(ctx, options.ResourceName)
			return err
		})

		if err := eg.Wait(); err != nil {
			return err
		}
	}

	hclFile, importFile, err := a.generator.GenerateRdsBlocks(clusters, instances, pgs)
	if err != nil {
		return err
	}

	err = a.writer.WriteFile("rds_generated.tf", hclFile)
	if err != nil {
		return err
	}
	return a.writer.WriteFile("rds_import.tf", importFile)
}

// BuildApp は依存関係を解決してAppを構築します。
func BuildApp(ctx context.Context) (*App, error) {
	awsCfg, err := aws.NewConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS config: %w", err)
	}

	// S3
	s3Client := aws.NewS3Client(awsCfg)
	s3Repo := s3.NewS3Repository(s3Client)
	s3Service := s3.NewBucketService(s3Repo)

	// EC2
	ec2Client := aws.NewEC2Client(awsCfg)
	ec2Repo := ec2.NewEC2Repository(ec2Client)
	ec2Service := ec2.NewEC2Service(ec2Repo)

	// ECS
	ecsClient := aws.NewECSClient(awsCfg)
	ecsRepo := ecs.NewECSRepository(ecsClient)
	ecsService := ecs.NewECSService(ecsRepo)

	// ELBv2
	elbv2Client := aws.NewELBV2Client(awsCfg)
	elbv2Repo := elbv2.NewELBV2Repository(elbv2Client)
	elbService := elbv2.NewELBV2Service(elbv2Repo)

	// IAM
	iamClient := aws.NewIAMClient(awsCfg)
	iamRepo := iam.NewIAMRepository(iamClient)
	iamService := iam.NewIAMService(iamRepo)

	// rds
	rdsClient := aws.NewRDSClient(awsCfg)
	rdsRepo := rds.NewRDSRepository(rdsClient)
	rdsService := rds.NewRDSService(rdsRepo)

	writer := writer.NewFileWriter()
	generator := hcl.NewHCLGenerator()

	app := NewApp(s3Service, ec2Service, ecsService, elbService, iamService, rdsService, writer, generator)

	return app, nil
} 
