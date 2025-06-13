// internal/aws/elasticloadbalancingv2/service.go
package elasticloadbalancingv2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
)

// --- Domain Models ---

type TargetGroup struct {
	Arn      string
	Name     string
	Port     int32
	Protocol types.ProtocolEnum
}

type Listener struct {
	Arn            string
	Port           int32
	Protocol       types.ProtocolEnum
	DefaultActions []types.Action // Simplified for now
}

type LoadBalancer struct {
	Arn          string
	Name         string
	Type         types.LoadBalancerTypeEnum
	Listeners    []Listener
	TargetGroups []TargetGroup
}

// --- Service Interface and Implementation ---

type ELBV2ServiceInterface interface {
	GetLoadBalancer(ctx context.Context, name string) (*LoadBalancer, error)
}

type ELBV2Service struct {
	repo ELBV2RepositoryInterface
}

func NewELBV2Service(repo ELBV2RepositoryInterface) *ELBV2Service {
	return &ELBV2Service{repo: repo}
}

func (s *ELBV2Service) GetLoadBalancer(ctx context.Context, name string) (*LoadBalancer, error) {
	// 1. Get Load Balancer
	awsLbs, err := s.repo.DescribeLoadBalancers(ctx, []string{name})
	if err != nil || len(awsLbs) == 0 {
		return nil, err
	}
	awsLb := awsLbs[0]
	lbArn := *awsLb.LoadBalancerArn

	// 2. Get Listeners
	awsListeners, err := s.repo.DescribeListeners(ctx, lbArn)
	if err != nil {
		return nil, err
	}

	var listeners []Listener
	for _, l := range awsListeners {
		listeners = append(listeners, Listener{
			Arn:            *l.ListenerArn,
			Port:           *l.Port,
			Protocol:       l.Protocol,
			DefaultActions: l.DefaultActions,
		})
	}

	// 3. Get Target Groups
	awsTgs, err := s.repo.DescribeTargetGroups(ctx, lbArn)
	if err != nil {
		return nil, err
	}

	var tgs []TargetGroup
	for _, tg := range awsTgs {
		tgs = append(tgs, TargetGroup{
			Arn:      *tg.TargetGroupArn,
			Name:     *tg.TargetGroupName,
			Port:     *tg.Port,
			Protocol: tg.Protocol,
		})
	}

	// 4. Assemble the final model
	lb := &LoadBalancer{
		Arn:          lbArn,
		Name:         *awsLb.LoadBalancerName,
		Type:         awsLb.Type,
		Listeners:    listeners,
		TargetGroups: tgs,
	}

	return lb, nil
} 
