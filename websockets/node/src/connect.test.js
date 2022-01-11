const tableName = 'example_websocket_manager'
process.env.WEBSOCKET_TABLE_NAME = tableName;

require('./test-helpers'); // Comment out if you want to see the console.logs in the tests
const connectHandler = require('./connect');

const connectionId = '1234';

const mockPromiseReturnValue = { promise: jest.fn().mockReturnValue({
  Items: [{ connectionId }]
}) };

const mockDocumentClientReturnValue = {
  delete: jest.fn(() => mockPromiseReturnValue),
  scan: jest.fn(() => mockPromiseReturnValue),
  put: jest.fn(() => mockPromiseReturnValue)
};

jest.mock('aws-sdk/clients/dynamodb', () => {
  return { DocumentClient: jest.fn(() => mockDocumentClientReturnValue) }
});

const deviceType = 'computer';

const event = {
  requestContext: { connectionId: connectionId },
  queryStringParameters: { deviceType: deviceType }
};

const mockDate = new Date(1466424490000)
const spy = jest
  .spyOn(global, 'Date')
  .mockImplementation(() => mockDate)

test('the table name and device type are set as the parameters for database row scan request', async () => {
  await connectHandler.handler(event);
  const mockedScanParams = {
    TableName: tableName,
    FilterExpression: '#deviceType = :deviceTypeValue',
    ExpressionAttributeNames: {
      "#deviceType": "deviceType",
    },
    ExpressionAttributeValues: {
      ':deviceTypeValue': deviceType
    }
  }
  expect(mockDocumentClientReturnValue.scan).toBeCalledOnce;
  expect(mockDocumentClientReturnValue.scan).toBeCalledWith(mockedScanParams);
});

test('the table name and connection ID are set as the parameters for database row deletion request', async () => {
  await connectHandler.handler(event);
  const mockedDeleteParams = {
    TableName: 'example_websocket_manager',
    Key: {
      connectionId: connectionId
    }
  }
  expect(mockDocumentClientReturnValue.delete).toBeCalledOnce;
  expect(mockDocumentClientReturnValue.delete).toBeCalledWith(mockedDeleteParams);
});

test('the table name and connection ID and device type are set as the parameters for database row add request', async () => {
  await connectHandler.handler(event);
  const mockedPutParams = {
    TableName: 'example_websocket_manager',
    Item: {
      connectionId: connectionId,
      timestamp: `${mockDate}`,
      deviceType: deviceType
    }
  }
  expect(mockDocumentClientReturnValue.put).toBeCalledOnce;
  expect(mockDocumentClientReturnValue.put).toBeCalledWith(mockedPutParams);
});

test('the table name and connection ID WITHOUT device type are set as the parameters for database row add request', async () => {
  const event = { requestContext: { connectionId: connectionId } };
  await connectHandler.handler(event);
  const mockedPutParams = {
    TableName: 'example_websocket_manager',
    Item: {
      connectionId: connectionId,
      timestamp: `${mockDate}`
    }
  }
  expect(mockDocumentClientReturnValue.put).toBeCalledOnce;
  expect(mockDocumentClientReturnValue.put).toBeCalledWith(mockedPutParams);
});

test('the response statusCode returned is 200 when connection ID has been deleted from the database', async () => {
  const response = await connectHandler.handler(event);
  expect(response.statusCode).toBe(200);
});

test('the response body returned is "Disconnected" when connection ID has been deleted from the database', async () => {
  const response = await connectHandler.handler(event);
  expect(response.body).toBe("Connection made with the connection ID of 1234");
});

test('the response statusCode returned is 500 when an error occurs deleting a connection ID from the database', async () => {
  mockPromiseReturnValue.promise = jest.fn(() => Promise.reject("Something went wrong!"));
  const response = await connectHandler.handler(event);
  expect(response.statusCode).toBe(500);
});

test('the response body returned is Failed to connect with an error when an error occurs deleting a connection ID from the database', async () => {
  mockPromiseReturnValue.promise = jest.fn(() => Promise.reject("Something went wrong!"));
  const response = await connectHandler.handler(event);
  expect(response.body).toBe("Failed to connect: \"Something went wrong!\"");
});
