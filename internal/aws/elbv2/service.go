package elbv2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"golang.org/x/sync/errgroup"
)

// --- Domain Models ---

type TargetGroup struct {
	Arn      string
	Name     string
	Port     int32
	Protocol types.ProtocolEnum
	VpcId    string
	TargetType string
	HealthCheck *HealthCheck
}

type HealthCheck struct {
	Enabled bool
	Path string
	Port string
	Protocol types.ProtocolEnum
	Interval int32
	Timeout int32
	HealthyThreshold int32
	UnhealthyThreshold int32
	Matcher string
}

type TargetGroupStickiness struct {
	Enabled  bool
	Duration int32
}

type DefaultActionForwardTargetGroup struct {
	Arn    string
	Weight *int64
}

type DefaultActionForward struct {
	TargetGroups []DefaultActionForwardTargetGroup
	Stickiness   *TargetGroupStickiness
}

type DefaultActionRedirect struct {
	Port       string
	Protocol   string
	StatusCode string
}

type DefaultAction struct {
	Type     types.ActionTypeEnum
	Forward  *DefaultActionForward
	Redirect *DefaultActionRedirect
}

type ListenerRuleActionForward struct {
	TargetGroupArn string
}

type ListenerRuleAction struct {
	Type    types.ActionTypeEnum
	Forward *ListenerRuleActionForward
}

type ListenerRuleCondition struct {
	Field  string
	Values []string
}

type ListenerRule struct {
	Arn        string
	Priority   string
	Actions    []ListenerRuleAction
	Conditions []ListenerRuleCondition
}

type Listener struct {
	Arn            string
	Port           int32
	Protocol       types.ProtocolEnum
	DefaultActions []DefaultAction
	CertificateArn *string
	Rules          []ListenerRule
}

type LoadBalancer struct {
	Arn          string
	Name         string
	Type         types.LoadBalancerTypeEnum
	Listeners    []Listener
	TargetGroups []TargetGroup
	Subnets      []string
	VpcId        string
}

// --- Service Interface and Implementation ---

type Service interface {
	GetLoadBalancer(ctx context.Context, name string) (*LoadBalancer, error)
	ListLoadBalancers(ctx context.Context, name string) ([]*LoadBalancer, error)
}

type ELBV2Service struct {
	repo ELBV2RepositoryInterface
}

func NewELBV2Service(repo ELBV2RepositoryInterface) *ELBV2Service {
	return &ELBV2Service{repo: repo}
}

func (s *ELBV2Service) ListLoadBalancers(ctx context.Context, name string) ([]*LoadBalancer, error) {
	var names []string
	if name != "" {
		names = append(names, name)
	}

	awsLbs, err := s.repo.DescribeLoadBalancers(ctx, names)
	if err != nil {
		return nil, err
	}

	var lbs []*LoadBalancer
	for _, awsLb := range awsLbs {
		lb, err := s.buildLoadBalancer(ctx, awsLb)
		if err != nil {
			return nil, err // Or handle more gracefully
		}
		lbs = append(lbs, lb)
	}

	return lbs, nil
}

func (s *ELBV2Service) GetLoadBalancer(ctx context.Context, name string) (*LoadBalancer, error) {
	lbs, err := s.ListLoadBalancers(ctx, name)
	if err != nil || len(lbs) == 0 {
		return nil, err
	}
	return lbs[0], nil
}

