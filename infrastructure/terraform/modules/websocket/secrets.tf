resource "aws_secretsmanager_secret" "authorization_secret_websocket_key" {
  name                    = "${var.PROJECT_NAME}-${var.ENVIRONMENT}-websocket-authorization-key"
  recovery_window_in_days = 0
}

resource "aws_secretsmanager_secret_version" "authorization_secret_websocket_version" {
  secret_id     = aws_secretsmanager_secret.authorization_secret_websocket_key.id
  secret_string = jsonencode(var.WEBSOCKET_AUTHORIZATION_SECRET_VALUE)
}
