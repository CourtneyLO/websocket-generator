const DynamoDB = require('aws-sdk/clients/dynamodb');

const { WEBSOCKET_MANAGER_TABLE_NAME } = process.env;

const scanDBForDevice = async (db, deviceType) => {
  try {
    const scanParams = {
      TableName: WEBSOCKET_MANAGER_TABLE_NAME,
      FilterExpression: '#deviceType = :deviceTypeValue',
      ExpressionAttributeNames: {
        "#deviceType": "deviceType",
      },
      ExpressionAttributeValues: {
        ':deviceTypeValue': deviceType
      },
    };

    const connectionsByDevice = await db.scan(scanParams).promise();
    return connectionsByDevice.Items;
  } catch (error) {
    console.error(`Failed to scan database for existing device types: ${deviceType}`, error);
    return [];
  }
};

const deleteExistingRowsWithDeviceType = async (db, items, deviceType) => {
  try {
   Promise.all(items.map(async ({ connectionId }) => {
      const deleteParams = {
        TableName: WEBSOCKET_MANAGER_TABLE_NAME,
        Key: { connectionId: connectionId },
      };

      await db.delete(deleteParams).promise();
      return;
    }));
    return;
  } catch (error) {
    console.error(`Failed to delete ${items.length} rows from the database for existing device types: ${deviceType}`, error);
  }
};

const addNewDeviceConnectionToDB = async (db, connectionId, deviceType) => {
  try {
    const putParams = {
      TableName: WEBSOCKET_MANAGER_TABLE_NAME,
      Item: {
        connectionId: connectionId,
        timestamp: `${new Date()}`
      }
    };

    if (deviceType) {
      putParams.Item.deviceType = deviceType;
    }

    await db.put(putParams).promise();
    console.log(`New connectionID ${connectionId} and device ${deviceType} have been added to the database`);
    return;
  } catch (error) {
    console.error(`Failed to add connectionId ${connectionId} and device ${deviceType} to the database`, error);
    return error;
  }
};

exports.handler = async function(event, context, callback) {
  console.log('Connect Handler');

  const deviceType = event.queryStringParameters && event.queryStringParameters.deviceType;

  if (!deviceType) {
    console.log('No deviceType was given as a query string parameter');
  }

  const db = new DynamoDB.DocumentClient();
  // Scan DB for any records with the same device type name as this connection
  const items = deviceType && await scanDBForDevice(db, deviceType);
  // Delete existing device type rows from the database so that when we add the new connection there will be only one connectionId for that device type
  deviceType && items.length && await deleteExistingRowsWithDeviceType(db, items, deviceType)
  // Add new connectionId and deviceType to the database
  const connectionId = event.requestContext.connectionId;
  const error = await addNewDeviceConnectionToDB(db, connectionId, deviceType);

  if (error) {
    return { statusCode: 500, body: "Failed to connect: " + JSON.stringify(error) };
  }

  return { statusCode: 200, body: `Connection made with the connection ID of ${connectionId}` };
};
