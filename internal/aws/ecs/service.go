// internal/aws/ecs/service.go
package ecs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

// --- Domain Models ---

type LoadBalancer struct {
	TargetGroupArn string
	ContainerName  string
	ContainerPort  int32
}

type DeploymentCircuitBreaker struct {
	Enable   bool
	Rollback bool
}

type NetworkConfiguration struct {
	Subnets        []string
	SecurityGroups []string
	AssignPublicIp bool
}

type ServiceDetail struct {
	Arn                              string
	Name                             string
	DesiredCount                     int32
	Tags                             map[string]string
	LoadBalancers                    []LoadBalancer
	TaskDefinitionArn                string
	EnableEcsManagedTags             bool
	EnableExecuteCommand             bool
	HealthCheckGracePeriodSeconds    int32
	DeploymentCircuitBreaker         *DeploymentCircuitBreaker
	NetworkConfiguration             *NetworkConfiguration
	PropagateTags                    string
	PlatformVersion                  string
	SchedulingStrategy               string
}

type Cluster struct {
	Arn      string
	Name     string
	Services []ServiceDetail
	Tags     map[string]string
}

// --- Service Interface and Implementation ---

type Service interface {
	GetClusters(ctx context.Context, clusterName, serviceName string) ([]Cluster, error)
}

type ECSService struct {
	repo ECSRepositoryInterface
}

func NewECSService(repo ECSRepositoryInterface) *ECSService {
	return &ECSService{repo: repo}
}

// GetClustersは指定されたECSクラスターとそのサービスを取得します。
func (s *ECSService) GetClusters(ctx context.Context, clusterName, serviceName string) ([]Cluster, error) {
	awsClusters, err := s.repo.DescribeClusters(ctx, []string{clusterName})
	if err != nil {
		return nil, err
	}
	if len(awsClusters) == 0 {
		return nil, nil // No cluster found
	}
	targetCluster := awsClusters[0]
	clusterArn := *targetCluster.ClusterArn

	// サービスを取得
	var serviceArns []string
	if serviceName != "" {
		serviceArns = []string{serviceName}
	} else {
		serviceArns, err = s.repo.ListServices(ctx, clusterArn)
		if err != nil {
			return nil, err
		}
	}

	awsServices, err := s.repo.DescribeServices(ctx, clusterArn, serviceArns)
	if err != nil {
		return nil, err
	}

	var services []ServiceDetail
	for _, awsService := range awsServices {
		// Tags
		tags := make(map[string]string)
		for _, tag := range awsService.Tags {
			tags[*tag.Key] = *tag.Value
		}

		// LoadBalancers
		var lbs []LoadBalancer
		for _, lb := range awsService.LoadBalancers {
			lbs = append(lbs, LoadBalancer{
				TargetGroupArn: *lb.TargetGroupArn,
				ContainerName:  *lb.ContainerName,
				ContainerPort:  *lb.ContainerPort,
			})
		}

		// DeploymentCircuitBreaker
		var circuitBreaker *DeploymentCircuitBreaker
		if awsService.DeploymentConfiguration != nil && awsService.DeploymentConfiguration.DeploymentCircuitBreaker != nil {
			breaker := awsService.DeploymentConfiguration.DeploymentCircuitBreaker
			circuitBreaker = &DeploymentCircuitBreaker{
				Enable:   breaker.Enable,
				Rollback: breaker.Rollback,
			}
		}

		// NetworkConfiguration
		var networkConfig *NetworkConfiguration
		if awsService.NetworkConfiguration != nil && awsService.NetworkConfiguration.AwsvpcConfiguration != nil {
			vpcConfig := awsService.NetworkConfiguration.AwsvpcConfiguration
			networkConfig = &NetworkConfiguration{
				Subnets:        vpcConfig.Subnets,
				SecurityGroups: vpcConfig.SecurityGroups,
				AssignPublicIp: vpcConfig.AssignPublicIp == types.AssignPublicIpEnabled,
			}
		}

		// Simple boolean and integer values
		healthCheckGracePeriodSeconds := int32(0)
		if awsService.HealthCheckGracePeriodSeconds != nil {
			healthCheckGracePeriodSeconds = *awsService.HealthCheckGracePeriodSeconds
		}

		service := ServiceDetail{
			Arn:                           *awsService.ServiceArn,
			Name:                          *awsService.ServiceName,
			DesiredCount:                  awsService.DesiredCount,
			Tags:                          tags,
			LoadBalancers:                 lbs,
			TaskDefinitionArn:             *awsService.TaskDefinition,
			EnableEcsManagedTags:          awsService.EnableECSManagedTags,
			EnableExecuteCommand:          awsService.EnableExecuteCommand,
			HealthCheckGracePeriodSeconds: healthCheckGracePeriodSeconds,
			DeploymentCircuitBreaker:      circuitBreaker,
			NetworkConfiguration:          networkConfig,
			PropagateTags:                 string(awsService.PropagateTags),
			PlatformVersion:               *awsService.PlatformVersion,
			SchedulingStrategy:            string(awsService.SchedulingStrategy),
		}

		for i, lb := range lbs {
			service.LoadBalancers[i] = lb
		}

		services = append(services, service)
	}

	cluster := Cluster{
		Arn:      clusterArn,
		Name:     *targetCluster.ClusterName,
		Services: services,
		Tags:     make(map[string]string),
	}

	for _, tag := range targetCluster.Tags {
		cluster.Tags[*tag.Key] = *tag.Value
	}

	return []Cluster{cluster}, nil
} 
