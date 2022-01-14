const tableName = 'example_websocket_manager'
process.env.WEBSOCKET_MANAGER_TABLE_NAME = tableName;

require('./test-helpers'); // Comment out if you want to see the console.logs in the tests
const disconnectHandler = require('./disconnect');

const mockDeleteReturnValue = { promise: jest.fn() };
const mockDocumentClientReturnValue = {
  delete: jest.fn(() => mockDeleteReturnValue)
};
jest.mock('aws-sdk/clients/dynamodb', () => {
  return { DocumentClient: jest.fn(() => mockDocumentClientReturnValue) }
});

describe('Disconnect Handler', () => {
  const connectionId = '1234';
  const event = { requestContext: { connectionId: connectionId } };

  afterEach(() => {
    mockDocumentClientReturnValue.delete.mockClear();
    mockDeleteReturnValue.promise.mockClear();
  });

  test('db.delete is called with TableName and Key.connectionId to remove a connectionId from the database', async () => {
    const response = await disconnectHandler.handler(event);
    const mockedDeleteParams = {
      TableName: tableName,
      Key: { connectionId: connectionId }
    }
    expect(mockDocumentClientReturnValue.delete).toBeCalledTimes(1);
    expect(mockDocumentClientReturnValue.delete).toBeCalledWith(mockedDeleteParams);
  });


  test('the response statusCode is 200 and the body returns a message when connection ID has been deleted from the database', async () => {
    const response = await disconnectHandler.handler(event);
    expect(response.statusCode).toBe(200);
    expect(response.body).toBe("Disconnected");
  });

  test('the response statusCode is 500 and the body is an error message when a connentionId fails to be deleted from the database', async () => {
    mockDeleteReturnValue.promise = jest.fn(() => Promise.reject("Something went wrong!"));
    const response = await disconnectHandler.handler(event);
    expect(response.statusCode).toBe(500);
    expect(response.body).toBe("Failed to disconnect: \"Something went wrong!\"");
  });
});
