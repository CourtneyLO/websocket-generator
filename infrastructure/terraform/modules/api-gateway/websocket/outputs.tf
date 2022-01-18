resource "aws_cloudformation_stack" "outputs" {
  name = "${var.PROJECT_NAME}-${var.ENVIRONMENT}-terraform-outputs"

  template_body = <<STACK
    {
      "Resources": {
        "IAMRoleArnAuthorizerLambda": {
          "Type": "AWS::SSM::Parameter",
          "Properties": {
            "Name": "${var.PROJECT_NAME}-${var.ENVIRONMENT}-iam-role-arn-authorization-lambda",
            "Type": "String",
            "Value": "${aws_iam_role.authorization_lambda.arn}"
          }
        },
        "IAMRoleArnConnectLambda": {
          "Type": "AWS::SSM::Parameter",
          "Properties": {
            "Name": "${var.PROJECT_NAME}-${var.ENVIRONMENT}-iam-role-arn-connect-lambda",
            "Type": "String",
            "Value": "${aws_iam_role.connect_lambda.arn}"
          }
        },
        "IAMRoleArnDisconnectLambda": {
          "Type": "AWS::SSM::Parameter",
          "Properties": {
            "Name": "${var.PROJECT_NAME}-${var.ENVIRONMENT}-iam-role-arn-disconnect-lambda",
            "Type": "String",
            "Value": "${aws_iam_role.disconnect_lambda.arn}"
          }
        },
        "IAMRoleArnDefaultLambda": {
          "Type": "AWS::SSM::Parameter",
          "Properties": {
            "Name": "${var.PROJECT_NAME}-${var.ENVIRONMENT}-iam-role-arn-default-lambda",
            "Type": "String",
            "Value": "${aws_iam_role.default_lambda.arn}"
          }
        },
        "DynamoDBTableNameWebsocketConnections": {
          "Type": "AWS::SSM::Parameter",
          "Properties": {
            "Name": "${var.PROJECT_NAME}-${var.ENVIRONMENT}-dynamodb-table-name-websocket-connections",
            "Type": "String",
            "Value": "${aws_dynamodb_table.dynamodb_websocket_manager.id}"
          }
        }
      },
      "Outputs": {
        "IAMRoleArnAuthorizerLambda": {
          "Value": "${aws_iam_role.authorization_lambda.arn}"
        },
        "IAMRoleArnConnectLambda": {
          "Value": "${aws_iam_role.connect_lambda.arn}"
        },
        "IAMRoleArnDisconnectLambda": {
          "Value": "${aws_iam_role.disconnect_lambda.arn}"
        },
        "IAMRoleArnDefaultLambda": {
          "Value": "${aws_iam_role.default_lambda.arn}"
        },
        "DynamoDBTableNameWebsocketConnections": {
          "Value": "${aws_dynamodb_table.dynamodb_websocket_manager.id}"
        }
      }
    }
  STACK
}
