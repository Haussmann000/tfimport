package iam

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

type IAMClientInterface interface {
	ListRoles(ctx context.Context, params *iam.ListRolesInput, optFns ...func(*iam.Options)) (*iam.ListRolesOutput, error)
	ListPolicies(ctx context.Context, params *iam.ListPoliciesInput, optFns ...func(*iam.Options)) (*iam.ListPoliciesOutput, error)
	ListAttachedRolePolicies(ctx context.Context, params *iam.ListAttachedRolePoliciesInput, optFns ...func(*iam.Options)) (*iam.ListAttachedRolePoliciesOutput, error)
	GetPolicy(ctx context.Context, params *iam.GetPolicyInput, optFns ...func(*iam.Options)) (*iam.GetPolicyOutput, error)
	GetPolicyVersion(ctx context.Context, params *iam.GetPolicyVersionInput, optFns ...func(*iam.Options)) (*iam.GetPolicyVersionOutput, error)
}

type IAMRepositoryInterface interface {
	ListRoles(ctx context.Context) ([]types.Role, error)
	ListPolicies(ctx context.Context, scope types.PolicyScopeType) ([]types.Policy, error)
	ListAttachedRolePolicies(ctx context.Context, roleName string) ([]types.AttachedPolicy, error)
	GetPolicy(ctx context.Context, policyArn string) (*types.Policy, error)
	GetPolicyVersion(ctx context.Context, policyArn string, versionId string) (*types.PolicyVersion, error)
}

type IAMRepository struct {
	client IAMClientInterface
}

func NewIAMRepository(client IAMClientInterface) *IAMRepository {
	return &IAMRepository{client: client}
}

func (r *IAMRepository) ListRoles(ctx context.Context) ([]types.Role, error) {
	var roles []types.Role
	paginator := iam.NewListRolesPaginator(r.client, &iam.ListRolesInput{})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		roles = append(roles, output.Roles...)
	}

	return roles, nil
}

func (r *IAMRepository) ListPolicies(ctx context.Context, scope types.PolicyScopeType) ([]types.Policy, error) {
	var policies []types.Policy
	paginator := iam.NewListPoliciesPaginator(r.client, &iam.ListPoliciesInput{
		Scope: scope,
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		policies = append(policies, output.Policies...)
	}

	return policies, nil
}

func (r *IAMRepository) ListAttachedRolePolicies(ctx context.Context, roleName string) ([]types.AttachedPolicy, error) {
	var attachedPolicies []types.AttachedPolicy
	paginator := iam.NewListAttachedRolePoliciesPaginator(r.client, &iam.ListAttachedRolePoliciesInput{
		RoleName: aws.String(roleName),
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		attachedPolicies = append(attachedPolicies, output.AttachedPolicies...)
	}

	return attachedPolicies, nil
}

func (r *IAMRepository) GetPolicy(ctx context.Context, policyArn string) (*types.Policy, error) {
	output, err := r.client.GetPolicy(ctx, &iam.GetPolicyInput{
		PolicyArn: aws.String(policyArn),
	})
	if err != nil {
		return nil, err
	}
	return output.Policy, nil
}

func (r *IAMRepository) GetPolicyVersion(ctx context.Context, policyArn string, versionId string) (*types.PolicyVersion, error) {
	output, err := r.client.GetPolicyVersion(ctx, &iam.GetPolicyVersionInput{
		PolicyArn: aws.String(policyArn),
		VersionId: aws.String(versionId),
	})
	if err != nil {
		return nil, err
	}
	return output.PolicyVersion, nil
} 
