resource "aws_security_group" "sg_0" {
  name        = "hrmc-stg-redis"
  description = "hrmc-stg redis traffic rule"
  tags = {
    IaC         = "terraform"
    Name        = "hrmc-stg-redis-sg"
    environment = "stg"
    service     = "hrmc"
  }
}