func (s *ELBV2Service) buildLoadBalancer(ctx context.Context, awsLb types.LoadBalancer) (*LoadBalancer, error) {
	lb := &LoadBalancer{
		Name:         *awsLb.LoadBalancerName,
		Arn:          *awsLb.LoadBalancerArn,
		Subnets:      getSubnetIDs(awsLb.AvailabilityZones),
		VpcId:        *awsLb.VpcId,
	}

	var eg errgroup.Group

	eg.Go(func() error {
		listeners, err := s.getListenersWithRules(ctx, lb.Arn)
		if err != nil {
			return err
		}
		lb.Listeners = listeners
		return nil
	})

	eg.Go(func() error {
		tgs, err := s.getTargetGroups(ctx, lb.Arn)
		if err != nil {
			return err
		}
		lb.TargetGroups = tgs
		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return lb, nil
}

func (s *ELBV2Service) getListenersWithRules(ctx context.Context, lbArn string) ([]Listener, error) {
	awsListeners, err := s.repo.DescribeListeners(ctx, lbArn)
	if err != nil {
		return nil, err
	}

	var listeners []Listener
	var eg errgroup.Group
	listenerChan := make(chan Listener, len(awsListeners))

	for _, l := range awsListeners {
		listener := l
		eg.Go(func() error {
			rules, err := s.repo.DescribeRules(ctx, *listener.ListenerArn)
			if err != nil {
				return err
			}

			listenerChan <- Listener{
				Arn:            *listener.ListenerArn,
				Port:           *listener.Port,
				Protocol:       listener.Protocol,
				CertificateArn: getCertificateArn(listener.Certificates),
				DefaultActions: convertActions(listener.DefaultActions),
				Rules:          convertRules(rules),
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}
	close(listenerChan)

	for l := range listenerChan {
		listeners = append(listeners, l)
	}

	return listeners, nil
}

func (s *ELBV2Service) getTargetGroups(ctx context.Context, lbArn string) ([]TargetGroup, error) {
	awsTgs, err := s.repo.DescribeTargetGroups(ctx, lbArn)
	if err != nil {
		return nil, err
	}
	var tgs []TargetGroup
	for _, tg := range awsTgs {
		tgs = append(tgs, TargetGroup{
			Name:        *tg.TargetGroupName,
			Arn:         *tg.TargetGroupArn,
			Port:        *tg.Port,
			Protocol:    tg.Protocol,
			TargetType:  string(tg.TargetType),
			VpcId:       *tg.VpcId,
			HealthCheck: convertHealthCheck(tg),
		})
	}
	return tgs, nil
}

func getSubnetIDs(zones []types.AvailabilityZone) []string {
	var subnetIDs []string
	for _, z := range zones {
		subnetIDs = append(subnetIDs, *z.SubnetId)
	}
	return subnetIDs
}

func getCertificateArn(certificates []types.Certificate) *string {
	if len(certificates) > 0 {
		return certificates[0].CertificateArn
	}
	return nil
}

func convertActions(actions []types.Action) []DefaultAction {
	var defaultActions []DefaultAction
	for _, da := range actions {
		action := DefaultAction{Type: da.Type}
		if da.Type == types.ActionTypeEnumForward && da.ForwardConfig != nil {
			forward := DefaultActionForward{}
			for _, tg := range da.ForwardConfig.TargetGroups {
				var weight *int64
				if tg.Weight != nil {
					w := int64(*tg.Weight)
					weight = &w
				}
				forward.TargetGroups = append(forward.TargetGroups, DefaultActionForwardTargetGroup{
					Arn:    *tg.TargetGroupArn,
					Weight: weight,
				})
			}
			if da.ForwardConfig.TargetGroupStickinessConfig != nil {
				forward.Stickiness = &TargetGroupStickiness{
					Enabled:  *da.ForwardConfig.TargetGroupStickinessConfig.Enabled,
					Duration: *da.ForwardConfig.TargetGroupStickinessConfig.DurationSeconds,
				}
			}
			action.Forward = &forward
		} else if da.Type == types.ActionTypeEnumRedirect && da.RedirectConfig != nil {
			action.Redirect = &DefaultActionRedirect{
				Port:       *da.RedirectConfig.Port,
				Protocol:   *da.RedirectConfig.Protocol,
				StatusCode: string(da.RedirectConfig.StatusCode),
			}
		}
		defaultActions = append(defaultActions, action)
	}
	return defaultActions
}

func convertRules(rules []types.Rule) []ListenerRule {
	var listenerRules []ListenerRule
	for _, r := range rules {
		if (r.IsDefault != nil && *r.IsDefault) || len(r.Actions) == 0 || len(r.Conditions) == 0 {
			continue
		}

		var actions []ListenerRuleAction
		for _, a := range r.Actions {
			if a.Type == types.ActionTypeEnumForward && a.ForwardConfig != nil && len(a.ForwardConfig.TargetGroups) > 0 {
				actions = append(actions, ListenerRuleAction{
					Type: a.Type,
					Forward: &ListenerRuleActionForward{
						TargetGroupArn: *a.ForwardConfig.TargetGroups[0].TargetGroupArn,
					},
				})
			}
		}

		var conditions []ListenerRuleCondition
		for _, c := range r.Conditions {
			if c.HostHeaderConfig != nil {
				conditions = append(conditions, ListenerRuleCondition{
					Field:  "host-header",
					Values: c.HostHeaderConfig.Values,
				})
			}
			if c.PathPatternConfig != nil {
				conditions = append(conditions, ListenerRuleCondition{
					Field:  "path-pattern",
					Values: c.PathPatternConfig.Values,
				})
			}
		}

		listenerRules = append(listenerRules, ListenerRule{
			Arn:        *r.RuleArn,
			Priority:   *r.Priority,
			Actions:    actions,
			Conditions: conditions,
		})
	}
	return listenerRules
}

func convertHealthCheck(tg types.TargetGroup) *HealthCheck {
	if tg.HealthCheckEnabled != nil && *tg.HealthCheckEnabled {
		return &HealthCheck{
			Enabled:            *tg.HealthCheckEnabled,
			Path:               *tg.HealthCheckPath,
			Port:               *tg.HealthCheckPort,
			Protocol:           tg.HealthCheckProtocol,
			Interval:           *tg.HealthCheckIntervalSeconds,
			Timeout:            *tg.HealthCheckTimeoutSeconds,
			HealthyThreshold:   *tg.HealthyThresholdCount,
			UnhealthyThreshold: *tg.UnhealthyThresholdCount,
			Matcher:            *tg.Matcher.HttpCode,
		}
	}
	return nil
} 
