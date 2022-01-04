# WebSocket Generator - WIP

This script will generate a fully functional WebSocket using AWS Gateway Api and DynamoDB.

The WebSocket routes will be as follows:

- $connect: Makes the WebSocket connection and saves the connectionID, timestamp and deviceTypeto a dynamoDB table. DeviceType helps to distinguish between different device connectivity
- $disconnect: Disconnects the WebSocket and removes the connectionID from the dynamoDB table.
- $default: This will contain the `__PING__` `__PONG__` functionality which ensures the connection can kept alive throughout it's use.

It will be built using Terraform and Serverless with choice of one of the following 6 languages:

- Node
- Typescript
- Python
- Golang
- Java
- Ruby

All Terraform and Serveless code will be added to you project in the directory of your choice. Once the coded is added you can alter or modified the code as you please. You can choose to continue using the build in scripts (see below)or use the Terraform and Serverless commands independently, it is up you. This is a quick first step to creating your own websocket.

## Prerequists

- You will need to have your aws account ID
- For a better user interface use a full screen terminal tab

## Config Questions

|   | Question                         | Description  | Default  | Examples |
|---|----------------------------------|--------------|----------|----------|
| 1 | Select your preferred language   |  The language you wish for your WebSocket lambdas to be written            |   Node      | Node, Typescript, Python, Golang, Java, Ruby   |
| 2  | Environment Name | The environment you wish to deploy your code | Development | Development, Staging, Production
| 3  | Infrastructure Code Location Path | The path should be only the path from the current directory. The current Directory path will be joined with the given path. Make sure to be in the directory (project) you wish to add the WebSocket infrastructure code too. | infrastructure | infrastructure |
| 4  | WebSocket Code Location Path | The path should be only the path from the current directory. The current Directory path will be joined with the given path. Make sure to be in the directory (project) you wish to add the WebSocket serverless code too. | websocket | websocket
| 5  | Terraform Resource Name | The prefix you wish to add to all your resource names. Resource name format will look as \<terraform-resource-name\>-\<environment\>-\<specifc-resource-name\> | your projects name | helloworld-development-websocket-manager|
| 6  | AWS Region | The AWS region where you want your infrastructure to be located | eu-west-2 | us-east-1 |
| 7  | AWS Account ID | The AWS account ID that you wish to use for your infrastructure | This is required, there is no default | 123456789000 |
| 8 | Authorization Key Name | The name of the authorizer query parameter on your WebSocket URL  | Authorization | ?Authorization=1234 |
