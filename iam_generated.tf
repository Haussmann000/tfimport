resource "aws_iam_policy" "hrmc-stg-next-allow-ssm-message" {
  name   = "hrmc-stg-next-allow-ssm-message"
  policy = "{\n    \"Statement\": [\n        {\n            \"Action\": [\n                \"ssmmessages:CreateControlChannel\",\n                \"ssmmessages:CreateDataChannel\",\n                \"ssmmessages:OpenControlChannel\",\n                \"ssmmessages:OpenDataChannel\"\n            ],\n            \"Effect\": \"Allow\",\n            \"Resource\": \"*\"\n        }\n    ],\n    \"Version\": \"2012-10-17\"\n}"
}

resource "aws_iam_role" "hrmc-stg-next-execution" {
  name               = "hrmc-stg-next-execution"
  assume_role_policy = "{\"Version\":\"2012-10-17\",\"Statement\":[{\"Sid\":\"\",\"Effect\":\"Allow\",\"Principal\":{\"Service\":\"ecs-tasks.amazonaws.com\"},\"Action\":\"sts:AssumeRole\"}]}"
}

resource "aws_iam_role" "hrmc-stg-next-task" {
  name               = "hrmc-stg-next-task"
  assume_role_policy = "{\"Version\":\"2012-10-17\",\"Statement\":[{\"Sid\":\"\",\"Effect\":\"Allow\",\"Principal\":{\"Service\":\"ecs-tasks.amazonaws.com\"},\"Action\":\"sts:AssumeRole\"}]}"
}

