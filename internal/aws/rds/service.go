// internal/aws/rds/service.go
package rds

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"golang.org/x/sync/errgroup"
)

// DBCluster はHCL生成に必要なDBクラスタの情報を保持します。
type DBCluster struct {
	Identifier              string
	Engine                  string
	EngineMode              string
	Tags                    map[string]string
	DBClusterParameterGroup string
	MemberIdentifiers       []string
}

// DBInstance はHCL生成に必要なDBインスタンスの情報を保持します。
type DBInstance struct {
	Identifier        string
	Engine            string
	InstanceClass     string
	Tags              map[string]string
	DBParameterGroups []string
}

// DBParameterGroup はHCL生成に必要なDBパラメータグループの情報を保持します。
type DBParameterGroup struct {
	Name   string
	Family string
	Tags   map[string]string
}

// Service はRDS関連のビジネスロジックを定義します。
type Service interface {
	ListDBClusters(ctx context.Context, dbClusterIdentifier string) ([]DBCluster, error)
	ListDBInstances(ctx context.Context, dbInstanceIdentifier string) ([]DBInstance, error)
	ListDBParameterGroups(ctx context.Context, dbParameterGroupName string) ([]DBParameterGroup, error)
}

// RDSService はServiceを実装します。
type RDSService struct {
	repo RDSRepositoryInterface
}

// NewRDSService は新しいRDSServiceを生成します。
func NewRDSService(repo RDSRepositoryInterface) *RDSService {
	return &RDSService{
		repo: repo,
	}
}

// ListDBClusters はDBクラスタのリストを取得し、ドメインオブジェクトに変換します。
func (s *RDSService) ListDBClusters(ctx context.Context, dbClusterIdentifier string) ([]DBCluster, error) {
	var id *string
	if dbClusterIdentifier != "" {
		id = aws.String(dbClusterIdentifier)
	}
	awsClusters, err := s.repo.DescribeDBClusters(ctx, id)
	if err != nil {
		return nil, err
	}

	var clusters []DBCluster
	for _, c := range awsClusters {
		var memberIdentifiers []string
		for _, member := range c.DBClusterMembers {
			memberIdentifiers = append(memberIdentifiers, *member.DBInstanceIdentifier)
		}

		var clusterPgName string
		if c.DBClusterParameterGroup != nil {
			clusterPgName = *c.DBClusterParameterGroup
		}

		clusters = append(clusters, DBCluster{
			Identifier:              *c.DBClusterIdentifier,
			Engine:                  *c.Engine,
			EngineMode:              *c.EngineMode,
			Tags:                    convertTags(c.TagList),
			DBClusterParameterGroup: clusterPgName,
			MemberIdentifiers:       memberIdentifiers,
		})
	}
	return clusters, nil
}

// ListDBInstances はDBインスタンスのリストを取得し、ドメインオブジェクトに変換します。
func (s *RDSService) ListDBInstances(ctx context.Context, dbInstanceIdentifier string) ([]DBInstance, error) {
	var id *string
	if dbInstanceIdentifier != "" {
		id = aws.String(dbInstanceIdentifier)
	}
	awsInstances, err := s.repo.DescribeDBInstances(ctx, id)
	if err != nil {
		return nil, err
	}

	var instances []DBInstance
	for _, i := range awsInstances {
		var pgNames []string
		for _, pg := range i.DBParameterGroups {
			pgNames = append(pgNames, *pg.DBParameterGroupName)
		}

		instances = append(instances, DBInstance{
			Identifier:        *i.DBInstanceIdentifier,
			Engine:            *i.Engine,
			InstanceClass:     *i.DBInstanceClass,
			Tags:              convertTags(i.TagList),
			DBParameterGroups: pgNames,
		})
	}
	return instances, nil
}

// ListDBParameterGroups はDBパラメータグループのリストを取得し、ドメインオブジェクトに変換します。
func (s *RDSService) ListDBParameterGroups(ctx context.Context, dbParameterGroupName string) ([]DBParameterGroup, error) {
	var name *string
	if dbParameterGroupName != "" {
		name = aws.String(dbParameterGroupName)
	}
	awsPGs, err := s.repo.DescribeDBParameterGroups(ctx, name)
	if err != nil {
		return nil, err
	}

	var eg errgroup.Group
	pgsWithTags := make(chan DBParameterGroup, len(awsPGs))

	for _, pg := range awsPGs {
		parameterGroup := pg
		eg.Go(func() error {
			tags, err := s.repo.ListTagsForResource(ctx, parameterGroup.DBParameterGroupArn)
			if err != nil {
				return err
			}
			pgDomain := DBParameterGroup{
				Name:   *parameterGroup.DBParameterGroupName,
				Family: *parameterGroup.DBParameterGroupFamily,
				Tags:   convertTags(tags),
			}
			pgsWithTags <- pgDomain
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}
	close(pgsWithTags)

	var pgs []DBParameterGroup
	for pg := range pgsWithTags {
		pgs = append(pgs, pg)
	}

	return pgs, nil
}

func convertTags(tags []types.Tag) map[string]string {
	m := make(map[string]string)
	for _, t := range tags {
		m[*t.Key] = *t.Value
	}
	return m
} 
