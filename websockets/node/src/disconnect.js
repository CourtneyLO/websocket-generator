const DynamoDB = require('aws-sdk/clients/dynamodb');

const { WEBSOCKET_MANAGER_TABLE_NAME } = process.env;

exports.handler = async function(event, context, callback) {
  console.log('Disconnect Handler');

  const db = new DynamoDB.DocumentClient();
  const deleteParams = {
    TableName: WEBSOCKET_MANAGER_TABLE_NAME,
    Key: {
      connectionId: event.requestContext.connectionId,
    }
  };

  try {
    await db.delete(deleteParams).promise();
    return {
      statusCode: 200,
      body: "Disconnected"
    };
  } catch (error) {
    console.error('Error', error);
    return {
      statusCode: 500,
      body: "Failed to disconnect: " + JSON.stringify(error),
    };
  }
};
