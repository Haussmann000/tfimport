// internal/aws/rds/repository.go
package rds

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
)

// RDSRepositoryInterface はRDSリソースへのアクセスを抽象化します。
type RDSRepositoryInterface interface {
	DescribeDBClusters(ctx context.Context, dbClusterIdentifier *string) ([]types.DBCluster, error)
	DescribeDBInstances(ctx context.Context, dbInstanceIdentifier *string) ([]types.DBInstance, error)
	DescribeDBParameterGroups(ctx context.Context, dbParameterGroupName *string) ([]types.DBParameterGroup, error)
	ListTagsForResource(ctx context.Context, resourceName *string) ([]types.Tag, error)
}

// RDSRepository はRDSRepositoryInterfaceを実装します。
type RDSRepository struct {
	client *rds.Client
}

// NewRDSRepository は新しいRDSRepositoryを生成します。
func NewRDSRepository(client *rds.Client) *RDSRepository {
	return &RDSRepository{
		client: client,
	}
}

// DescribeDBClusters はAWSからDBClusterのリストを取得します。
func (r *RDSRepository) DescribeDBClusters(ctx context.Context, dbClusterIdentifier *string) ([]types.DBCluster, error) {
	input := &rds.DescribeDBClustersInput{
		DBClusterIdentifier: dbClusterIdentifier,
	}
	result, err := r.client.DescribeDBClusters(ctx, input)
	if err != nil {
		return nil, err
	}
	return result.DBClusters, nil
}

// DescribeDBInstances はAWSからDBInstanceのリストを取得します。
func (r *RDSRepository) DescribeDBInstances(ctx context.Context, dbInstanceIdentifier *string) ([]types.DBInstance, error) {
	input := &rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: dbInstanceIdentifier,
	}
	result, err := r.client.DescribeDBInstances(ctx, input)
	if err != nil {
		return nil, err
	}
	return result.DBInstances, nil
}

// DescribeDBParameterGroups はAWSからDBParameterGroupのリストを取得します。
func (r *RDSRepository) DescribeDBParameterGroups(ctx context.Context, dbParameterGroupName *string) ([]types.DBParameterGroup, error) {
	input := &rds.DescribeDBParameterGroupsInput{
		DBParameterGroupName: dbParameterGroupName,
	}
	result, err := r.client.DescribeDBParameterGroups(ctx, input)
	if err != nil {
		return nil, err
	}
	return result.DBParameterGroups, nil
}

// ListTagsForResource はAWSからリソースのタグリストを取得します。
func (r *RDSRepository) ListTagsForResource(ctx context.Context, resourceName *string) ([]types.Tag, error) {
	input := &rds.ListTagsForResourceInput{
		ResourceName: resourceName,
	}
	result, err := r.client.ListTagsForResource(ctx, input)
	if err != nil {
		return nil, err
	}
	return result.TagList, nil
} 
