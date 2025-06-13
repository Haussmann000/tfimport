package elbv2

import (
	"context"

	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
)

// ELBV2RepositoryInterface はELBV2リソースへのアクセスを抽象化します。
type ELBV2RepositoryInterface interface {
	DescribeLoadBalancers(ctx context.Context, names []string) ([]types.LoadBalancer, error)
	DescribeListeners(ctx context.Context, loadBalancerArn string) ([]types.Listener, error)
	DescribeRules(ctx context.Context, listenerArn string) ([]types.Rule, error)
	DescribeTargetGroups(ctx context.Context, loadBalancerArn string) ([]types.TargetGroup, error)
}

// ELBV2Repository はELBV2RepositoryInterfaceを実装します。
type ELBV2Repository struct {
	client *elbv2.Client
}

// NewELBV2Repository は新しいELBV2Repositoryを生成します。
func NewELBV2Repository(client *elbv2.Client) *ELBV2Repository {
	return &ELBV2Repository{client: client}
}

func (r *ELBV2Repository) DescribeLoadBalancers(ctx context.Context, names []string) ([]types.LoadBalancer, error) {
	input := &elbv2.DescribeLoadBalancersInput{
		Names: names,
	}
	result, err := r.client.DescribeLoadBalancers(ctx, input)
	if err != nil {
		return nil, err
	}
	return result.LoadBalancers, nil
}

func (r *ELBV2Repository) DescribeListeners(ctx context.Context, loadBalancerArn string) ([]types.Listener, error) {
	input := &elbv2.DescribeListenersInput{
		LoadBalancerArn: &loadBalancerArn,
	}
	result, err := r.client.DescribeListeners(ctx, input)
	if err != nil {
		return nil, err
	}
	return result.Listeners, nil
}

func (r *ELBV2Repository) DescribeRules(ctx context.Context, listenerArn string) ([]types.Rule, error) {
	input := &elbv2.DescribeRulesInput{
		ListenerArn: &listenerArn,
	}

	var rules []types.Rule
	paginator := elbv2.NewDescribeRulesPaginator(r.client, input)
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		rules = append(rules, output.Rules...)
	}
	return rules, nil
}

func (r *ELBV2Repository) DescribeTargetGroups(ctx context.Context, loadBalancerArn string) ([]types.TargetGroup, error) {
	input := &elbv2.DescribeTargetGroupsInput{
		LoadBalancerArn: &loadBalancerArn,
	}
	result, err := r.client.DescribeTargetGroups(ctx, input)
	if err != nil {
		return nil, err
	}
	return result.TargetGroups, nil
} 
