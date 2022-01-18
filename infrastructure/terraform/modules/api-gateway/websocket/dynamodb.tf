resource "aws_dynamodb_table" "dynamodb_websocket_manager" {
  name            = "${var.PROJECT_NAME}-${var.ENVIRONMENT}-websocket-manager"
  billing_mode    = "PROVISIONED"
  read_capacity   = 5
  write_capacity  = 5
  hash_key        = "connectionId"

  attribute {
    name = "connectionId"
    type = "S"
  }
}
