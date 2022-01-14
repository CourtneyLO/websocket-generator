const ApiGatewayManagementApi = require('aws-sdk/clients/apigatewaymanagementapi');
const DynamoDB = require('aws-sdk/clients/dynamodb');

const { WEBSOCKET_MANAGER_TABLE_NAME, WEBSOCKET_URL } = process.env;

const deleteRowByConnectionId = async (connectionId, db) => {
  try {
    const deleteParams = {
      TableName: WEBSOCKET_MANAGER_TABLE_NAME,
      Key: { connectionId }
    };

    console.log(`Success: ConnectionId ${connectionId} has been removed from the database`);
    return await db.delete(deleteParams).promise();
  } catch (error) {
    console.error(`Error: ConnectionId ${connectionId} has failed to be removed from the database`);
    throw error;
  }
};

const postToConnectionById = async (connectionId, body, api, db) => {
  try {
    return await api.postToConnection({ ConnectionId: connectionId, Data: JSON.stringify(body) }).promise();
  } catch (error) {
    if (error.statusCode === 410) {
      console.log(`Info: Removing stale connectionId ${connectionId}`, error);
      await deleteRowByConnectionId(connectionId, db);
    }
    console.error(`Error: postToConnection request failed with connectionId ${connectionId}`, error);
    throw error;
  }
};

exports.handler = async function(event) {
  console.log('Default Handler');
  const body = event.body;

	if (body === '__ping__') {
    console.log('Success: __pong__');
    return { statusCode: 200, body: '__pong__' };
  }

  const db = new DynamoDB.DocumentClient();
  let connections;
  try {
    connections = await db.scan({ TableName: WEBSOCKET_MANAGER_TABLE_NAME }).promise();
  } catch (error) {
    console.error(`Error: Failed to scan ${WEBSOCKET_MANAGER_TABLE_NAME} for connection IDs`, error);
    return { statusCode: 500, body: error.stack };
  }

  try {
    const api = new ApiGatewayManagementApi({ endpoint: WEBSOCKET_URL });
    const postCalls = await Promise.all(connections.Items.map(async ({ connectionId }) => {
      return await postToConnectionById(connectionId, body, api, db);
    }));

    console.log('Success: Message sent!');
    return { statusCode: 200, body: 'Message sent!' };
  } catch (error) {
    console.error("Error: Message has not been sent", error);
    return { statusCode: 500, body: error.stack };
  }
};
