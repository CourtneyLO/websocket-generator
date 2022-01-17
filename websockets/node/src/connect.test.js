const tableName = 'example_websocket_manager';
process.env.WEBSOCKET_MANAGER_TABLE_NAME = tableName;

// Comment out the line below if you want to see the console.logs in the tests
require('./test-helpers/silence-console-logs');
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

describe('Connect Handler', () => {
  const deviceType = 'computer';
  const event = {
    requestContext: { connectionId: connectionId },
    queryStringParameters: { deviceType: deviceType }
  };
  const mockDate = new Date(1466424490000)
  const spy = jest
    .spyOn(global, 'Date')
    .mockImplementation(() => mockDate)

  afterEach(() => {
    mockDocumentClientReturnValue.scan.mockClear();
    mockDocumentClientReturnValue.delete.mockClear();
    mockDocumentClientReturnValue.put.mockClear();
    mockPromiseReturnValue.promise.mockClear();
  });

  describe('Success', () => {
    test('db.scan is called with TableName, FilterExpression, ExpressionAttributeNames and ExpressionAttributeValues to get items from the database', async () => {
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
      expect(mockDocumentClientReturnValue.scan).toBeCalledTimes(1);
      expect(mockDocumentClientReturnValue.scan).toBeCalledWith(mockedScanParams);
    });

    test('db.delete is called with TableName and Key.ConnectionId to delete a row from the database', async () => {
      await connectHandler.handler(event);
      const mockedDeleteParams = {
        TableName: tableName,
        Key: { connectionId: connectionId }
      }
      expect(mockDocumentClientReturnValue.delete).toBeCalledTimes(1);
      expect(mockDocumentClientReturnValue.delete).toBeCalledWith(mockedDeleteParams);
    });

    test('db.put is called with TableName and Item to add a row to the database', async () => {
      await connectHandler.handler(event);
      const mockedPutParams = {
        TableName: tableName,
        Item: {
          connectionId: connectionId,
          timestamp: `${mockDate}`,
          deviceType: deviceType
        }
      }
      expect(mockDocumentClientReturnValue.put).toBeCalledTimes(1);
      expect(mockDocumentClientReturnValue.put).toBeCalledWith(mockedPutParams);
    });

    test('db.put is called with WITHOUT device type as it was not given in the url', async () => {
      const event = { requestContext: { connectionId: connectionId } };
      await connectHandler.handler(event);
      const mockedPutParams = {
        TableName: tableName,
        Item: {
          connectionId: connectionId,
          timestamp: `${mockDate}`
        }
      }
      expect(mockDocumentClientReturnValue.put).toBeCalledTimes(1);
      expect(mockDocumentClientReturnValue.put).toBeCalledWith(mockedPutParams);
    });
  });
});
