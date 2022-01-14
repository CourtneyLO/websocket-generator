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

  test('the table name and connection ID are set as the parameters for database row deletion request', async () => {
    const response = await disconnectHandler.handler(event);
    const mockedDeleteParams = {
      TableName: tableName,
      Key: {
        connectionId: connectionId
      }
    }
    expect(mockDocumentClientReturnValue.delete).toBeCalledTimes(1);
    expect(mockDocumentClientReturnValue.delete).toBeCalledWith(mockedDeleteParams);
  });


  test('the response statusCode returned is 200 when connection ID has been deleted from the database', async () => {
    const response = await disconnectHandler.handler(event);
    expect(response.statusCode).toBe(200);
  });

  test('the response body returned is "Disconnected" when connection ID has been deleted from the database', async () => {
    const response = await disconnectHandler.handler(event);
    expect(response.body).toBe("Disconnected");
  });

  test('the response statusCode returned is 500 when an error occurs deleting a connection ID from the database', async () => {
    mockDeleteReturnValue.promise = jest.fn(() => Promise.reject("Something went wrong!"));
    const response = await disconnectHandler.handler(event);
    expect(response.statusCode).toBe(500);
  });

  test('the response body returned is Failed to disconnect with an error when an error occurs deleting a connection ID from the database', async () => {
    mockDeleteReturnValue.promise = jest.fn(() => Promise.reject("Something went wrong!"));
    const response = await disconnectHandler.handler(event);
    expect(response.body).toBe("Failed to disconnect: \"Something went wrong!\"");
  });

});
