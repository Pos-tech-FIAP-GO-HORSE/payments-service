provider "aws" {
  region = "us-east-1"
}

#Lambda role
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

#Lambda policy
resource "aws_iam_policy" "lambda_all_policies" {
  name = "lambda-all-policies"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      #Permissões para SNS
      {
        Action   = "sns:*"
        Effect   = "Allow"
        Resource = aws_sns_topic.orders-service-events-order-created.arn
      },
      {
        Action   = "lambda:InvokeFunction"
        Effect   = "Allow"
        Resource = aws_lambda_function.payments-service-lambda.arn
      },
      #Permissões para S3
      {
        Action   = "s3:GetObject"
        Effect   = "Allow"
        Resource = "arn:aws:s3:::payments-service/*"
      },
      #Permissões para o usuário visualizar as roles e políticas
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
      #Permissões para listar usuários
      {
        Action   = [
          "iam:ListUsers",
        ]
        Effect   = "Allow"
        Resource = "*"
      },
      #Permissões para escrever em um tópico
      {
        Action   = [
          "sns:Publish"
        ]
        Effect   = "Allow"
        Resource = aws_sns_topic.payments-service-events-payment-created.arn
      }
    ]
  })
}

#Define the lambda function
resource "aws_lambda_function" "payments-service-lambda" {
  function_name = "payments-service-lambda"
  handler       = "main"
  runtime       = "provided.al2"
  s3_bucket     = "payments-service"
  s3_key        = "function.zip"
  role          = aws_iam_role.lambda_exec_role_payments_service.arn

  environment {
    variables = {
      DB_URI             = "mongodb+srv://admin:admin123@payment-cluster.sl4mh.mongodb.net/pos-tech-fiap?retryWrites=true&w=majority"
      DB_NAME            = "pos-tech-fiap"
      DB_COLLECTION_NAME = "payments"
      TOKEN_MERCADO_PAGO = "TEST-2373946154784631-101516-50ff7f4dcdff3aec43372568c77990e3-175794680"
      SNS_TOPIC_ARN      = "arn:aws:sns:us-east-1:537124948968:payments-service-events-payment-created"
    }
  }
}

#Define the sns topic for receive create payment events
resource "aws_sns_topic" "orders-service-events-order-created" {
  name = "orders-service-events-order-created"
}

#Define the sns topic for send payment created events
resource "aws_sns_topic" "payments-service-events-payment-created" {
  name = "payments-service-events-payment-created"
}

#Define the signature of the sns topic in the lambda to be invoked
resource "aws_sns_topic_subscription" "my_sns_subscription" {
  topic_arn = aws_sns_topic.orders-service-events-order-created.arn
  protocol  = "lambda"
  endpoint  = aws_lambda_function.payments-service-lambda.arn
}

# Associar a política ao usuário criado via AWS CLI (usando o ARN do usuário)
resource "aws_iam_policy_attachment" "attach_policy_to_user" {
  name       = "attach-github-payments-service-policy"
  policy_arn = aws_iam_policy.lambda_all_policies.arn
  users      = ["github-payments-service"]
}

#Anexar a política IAMReadPermissions ao usuário
 resource "aws_iam_policy_attachment" "attach_iam_read_policy" {
   name       = "attach-iam-read-policy"
   policy_arn = aws_iam_policy.lambda_all_policies.arn
   users      = ["github-payments-service"]
 }

# Anexar a política ao usuário
resource "aws_iam_user_policy_attachment" "sns_attach_policy" {
  user       = aws_iam_user.github_payments_service.name
  policy_arn = aws_iam_policy.lambda_all_policies.arn
}

resource "aws_iam_user" "github_payments_service" {
  name = "github-payments-service"
}
