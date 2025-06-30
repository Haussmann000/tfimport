package hcl

import (
	"fmt"
	"strings"

	"github.com/Haussmann000/tfimport/internal/aws/ec2"
	"github.com/Haussmann000/tfimport/internal/aws/ecs"
	"github.com/Haussmann000/tfimport/internal/aws/elbv2"
	"github.com/Haussmann000/tfimport/internal/aws/iam"
	"github.com/Haussmann000/tfimport/internal/aws/s3"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

// HCLGenerator はHCLブロックを生成します。
type HCLGenerator struct{}

// NewHCLGenerator は新しいHCLGeneratorを生成します。
func NewHCLGenerator() *HCLGenerator {
	return &HCLGenerator{}
}

// GenerateVpcBlocks はVPCリソースのresourceブロックとimportブロックを生成します。
func (g *HCLGenerator) GenerateVpcBlocks(vpcs []ec2.Vpc) (*hclwrite.File, *hclwrite.File, error) {
	resourceFile := hclwrite.NewEmptyFile()
	importFile := hclwrite.NewEmptyFile()
	resourceBody := resourceFile.Body()
	importBody := importFile.Body()

	for i, vpc := range vpcs {
		resourceName := fmt.Sprintf("vpc_%d", i)
		resourceType := "aws_vpc"

		// importブロックの生成
		importBlock := importBody.AppendNewBlock("import", nil)
		importBody.AppendNewline()
		importBlock.Body().SetAttributeValue("to", cty.StringVal(resourceType+"."+resourceName))
		importBlock.Body().SetAttributeValue("id", cty.StringVal(vpc.ID))

		// resourceブロックの生成
		vpcBlock := resourceBody.AppendNewBlock("resource", []string{resourceType, resourceName})
		resourceBody.AppendNewline()
		vpcBlock.Body().SetAttributeValue("cidr_block", cty.StringVal(vpc.CidrBlock))

		if len(vpc.Tags) > 0 {
			tagMap := make(map[string]cty.Value)
			for k, v := range vpc.Tags {
				tagMap[k] = cty.StringVal(v)
			}
			vpcBlock.Body().SetAttributeValue("tags", cty.MapVal(tagMap))
		}
	}

	return resourceFile, importFile, nil
}

// GenerateSecurityGroupBlocks はSecurity Groupリソースのresourceブロックとimportブロックを生成します。
func (g *HCLGenerator) GenerateSecurityGroupBlocks(sgs []ec2.SecurityGroup) (*hclwrite.File, *hclwrite.File, error) {
	resourceFile := hclwrite.NewEmptyFile()
	importFile := hclwrite.NewEmptyFile()
	resourceBody := resourceFile.Body()
	importBody := importFile.Body()

	for i, sg := range sgs {
		resourceName := fmt.Sprintf("sg_%d", i)
		resourceType := "aws_security_group"

		// importブロックの生成
		importBlock := importBody.AppendNewBlock("import", nil)
		importBody.AppendNewline()
		importBlock.Body().SetAttributeValue("to", cty.StringVal(resourceType+"."+resourceName))
		importBlock.Body().SetAttributeValue("id", cty.StringVal(sg.ID))

		// resourceブロックの生成
		sgBlock := resourceBody.AppendNewBlock("resource", []string{resourceType, resourceName})
		resourceBody.AppendNewline()
		sgBlock.Body().SetAttributeValue("name", cty.StringVal(sg.Name))
		sgBlock.Body().SetAttributeValue("description", cty.StringVal(sg.Description))

		if len(sg.Tags) > 0 {
			tagMap := make(map[string]cty.Value)
			for k, v := range sg.Tags {
				tagMap[k] = cty.StringVal(v)
			}
			sgBlock.Body().SetAttributeValue("tags", cty.MapVal(tagMap))
		}
	}

	return resourceFile, importFile, nil
}

// GenerateS3BucketBlocks はS3バケットリソースのresourceブロックとimportブロックを生成します。
func (g *HCLGenerator) GenerateS3BucketBlocks(buckets []s3.Bucket) (*hclwrite.File, *hclwrite.File, error) {
	resourceFile := hclwrite.NewEmptyFile()
	importFile := hclwrite.NewEmptyFile()
	resourceBody := resourceFile.Body()
	importBody := importFile.Body()

	for _, bucket := range buckets {
		resourceName := bucket.Name
		resourceType := "aws_s3_bucket"

		// importブロックの生成
		importBlock := importBody.AppendNewBlock("import", nil)
		importBody.AppendNewline()
		importBlock.Body().SetAttributeValue("to", cty.StringVal(resourceType+"."+resourceName))
		importBlock.Body().SetAttributeValue("id", cty.StringVal(bucket.Name))

		// resourceブロックの生成
		bucketBlock := resourceBody.AppendNewBlock("resource", []string{resourceType, resourceName})
		resourceBody.AppendNewline()
		bucketBlock.Body().SetAttributeValue("bucket", cty.StringVal(bucket.Name))
	}

	return resourceFile, importFile, nil
}

// GenerateEcsBlocks はECSリソースのresourceブロックとimportブロックを生成します。
func (g *HCLGenerator) GenerateEcsBlocks(clusters []ecs.Cluster) (*hclwrite.File, *hclwrite.File, error) {
	resourceFile := hclwrite.NewEmptyFile()
	importFile := hclwrite.NewEmptyFile()
	resourceBody := resourceFile.Body()
	importBody := importFile.Body()

	for _, cluster := range clusters {
		// Cluster
		clusterResourceType := "aws_ecs_cluster"
		clusterResourceName := cluster.Name
		g.appendImportBlock(importBody, clusterResourceType+"."+clusterResourceName, cluster.Name)
		clusterBlock := g.appendResourceBlock(resourceBody, clusterResourceType, clusterResourceName)
		clusterBlock.Body().SetAttributeValue("name", cty.StringVal(cluster.Name))
		if len(cluster.Tags) > 0 {
			tagMap := make(map[string]cty.Value)
			for k, v := range cluster.Tags {
				tagMap[k] = cty.StringVal(v)
			}
			clusterBlock.Body().SetAttributeValue("tags", cty.MapVal(tagMap))
		}

		// Services
		for _, service := range cluster.Services {
			serviceResourceType := "aws_ecs_service"
			serviceResourceName := service.Name
			importId := fmt.Sprintf("%s/%s", cluster.Name, service.Name)
			g.appendImportBlock(importBody, serviceResourceType+"."+serviceResourceName, importId)
			serviceBlock := g.appendResourceBlock(resourceBody, serviceResourceType, serviceResourceName)
			serviceBlock.Body().SetAttributeValue("name", cty.StringVal(service.Name))
			serviceBlock.Body().SetAttributeValue("task_definition", cty.StringVal(service.TaskDefinitionArn))
			serviceBlock.Body().SetAttributeValue("desired_count", cty.NumberIntVal(int64(service.DesiredCount)))
			serviceBlock.Body().SetAttributeValue("enable_ecs_managed_tags", cty.BoolVal(service.EnableEcsManagedTags))
			serviceBlock.Body().SetAttributeValue("enable_execute_command", cty.BoolVal(service.EnableExecuteCommand))
			serviceBlock.Body().SetAttributeValue("health_check_grace_period_seconds", cty.NumberIntVal(int64(service.HealthCheckGracePeriodSeconds)))
			serviceBlock.Body().SetAttributeValue("propagate_tags", cty.StringVal(service.PropagateTags))
			serviceBlock.Body().SetAttributeValue("platform_version", cty.StringVal(service.PlatformVersion))
			serviceBlock.Body().SetAttributeValue("scheduling_strategy", cty.StringVal(service.SchedulingStrategy))

			// Deployment Circuit Breaker
			if service.DeploymentCircuitBreaker != nil {
				breakerBlock := serviceBlock.Body().AppendNewBlock("deployment_circuit_breaker", nil)
				breakerBlock.Body().SetAttributeValue("enable", cty.BoolVal(service.DeploymentCircuitBreaker.Enable))
				breakerBlock.Body().SetAttributeValue("rollback", cty.BoolVal(service.DeploymentCircuitBreaker.Rollback))
			}

			// Network Configuration
			if service.NetworkConfiguration != nil {
				netBlock := serviceBlock.Body().AppendNewBlock("network_configuration", nil)
				var subnetVals []cty.Value
				for _, subnet := range service.NetworkConfiguration.Subnets {
					subnetVals = append(subnetVals, cty.StringVal(subnet))
				}
				netBlock.Body().SetAttributeValue("subnets", cty.ListVal(subnetVals)) // Use ListVal, TF can handle it.
				if len(service.NetworkConfiguration.SecurityGroups) > 0 {
					var sgVals []cty.Value
					for _, sg := range service.NetworkConfiguration.SecurityGroups {
						sgVals = append(sgVals, cty.StringVal(sg))
					}
					netBlock.Body().SetAttributeValue("security_groups", cty.ListVal(sgVals))
				}
				netBlock.Body().SetAttributeValue("assign_public_ip", cty.BoolVal(service.NetworkConfiguration.AssignPublicIp))
			}

			// Tags
			if len(service.Tags) > 0 {
				tagMap := make(map[string]cty.Value)
				for k, v := range service.Tags {
					tagMap[k] = cty.StringVal(v)
				}
				serviceBlock.Body().SetAttributeValue("tags", cty.MapVal(tagMap))
			}

			// Load Balancers
			if len(service.LoadBalancers) > 0 {
				for _, lb := range service.LoadBalancers {
					lbBlock := serviceBlock.Body().AppendNewBlock("load_balancer", nil)
					lbBlock.Body().SetAttributeValue("target_group_arn", cty.StringVal(lb.TargetGroupArn))
					lbBlock.Body().SetAttributeValue("container_name", cty.StringVal(lb.ContainerName))
					lbBlock.Body().SetAttributeValue("container_port", cty.NumberIntVal(int64(lb.ContainerPort)))
				}
			}

			// Cluster Reference
			clusterRefParts := strings.SplitN("aws_ecs_cluster."+clusterResourceName, ".", 2)
			clusterTraversal := hcl.Traversal{
				hcl.TraverseRoot{Name: clusterRefParts[0]},
				hcl.TraverseAttr{Name: clusterRefParts[1]},
				hcl.TraverseAttr{Name: "arn"},
			}
			serviceBlock.Body().SetAttributeRaw("cluster", hclwrite.TokensForTraversal(clusterTraversal))
		}
	}

	return resourceFile, importFile, nil
}

// GenerateElbBlocks はELBv2リソースのresourceブロックとimportブロックを生成します。
func (g *HCLGenerator) GenerateElbBlocks(lbs []*elbv2.LoadBalancer) (*hclwrite.File, *hclwrite.File, error) {
	resourceFile := hclwrite.NewEmptyFile()
	importFile := hclwrite.NewEmptyFile()
	resourceBody := resourceFile.Body()
	importBody := importFile.Body()

	for _, lb := range lbs {
		// Target Groups
		tgRefs := make(map[string]hcl.Traversal)
		for _, tg := range lb.TargetGroups {
			tgResourceType := "aws_lb_target_group"
			tgResourceName := tg.Name
			g.appendImportBlock(importBody, tgResourceType+"."+tgResourceName, tg.Arn)
			tgBlock := g.appendResourceBlock(resourceBody, tgResourceType, tgResourceName)
			tgBlock.Body().SetAttributeValue("name", cty.StringVal(tg.Name))
			tgBlock.Body().SetAttributeValue("port", cty.NumberIntVal(int64(tg.Port)))
			tgBlock.Body().SetAttributeValue("protocol", cty.StringVal(string(tg.Protocol)))
			tgBlock.Body().SetAttributeValue("vpc_id", cty.StringVal(tg.VpcId))
			tgBlock.Body().SetAttributeValue("target_type", cty.StringVal(tg.TargetType))

			if tg.HealthCheck != nil {
				hcBlock := tgBlock.Body().AppendNewBlock("health_check", nil)
				hcBlock.Body().SetAttributeValue("enabled", cty.BoolVal(tg.HealthCheck.Enabled))
				hcBlock.Body().SetAttributeValue("path", cty.StringVal(tg.HealthCheck.Path))
				hcBlock.Body().SetAttributeValue("port", cty.StringVal(tg.HealthCheck.Port))
				hcBlock.Body().SetAttributeValue("protocol", cty.StringVal(string(tg.HealthCheck.Protocol)))
				hcBlock.Body().SetAttributeValue("interval", cty.NumberIntVal(int64(tg.HealthCheck.Interval)))
				hcBlock.Body().SetAttributeValue("timeout", cty.NumberIntVal(int64(tg.HealthCheck.Timeout)))
				hcBlock.Body().SetAttributeValue("healthy_threshold", cty.NumberIntVal(int64(tg.HealthCheck.HealthyThreshold)))
				hcBlock.Body().SetAttributeValue("unhealthy_threshold", cty.NumberIntVal(int64(tg.HealthCheck.UnhealthyThreshold)))
				hcBlock.Body().SetAttributeValue("matcher", cty.StringVal(tg.HealthCheck.Matcher))
			}

			tgRefParts := strings.SplitN(tgResourceType+"."+tgResourceName, ".", 2)
			tgRefs[tg.Arn] = hcl.Traversal{
				hcl.TraverseRoot{Name: tgRefParts[0]},
				hcl.TraverseAttr{Name: tgRefParts[1]},
			}
		}

		// Load Balancer
		lbResourceType := "aws_lb"
		lbResourceName := lb.Name
		g.appendImportBlock(importBody, lbResourceType+"."+lbResourceName, lb.Arn)
		lbBlock := g.appendResourceBlock(resourceBody, lbResourceType, lbResourceName)
		lbBlock.Body().SetAttributeValue("name", cty.StringVal(lb.Name))
		lbBlock.Body().SetAttributeValue("load_balancer_type", cty.StringVal(string(lb.Type)))
		subnetVals := []cty.Value{}
		for _, subnet := range lb.Subnets {
			subnetVals = append(subnetVals, cty.StringVal(subnet))
		}
		lbBlock.Body().SetAttributeValue("subnets", cty.ListVal(subnetVals))

		lbRefParts := strings.SplitN(lbResourceType+"."+lbResourceName, ".", 2)
		lbTraversal := hcl.Traversal{
			hcl.TraverseRoot{Name: lbRefParts[0]},
			hcl.TraverseAttr{Name: lbRefParts[1]},
		}

		// Listeners
		listenerRefs := make(map[string]hcl.Traversal)
		for _, listener := range lb.Listeners {
			listenerResourceType := "aws_lb_listener"
			listenerResourceName := fmt.Sprintf("%s_%d", lb.Name, listener.Port)
			g.appendImportBlock(importBody, listenerResourceType+"."+listenerResourceName, listener.Arn)
			listenerBlock := g.appendResourceBlock(resourceBody, listenerResourceType, listenerResourceName)

			listenerRefParts := strings.SplitN(listenerResourceType+"."+listenerResourceName, ".", 2)
			listenerTraversal := hcl.Traversal{
				hcl.TraverseRoot{Name: listenerRefParts[0]},
				hcl.TraverseAttr{Name: listenerRefParts[1]},
			}
			listenerRefs[listener.Arn] = listenerTraversal

			listenerBlock.Body().SetAttributeRaw("load_balancer_arn", hclwrite.TokensForTraversal(append(lbTraversal, hcl.TraverseAttr{Name: "arn"})))
			listenerBlock.Body().SetAttributeValue("port", cty.NumberIntVal(int64(listener.Port)))
			listenerBlock.Body().SetAttributeValue("protocol", cty.StringVal(string(listener.Protocol)))
			if listener.CertificateArn != nil {
				listenerBlock.Body().SetAttributeValue("certificate_arn", cty.StringVal(*listener.CertificateArn))
			}

			// Default Action
			for _, da := range listener.DefaultActions {
				daBlock := listenerBlock.Body().AppendNewBlock("default_action", nil)
				daBlock.Body().SetAttributeValue("type", cty.StringVal(string(da.Type)))

				if da.Type == "forward" && da.Forward != nil {
					forwardBlock := daBlock.Body().AppendNewBlock("forward", nil)
					for _, tg := range da.Forward.TargetGroups {
						tgBlock := forwardBlock.Body().AppendNewBlock("target_group", nil)
						tgTraversal := tgRefs[tg.Arn]
						tgBlock.Body().SetAttributeRaw("arn", hclwrite.TokensForTraversal(append(tgTraversal, hcl.TraverseAttr{Name: "arn"})))
						if tg.Weight != nil {
							tgBlock.Body().SetAttributeValue("weight", cty.NumberIntVal(*tg.Weight))
						}
					}
					if da.Forward.Stickiness != nil {
						stickinessBlock := forwardBlock.Body().AppendNewBlock("stickiness", nil)
						stickinessBlock.Body().SetAttributeValue("enabled", cty.BoolVal(da.Forward.Stickiness.Enabled))
						stickinessBlock.Body().SetAttributeValue("duration", cty.NumberIntVal(int64(da.Forward.Stickiness.Duration)))
					}
				} else if da.Type == "redirect" && da.Redirect != nil {
					redirectBlock := daBlock.Body().AppendNewBlock("redirect", nil)
					redirectBlock.Body().SetAttributeValue("port", cty.StringVal(da.Redirect.Port))
					redirectBlock.Body().SetAttributeValue("protocol", cty.StringVal(da.Redirect.Protocol))
					redirectBlock.Body().SetAttributeValue("status_code", cty.StringVal(da.Redirect.StatusCode))
				}
			}

			// Listener Rules
			for _, rule := range listener.Rules {
				ruleResourceType := "aws_lb_listener_rule"
				// A bit of a hack to create a unique name. A better approach might be needed.
				ruleResourceName := fmt.Sprintf("%s_rule_%s", listenerResourceName, rule.Priority)
				g.appendImportBlock(importBody, ruleResourceType+"."+ruleResourceName, rule.Arn)
				ruleBlock := g.appendResourceBlock(resourceBody, ruleResourceType, ruleResourceName)
				listenerRef := listenerRefs[listener.Arn]
				ruleBlock.Body().SetAttributeRaw("listener_arn", hclwrite.TokensForTraversal(append(listenerRef, hcl.TraverseAttr{Name: "arn"})))
				ruleBlock.Body().SetAttributeValue("priority", cty.StringVal(rule.Priority))

				// Actions
				for _, action := range rule.Actions {
					actionBlock := ruleBlock.Body().AppendNewBlock("action", nil)
					actionBlock.Body().SetAttributeValue("type", cty.StringVal(string(action.Type)))
					if action.Forward != nil {
						tgTraversal := tgRefs[action.Forward.TargetGroupArn]
						actionBlock.Body().SetAttributeRaw("target_group_arn", hclwrite.TokensForTraversal(append(tgTraversal, hcl.TraverseAttr{Name: "arn"})))
					}
				}

				// Conditions
				for _, cond := range rule.Conditions {
					condBlock := ruleBlock.Body().AppendNewBlock("condition", nil)
					if cond.Field == "host-header" {
						hostHeaderBlock := condBlock.Body().AppendNewBlock("host_header", nil)
						var vals []cty.Value
						for _, v := range cond.Values {
							vals = append(vals, cty.StringVal(v))
						}
						hostHeaderBlock.Body().SetAttributeValue("values", cty.ListVal(vals))
					}
					if cond.Field == "path-pattern" {
						pathPatternBlock := condBlock.Body().AppendNewBlock("path_pattern", nil)
						var vals []cty.Value
						for _, v := range cond.Values {
							vals = append(vals, cty.StringVal(v))
						}
						pathPatternBlock.Body().SetAttributeValue("values", cty.ListVal(vals))
					}
				}
			}
		}
	}

	return resourceFile, importFile, nil
}

func (g *HCLGenerator) appendImportBlock(body *hclwrite.Body, to, id string) {
	importBlock := body.AppendNewBlock("import", nil)
	// "to" is a resource address, not a string. e.g., aws_ecs_cluster.my_cluster
	parts := strings.SplitN(to, ".", 2)
	traversal := hcl.Traversal{
		hcl.TraverseRoot{Name: parts[0]},
		hcl.TraverseAttr{Name: parts[1]},
	}
	importBlock.Body().SetAttributeRaw("to", hclwrite.TokensForTraversal(traversal))
	importBlock.Body().SetAttributeValue("id", cty.StringVal(id))
	body.AppendNewline()
}

func (g *HCLGenerator) appendResourceBlock(body *hclwrite.Body, resourceType, resourceName string) *hclwrite.Block {
	block := body.AppendNewBlock("resource", []string{resourceType, resourceName})
	body.AppendNewline()
	return block
}

// GenerateIamBlocks はIAMリソースのresourceブロックとimportブロックを生成します。
func (g *HCLGenerator) GenerateIamBlocks(policies []iam.Policy, roles []iam.Role) (*hclwrite.File, *hclwrite.File, error) {
	resourceFile := hclwrite.NewEmptyFile()
	importFile := hclwrite.NewEmptyFile()
	resourceBody := resourceFile.Body()
	importBody := importFile.Body()

	policyData := make(map[string]struct {
		Name      string
		Traversal hcl.Traversal
	})

	// IAM Policies
	for _, p := range policies {
		resourceType := "aws_iam_policy"
		resourceName := p.Name
		g.appendImportBlock(importBody, resourceType+"."+resourceName, p.Arn)
		policyBlock := g.appendResourceBlock(resourceBody, resourceType, resourceName)
		policyBlock.Body().SetAttributeValue("name", cty.StringVal(p.Name))
		policyBlock.Body().SetAttributeValue("policy", cty.StringVal(p.PolicyDocument))

		refParts := strings.SplitN(resourceType+"."+resourceName, ".", 2)
		policyData[p.Arn] = struct {
			Name      string
			Traversal hcl.Traversal
		}{
			Name: resourceName,
			Traversal: hcl.Traversal{
				hcl.TraverseRoot{Name: refParts[0]},
				hcl.TraverseAttr{Name: refParts[1]},
			},
		}
	}

	// IAM Roles
	for _, r := range roles {
		roleResourceType := "aws_iam_role"
		roleResourceName := r.Name
		g.appendImportBlock(importBody, roleResourceType+"."+roleResourceName, r.Name) // IAM Role ID is its name
		roleBlock := g.appendResourceBlock(resourceBody, roleResourceType, roleResourceName)
		roleBlock.Body().SetAttributeValue("name", cty.StringVal(r.Name))
		roleBlock.Body().SetAttributeValue("assume_role_policy", cty.StringVal(r.AssumeRolePolicy))

		// IAM Role Policy Attachments
		for _, policyArn := range r.AttachedPolicyArns {
			if data, ok := policyData[policyArn]; ok {
				attachmentResourceName := fmt.Sprintf("%s_%s_attachment", r.Name, data.Name)
				attachmentResourceType := "aws_iam_role_policy_attachment"
				// No import block for attachments, they are managed with the role.
				attachmentBlock := g.appendResourceBlock(resourceBody, attachmentResourceType, attachmentResourceName)

				roleRefParts := strings.SplitN(roleResourceType+"."+roleResourceName, ".", 2)
				roleTraversal := hcl.Traversal{
					hcl.TraverseRoot{Name: roleRefParts[0]},
					hcl.TraverseAttr{Name: roleRefParts[1]},
				}

				policyTraversal := data.Traversal

				attachmentBlock.Body().SetAttributeRaw("role", hclwrite.TokensForTraversal(append(roleTraversal, hcl.TraverseAttr{Name: "name"})))
				attachmentBlock.Body().SetAttributeRaw("policy_arn", hclwrite.TokensForTraversal(append(policyTraversal, hcl.TraverseAttr{Name: "arn"})))
			}
		}
	}

	return resourceFile, importFile, nil
} 
