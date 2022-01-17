const DynamoDB = require('aws-sdk/clients/dynamodb');

const { WEBSOCKET_MANAGER_TABLE_NAME } = process.env;

exports.handler = async function(event, context, callback) {
  console.log('Disconnect Handler');
  const connectionId = event.requestContext.connectionId;

  const db = new DynamoDB.DocumentClient();
  const deleteParams = {
    TableName: WEBSOCKET_MANAGER_TABLE_NAME,
    Key: { connectionId }
  };

  try {
    await db.delete(deleteParams).promise();
    console.log(`Success: ConnectionId ${connectionId} was deleted from the database`);
  } catch (error) {
    console.error(`Error: ConnectionId ${connectionId} failed to be deleted from the database`, error);
  }
};
