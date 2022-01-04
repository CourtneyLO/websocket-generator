const DynamoDB = require('aws-sdk/clients/dynamodb');

exports.handler = async function(event, context, callback) {
  console.log('DISCONNECT HANDLER');

  const db = new DynamoDB.DocumentClient();
  const deleteParams = {
    TableName: process.env.WEBSOCKET_TABLE_NAME,
    Key: {
      connectionId: event.requestContext.connectionId,
    }
  };

  try {
    await db.delete(deleteParams).promise();
    console.log('DISCONNECTED');
    return {
      statusCode: 200,
      body: "Disconnected"
    };
  } catch (error) {
    console.error('ERROR', error);
    return {
      statusCode: 501,
      body: "Failed to disconnect: " + JSON.stringify(error),
    };
  }
};
