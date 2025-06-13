// internal/aws/elasticloadbalancingv2/repository.go
package elasticloadbalancingv2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
)

// ELBV2RepositoryInterface はELBV2リソースへのアクセスを抽象化します。
type ELBV2RepositoryInterface interface {
	DescribeLoadBalancers(ctx context.Context, names []string) ([]types.LoadBalancer, error)
	DescribeListeners(ctx context.Context, loadBalancerArn string) ([]types.Listener, error)
	DescribeTargetGroups(ctx context.Context, loadBalancerArn string) ([]types.TargetGroup, error)
}

// ELBV2Repository はELBV2RepositoryInterfaceを実装します。
type ELBV2Repository struct {
	client *elasticloadbalancingv2.Client
}

// NewELBV2Repository は新しいELBV2Repositoryを生成します。
func NewELBV2Repository(client *elasticloadbalancingv2.Client) *ELBV2Repository {
	return &ELBV2Repository{client: client}
}

func (r *ELBV2Repository) DescribeLoadBalancers(ctx context.Context, names []string) ([]types.LoadBalancer, error) {
	input := &elasticloadbalancingv2.DescribeLoadBalancersInput{
		Names: names,
	}
	result, err := r.client.DescribeLoadBalancers(ctx, input)
	if err != nil {
		return nil, err
	}
	return result.LoadBalancers, nil
}

func (r *ELBV2Repository) DescribeListeners(ctx context.Context, loadBalancerArn string) ([]types.Listener, error) {
	input := &elasticloadbalancingv2.DescribeListenersInput{
		LoadBalancerArn: &loadBalancerArn,
	}
	result, err := r.client.DescribeListeners(ctx, input)
	if err != nil {
		return nil, err
	}
	return result.Listeners, nil
}

func (r *ELBV2Repository) DescribeTargetGroups(ctx context.Context, loadBalancerArn string) ([]types.TargetGroup, error) {
	input := &elasticloadbalancingv2.DescribeTargetGroupsInput{
		LoadBalancerArn: &loadBalancerArn,
	}
	result, err := r.client.DescribeTargetGroups(ctx, input)
	if err != nil {
		return nil, err
	}
	return result.TargetGroups, nil
} 
