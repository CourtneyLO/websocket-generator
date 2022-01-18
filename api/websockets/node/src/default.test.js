const tableName = 'example_websocket_manager'
process.env.WEBSOCKET_MANAGER_TABLE_NAME = tableName;
const websocketURL = 'https://example.com'
process.env.WEBSOCKET_URL = websocketURL;

// Comment out the line below if you want to see the console.logs in the tests
require('./test-helpers/silence-console-logs');
const defaultHandler = require('./default');

const deviceType = 'computer';
const connectionId = '1234';
const mockScanPromiseValue = { promise: jest.fn().mockReturnValue({
  Items: [{ connectionId, deviceType }]
}) };
const mockDeletePromiseValue = { promise: jest.fn() }
const mockDocumentClientReturnValue = {
  scan: jest.fn(() => mockScanPromiseValue),
  delete: jest.fn(() => mockDeletePromiseValue)
};
jest.mock('aws-sdk/clients/dynamodb', () => {
  return { DocumentClient: jest.fn(() => mockDocumentClientReturnValue) }
});

const mockPromise = { promise: jest.fn().mockReturnValue('success') };
const mockApiGatewayManagementApi = { postToConnection: jest.fn(() => mockPromise) }
const ApiGatewayManagementApi = require('aws-sdk/clients/apigatewaymanagementapi');
jest.mock('aws-sdk/clients/apigatewaymanagementapi', () => {
  return jest.fn().mockImplementation(() => {
    return mockApiGatewayManagementApi;
  })
});

describe('Default Lambda', () => {
  describe('__ping__, __pong__', () => {
    const event = { body: '__ping__' };

    test('the response statusCode returned is 200 when __ping__ is sent in the request body', async () => {
      const response = await defaultHandler.handler(event);
      expect(response.statusCode).toBe(200);
    });

    test('the response body returned is __pong__ when __ping__ is sent in the request body', async () => {
      const response = await defaultHandler.handler(event);
      expect(response.body).toBe("__pong__");
    });
  });

  describe('Generate Socket Message', () => {
    const message = 'Hello World';
    const event = { body: message };

    beforeEach(() => {
      mockScanPromiseValue.promise = jest.fn().mockReturnValue({
        Items: [{ connectionId }]
      });
    });

    afterEach(() => {
      ApiGatewayManagementApi.mockClear();
      mockDocumentClientReturnValue.scan.mockClear();
      mockDocumentClientReturnValue.delete.mockClear();
      mockApiGatewayManagementApi.postToConnection.mockClear();
      mockScanPromiseValue.promise.mockClear();
      mockDeletePromiseValue.promise.mockClear();
      mockPromise.promise.mockClear();
    });

    describe('Success', () => {
      test('db.scan is called with the table name', async () => {
        await defaultHandler.handler(event);
        const mockedScanParams = {
          TableName: tableName,
        }
        expect(mockDocumentClientReturnValue.scan).toBeCalledTimes(1);
        expect(mockDocumentClientReturnValue.scan).toBeCalledWith(mockedScanParams);
      });

      test('ApiGatewayManagementApi is initialised with an endpoint url', async () => {
        await defaultHandler.handler(event);
        const mockedApiParams = {
          endpoint: websocketURL,
        }
        expect(ApiGatewayManagementApi).toBeCalledTimes(1);
        expect(ApiGatewayManagementApi).toBeCalledWith(mockedApiParams);
      });

      test('postToConnection is called a ConnectionId and Data', async () => {
        await defaultHandler.handler(event);
        expect(mockApiGatewayManagementApi.postToConnection).toBeCalledTimes(1);
        expect(mockApiGatewayManagementApi.postToConnection).toBeCalledWith({
          ConnectionId: connectionId,
          Data: JSON.stringify(message)
        });
      });

      test('the response has a statusCode of 200 and a body containing a single active connection, no deleted connections and the message that was sent to the client', async () => {
        const response = await defaultHandler.handler(event);
        expect(response.statusCode).toBe(200);
        expect(JSON.parse(response.body)).toEqual({
          activeConnections: [{
            connectionId: '1234',
            deviceType: undefined,
            response: "Message sent to connection"
          }],
          message: 'Hello World',
          deletedConnections: []
        });
      });

      test('the response has a single deleted connection, 2 active connections and the message that was sent to the clients', async () => {
        mockScanPromiseValue.promise = jest.fn().mockReturnValue({
          Items: [
            { connectionId, deviceType },
            { connectionId: '456', deviceType: 'tablet' },
            { connectionId: '789', deviceType: 'mobile' }
          ]
        });
        mockDeletePromiseValue.promise = jest.fn();
        mockPromise.promise = jest.fn()
          .mockReturnValueOnce(Promise.resolve({}))
          .mockReturnValueOnce(Promise.reject({ statusCode: 410, stack: 'Stale connection' }))
          .mockReturnValueOnce(Promise.resolve({}))

        const response = await defaultHandler.handler(event);
        expect(response.statusCode).toBe(200);
        expect(JSON.parse(response.body)).toEqual({
          activeConnections: [
            {
              connectionId: '1234',
              deviceType: 'computer',
              response: 'Message sent to connection'
            },
            {
              connectionId: '789',
              deviceType: 'mobile',
              response: 'Message sent to connection'
            }
          ],
          deletedConnections: [
            {
              connectionId: '456',
              deviceType: 'tablet',
              response: 'Deleted stale ConnectionId'
            }
          ],
          message: 'Hello World'
        });
      });
    });

    describe('Error', () => {
      describe('response statusCode is 500 and the body contains an error message', () => {
        test('a db.scan functionality error', async () => {
          mockScanPromiseValue.promise = jest.fn().mockReturnValueOnce(Promise.reject({ stack: 'Something went wrong with db.scan!' }));
          response = await defaultHandler.handler(event);
          expect(response.statusCode).toBe(500);
          expect(response.body).toBe('Something went wrong with db.scan!');
        });

        test('a postToConnection functionality error', async () => {
          mockPromise.promise = jest.fn(() => Promise.reject({ stack: 'Something went wrong with postToConnection!' }));
          response = await defaultHandler.handler(event);
          expect(response.statusCode).toBe(500);
          expect(response.body).toBe('Something went wrong with postToConnection!');
        });

        test('a db.delete functionality error', async () => {
          mockDeletePromiseValue.promise = jest.fn().mockReturnValueOnce(() => Promise.reject({ stack: 'Something went wrong with db.delete!' }))
          response = await defaultHandler.handler(event);
          expect(response.statusCode).toBe(500);
          expect(response.body).toBe('Something went wrong with postToConnection!');
        });
      });

      describe('postToConnection error returns statusCode 410', () => {
        test('db.delete is called with table name and connectionId to delete a stale connection', async () => {
          const mockedDeleteParams = {
            TableName: tableName,
            Key: { connectionId: connectionId }
          };
          mockPromise.promise = jest.fn(() => Promise.reject({ statusCode: 410, stack: 'Stale connection' }));
          await defaultHandler.handler(event);
          expect(mockDocumentClientReturnValue.delete).toBeCalledTimes(1);
          expect(mockDocumentClientReturnValue.delete).toBeCalledWith(mockedDeleteParams);
        });
      })
    });
  });
});
