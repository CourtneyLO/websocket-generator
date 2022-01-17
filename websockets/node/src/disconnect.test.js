const tableName = 'example_websocket_manager'
process.env.WEBSOCKET_MANAGER_TABLE_NAME = tableName;

// Comment out the line below if you want to see the console.logs in the tests
require('./test-helpers/silence-console-logs');
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
});
