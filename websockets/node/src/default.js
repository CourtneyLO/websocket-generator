const ApiGatewayManagementApi = require('aws-sdk/clients/apigatewaymanagementapi');
const DynamoDB = require('aws-sdk/clients/dynamodb');

const { WEBSOCKET_MANAGER_TABLE_NAME, WEBSOCKET_URL } = process.env;
const ACTIVE_CONNECTION = 'Message sent to connection';
const STALE_CONNECTION = 'Deleted stale ConnectionId';

const deleteRowByConnectionId = async (connectionId, db) => {
  try {
    const deleteParams = {
      TableName: WEBSOCKET_MANAGER_TABLE_NAME,
      Key: { connectionId }
    };

    console.log(`Success: ConnectionId ${connectionId} has been removed from the database`);
    await db.delete(deleteParams).promise();
    return STALE_CONNECTION
  } catch (error) {
    console.error(`Error: ConnectionId ${connectionId} has failed to be removed from the database`);
    throw error;
  }
};

const postToConnectionById = async (connectionId, body, api, db) => {
  try {
    await api.postToConnection({ ConnectionId: connectionId, Data: JSON.stringify(body) }).promise();
    console.log(`Success: Message sent to connectionId ${connectionId}`);
    return ACTIVE_CONNECTION;
  } catch (error) {
    if (error.statusCode === 410) {
      console.log(`Info: Removing stale connectionId ${connectionId}`, error);
      return deleteRowByConnectionId(connectionId, db);
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
    console.log(`Success: ${connections.Count} connectionIds scanned in the database`);
  } catch (error) {
    console.error(`Error: Failed to scan ${WEBSOCKET_MANAGER_TABLE_NAME} for connection IDs`, error);
    return { statusCode: 500, body: error.stack };
  }

  try {
    const api = new ApiGatewayManagementApi({ endpoint: WEBSOCKET_URL });
    const postCalls = await Promise.all(connections.Items.map(async ({ connectionId, deviceType }) => {
      const response = await postToConnectionById(connectionId, body, api, db);
      return { connectionId, deviceType, response };
    }));

    const activeConnections = postCalls.filter(({ response }) => response === ACTIVE_CONNECTION);
    const deletedConnections = postCalls.filter(({ response }) => response === STALE_CONNECTION);

    console.log('Success: Message sent!', { activeConnections, deletedConnections, message: body });
    return { statusCode: 200, body: { activeConnections, deletedConnections, message: body } };
  } catch (error) {
    console.error("Error: Message has not been sent", error);
    return { statusCode: 500, body: error.stack };
  }
};
