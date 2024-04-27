package lib

const AWS_VPC_TMPL = `
resource "aws_vpc" "vpc_{{.Index}}" {
  cidr_block = "{{.CidrBlock}}"
  tags = {{parsetags .Tags}}
}
`
