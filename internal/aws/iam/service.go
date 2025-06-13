package iam

import (
	"context"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"golang.org/x/sync/errgroup"
)

type Policy struct {
	Name           string
	Arn            string
	PolicyDocument string
}

type Role struct {
	Name                 string
	Arn                  string
	AssumeRolePolicy     string
	AttachedPolicyArns   []string
}

type Service interface {
	ListRoles(ctx context.Context, nameContains string) ([]Role, error)
	ListPolicies(ctx context.Context, nameContains string) ([]Policy, error)
}

type IAMService struct {
	iamRepo IAMRepositoryInterface
}

func NewIAMService(repo IAMRepositoryInterface) *IAMService {
	return &IAMService{iamRepo: repo}
}

func (s *IAMService) ListRoles(ctx context.Context, nameContains string) ([]Role, error) {
	awsRoles, err := s.iamRepo.ListRoles(ctx)
	if err != nil {
		return nil, err
	}

	var filteredRoles []types.Role
	if nameContains != "" {
		for _, r := range awsRoles {
			if strings.Contains(*r.RoleName, nameContains) {
				filteredRoles = append(filteredRoles, r)
			}
		}
	} else {
		filteredRoles = awsRoles
	}

	var roles []Role
	for _, r := range filteredRoles {
		assumeRolePolicy, err := url.QueryUnescape(*r.AssumeRolePolicyDocument)
		if err != nil {
			return nil, err
		}

		attachedPolicies, err := s.iamRepo.ListAttachedRolePolicies(ctx, *r.RoleName)
		if err != nil {
			return nil, err
		}

		var attachedPolicyArns []string
		for _, p := range attachedPolicies {
			attachedPolicyArns = append(attachedPolicyArns, *p.PolicyArn)
		}

		roles = append(roles, Role{
			Name:                 *r.RoleName,
			Arn:                  *r.Arn,
			AssumeRolePolicy:     assumeRolePolicy,
			AttachedPolicyArns:   attachedPolicyArns,
		})
	}

	return roles, nil
}

func (s *IAMService) ListPolicies(ctx context.Context, nameContains string) ([]Policy, error) {
	awsPolicies, err := s.iamRepo.ListPolicies(ctx, types.PolicyScopeTypeLocal)
	if err != nil {
		return nil, err
	}

	var filteredPolicies []types.Policy
	if nameContains != "" {
		for _, p := range awsPolicies {
			if strings.Contains(*p.PolicyName, nameContains) {
				filteredPolicies = append(filteredPolicies, p)
			}
		}
	} else {
		filteredPolicies = awsPolicies
	}

	var policies []Policy
	var eg errgroup.Group
	policyChan := make(chan Policy, len(filteredPolicies))

	for _, p := range filteredPolicies {
		policy := p
		eg.Go(func() error {
			policyVersion, err := s.iamRepo.GetPolicyVersion(ctx, *policy.Arn, *policy.DefaultVersionId)
			if err != nil {
				return err
			}
			policyDocument, err := url.QueryUnescape(*policyVersion.Document)
			if err != nil {
				return err
			}
			policyChan <- Policy{
				Name:           *policy.PolicyName,
				Arn:            *policy.Arn,
				PolicyDocument: policyDocument,
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}
	close(policyChan)

	for p := range policyChan {
		policies = append(policies, p)
	}

	return policies, nil
} 
