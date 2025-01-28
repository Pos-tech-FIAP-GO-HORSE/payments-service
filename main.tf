provider "aws" {
  region = "us-east-1"
}

resource "aws_iam_role" "lambda_exec_role_payments_service" {
  name = "lambda_exec_role_payments_service"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action    = "sts:AssumeRole"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
        Effect    = "Allow"
      },
    ]
  })
}

resource "aws_iam_policy" "lambda_s3_access_policy" {
  name        = "lambda-s3-access-policy"
  description = "Política para permitir que Lambda acesse objetos no S3"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action   = "s3:GetObject"
        Effect   = "Allow"
        Resource = "arn:aws:s3:::payments-service/*"
      }
    ]
  })
}

# Criação de uma política que permite o acesso a GetRole e GetPolicy
resource "aws_iam_policy" "iam_read_permissions" {
  name        = "IAMReadPermissions"
  description = "Permissão para consultar roles e policies do IAM"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action   = [
          "iam:GetRole",
          "iam:GetPolicy",
          "iam:ListRoles",
          "iam:ListPolicies"
        ]
        Effect   = "Allow"
        Resource = "*"
      },
    ]
  })
}

# Anexar a política IAMReadPermissions ao usuário
resource "aws_iam_policy_attachment" "attach_iam_read_policy" {
  name       = "attach-iam-read-policy"
  policy_arn = aws_iam_policy.iam_read_permissions.arn
  users      = ["github-payments-service"]
}

resource "aws_iam_policy" "list_users_policy" {
  name        = "ListUsersPolicy"
  description = "Permissão para listar os usuários no IAM"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action   = "iam:ListUsers"
        Effect   = "Allow"
        Resource = "*"
      },
    ]
  })
}

resource "aws_iam_policy_attachment" "attach_list_users_policy" {
  name       = "attach-list-users-policy"
  policy_arn = aws_iam_policy.list_users_policy.arn
  users      = ["github-payments-service"]
}

resource "aws_iam_role_policy_attachment" "lambda_s3_policy_attachment" {
  policy_arn = aws_iam_policy.lambda_s3_access_policy.arn
  role       = aws_iam_role.lambda_exec_role_payments_service.name
}

resource "aws_lambda_function" "payments-service-lambda" {
  function_name = "payments-service-lambda"
  handler       = "main"
  runtime       = "provided.al2"
  s3_bucket     = "payments-service"
  s3_key        = "function.zip"
  role          = aws_iam_role.lambda_exec_role_payments_service.arn
}

resource "aws_sns_topic" "payments-service-sns-topic" {
  name = "payments-service-sns-topic"
}

resource "aws_lambda_permission" "allow_sns_invocation" {
  statement_id  = "AllowExecutionFromSNS"
  action        = "lambda:InvokeFunction"
  principal     = "sns.amazonaws.com"
  function_name = aws_lambda_function.payments-service-lambda.function_name
  source_arn    = aws_sns_topic.payments-service-sns-topic.arn
}

resource "aws_sns_topic_subscription" "my_sns_subscription" {
  topic_arn = aws_sns_topic.payments-service-sns-topic.arn
  protocol  = "lambda"
  endpoint  = aws_lambda_function.payments-service-lambda.arn
}