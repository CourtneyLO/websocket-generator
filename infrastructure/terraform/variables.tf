variable "ENVIRONMENT" {
  type = string
}

variable "PROJECT_NAME" {
  type = string
}

variable "AWS_REGION" {
  type = string
}

variable "AWS_ACCOUNT_ID" {
  type = string
}

variable "WEBSOCKET_AUTHORIZATION_SECRET_VALUE" {
  type      = string
  sensitive = true
}
