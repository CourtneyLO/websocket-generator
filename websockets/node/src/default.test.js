const tableName = 'example_websocket_manager'
process.env.WEBSOCKET_MANAGER_TABLE_NAME = tableName;
const websocketURL = 'https://example.com'
process.env.WEBSOCKET_URL = websocketURL;

require('./test-helpers'); // Comment out if you want to see the console.logs in the tests
const defaultHandler = require('./default');

const connectionId = '1234';
const mockScanPromiseValue = { promise: jest.fn().mockReturnValue({
  Items: [{ connectionId }]
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

      test('the response has a statusCode of 200 and there is a message in the body', async () => {
        const response = await defaultHandler.handler(event);
        expect(response.statusCode).toBe(200);
        expect(response.body).toBe("Message sent!");
      });
    });

    describe('Error', () => {
      describe('response statusCode is 500 and the body contains an error message', () => {
        test('a db.scan functionality error', async () => {
          mockScanPromiseValue.promise = jest.fn(() => Promise.reject({ stack: 'Something went wrong with db.scan!' }));
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
          mockDeletePromiseValue.promise = jest.fn(() => Promise.reject({ stack: 'Something went wrong with db.delete!' }));
          response = await defaultHandler.handler(event);
          expect(response.statusCode).toBe(500);
          expect(response.body).toBe('Something went wrong with postToConnection!');
        });

        test('a 410 postToConnection error', async () => {
          mockDeletePromiseValue.promise = jest.fn();
          mockPromise.promise = jest.fn(() => Promise.reject({ statusCode: 410, stack: 'Stale connection' }));
          response = await defaultHandler.handler(event);
          expect(response.statusCode).toBe(500);
          expect(response.body).toBe('Stale connection');
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
