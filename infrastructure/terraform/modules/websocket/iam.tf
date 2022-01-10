# Default assume-role-policy for Lambda
data "aws_iam_policy_document" "lambda_assume_role_policy" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

# IAM role for the Authorizer lambda
resource "aws_iam_role" "authorization_lambda" {
  name               = "${var.PROJECT_NAME}-${var.ENVIRONMENT}-AuthorizerLambda"
  assume_role_policy = data.aws_iam_policy_document.lambda_assume_role_policy.json
}

resource "aws_iam_role_policy" "authorization_lambda" {
  name   = "${var.PROJECT_NAME}-${var.ENVIRONMENT}-AuthorizerLambda"
  role   = aws_iam_role.authorization_lambda.id
  policy = data.aws_iam_policy_document.authorization_lambda.json
}

data "aws_iam_policy_document" "authorization_lambda" {
  statement {
    actions   = ["secretsmanager:GetSecretValue"]
    resources = ["*"]
  }

  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]
    resources = [
      "arn:aws:logs:${var.AWS_REGION}:${var.AWS_ACCOUNT_ID}:log-group:*:*"
    ]
  }
}

# IAM role for the Connect lambda
resource "aws_iam_role" "connect_lambda" {
  name               = "${var.PROJECT_NAME}-${var.ENVIRONMENT}-ConnectLambda"
  assume_role_policy = data.aws_iam_policy_document.lambda_assume_role_policy.json
}

resource "aws_iam_role_policy" "connect_lambda" {
  name   = "${var.PROJECT_NAME}-${var.ENVIRONMENT}-ConnectLambda"
  role   = aws_iam_role.connect_lambda.id
  policy = data.aws_iam_policy_document.connect_lambda.json
}

data "aws_iam_policy_document" "connect_lambda" {
  statement {
    actions = [
      "dynamodb:PutItem",
      "dynamodb:Scan",
      "dynamodb:DeleteItem"
    ]
    resources = [aws_dynamodb_table.dynamodb_websocket_manager.arn]
  }

  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]
    resources = ["arn:aws:logs:${var.AWS_REGION}:${var.AWS_ACCOUNT_ID}:log-group:*:*"]
  }
}

# IAM role for the Disconnect lambda
resource "aws_iam_role" "disconnect_lambda" {
  name               = "${var.PROJECT_NAME}-${var.ENVIRONMENT}-DisconnectLambda"
  assume_role_policy = data.aws_iam_policy_document.lambda_assume_role_policy.json
}

resource "aws_iam_role_policy" "disconnect_lambda" {
  name   = "${var.PROJECT_NAME}-${var.ENVIRONMENT}-DisconnectLambda"
  role   = aws_iam_role.disconnect_lambda.id
  policy = data.aws_iam_policy_document.disconnect_lambda.json
}

data "aws_iam_policy_document" "disconnect_lambda" {
  statement {
    actions   = ["dynamodb:DeleteItem"]
    resources = [aws_dynamodb_table.dynamodb_websocket_manager.arn]
  }

  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]
    resources = ["arn:aws:logs:${var.AWS_REGION}:${var.AWS_ACCOUNT_ID}:log-group:*:*"]
  }
}

# IAM role for the Default lambda
resource "aws_iam_role" "default_lambda" {
  name               = "${var.PROJECT_NAME}-${var.ENVIRONMENT}-DefaultLambda"
  assume_role_policy = data.aws_iam_policy_document.lambda_assume_role_policy.json
}

resource "aws_iam_role_policy" "default_lambda" {
  name   = "${var.PROJECT_NAME}-${var.ENVIRONMENT}-DefaultLambda"
  role   = aws_iam_role.default_lambda.id
  policy = data.aws_iam_policy_document.default_lambda.json
}

data "aws_iam_policy_document" "default_lambda" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]
    resources = ["arn:aws:logs:${var.AWS_REGION}:${var.AWS_ACCOUNT_ID}:log-group:*:*"]
  }
}
