package hcl

import (
	"fmt"

	"github.com/Haussmann000/tfimport/internal/aws/ec2"
	"github.com/Haussmann000/tfimport/internal/aws/s3"
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
