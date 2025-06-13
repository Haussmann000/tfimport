// internal/aws/ecs/repository.go
package ecs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

// ECSRepositoryInterface はECSリソースへのアクセスを抽象化します。
type ECSRepositoryInterface interface {
	ListClusters(ctx context.Context) ([]string, error)
	DescribeClusters(ctx context.Context, clusterNames []string) ([]types.Cluster, error)
	ListServices(ctx context.Context, clusterArn string) ([]string, error)
	DescribeServices(ctx context.Context, clusterArn string, serviceArns []string) ([]types.Service, error)
	DescribeCluster(ctx context.Context, clusterName string) (*types.Cluster, error)
}

// ECSRepository はECSRepositoryInterfaceを実装します。
type ECSRepository struct {
	client *ecs.Client
}

// NewECSRepository は新しいECSRepositoryを生成します。
func NewECSRepository(client *ecs.Client) *ECSRepository {
	return &ECSRepository{client: client}
}

func (r *ECSRepository) ListClusters(ctx context.Context) ([]string, error) {
	input := &ecs.ListClustersInput{}
	result, err := r.client.ListClusters(ctx, input)
	if err != nil {
		return nil, err
	}
	return result.ClusterArns, nil
}

func (r *ECSRepository) DescribeClusters(ctx context.Context, clusterNames []string) ([]types.Cluster, error) {
	if len(clusterNames) == 0 {
		return nil, nil
	}
	input := &ecs.DescribeClustersInput{
		Clusters: clusterNames,
		Include: []types.ClusterField{
			types.ClusterFieldTags,
		},
	}
	result, err := r.client.DescribeClusters(ctx, input)
	if err != nil {
		return nil, err
	}
	return result.Clusters, nil
}

func (r *ECSRepository) ListServices(ctx context.Context, clusterArn string) ([]string, error) {
	input := &ecs.ListServicesInput{
		Cluster: &clusterArn,
	}
	result, err := r.client.ListServices(ctx, input)
	if err != nil {
		return nil, err
	}
	return result.ServiceArns, nil
}

func (r *ECSRepository) DescribeServices(ctx context.Context, clusterArn string, serviceArns []string) ([]types.Service, error) {
	if len(serviceArns) == 0 {
		return nil, nil
	}
	input := &ecs.DescribeServicesInput{
		Cluster:  &clusterArn,
		Services: serviceArns,
	}
	result, err := r.client.DescribeServices(ctx, input)
	if err != nil {
		return nil, err
	}
	return result.Services, nil
}

func (r *ECSRepository) DescribeCluster(ctx context.Context, clusterName string) (*types.Cluster, error) {
	input := &ecs.DescribeClustersInput{
		Clusters: []string{clusterName},
		Include: []types.ClusterField{
			types.ClusterFieldTags,
		},
	}
	result, err := r.client.DescribeClusters(ctx, input)
	if err != nil {
		return nil, err
	}
	if len(result.Clusters) == 0 {
		return nil, nil
	}
	return &result.Clusters[0], nil
} 
