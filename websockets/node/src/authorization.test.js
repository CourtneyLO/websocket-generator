const awsRegion = 'eu-west-2';
process.env.AWS_REGION_VALUE = awsRegion;
const authorizationKey = 'authorization';
process.env.AUTHORIZATION_KEY = authorizationKey;
const autorizationSecretName = 'example_secret';
process.env.AUTHORIZATION_SECRET_NAME = autorizationSecretName;

const secretString = '123';
const mockPromiseReturnValue = { promise: jest.fn().mockReturnValue({
  SecretString: secretString
}) };
const mockSecretManagerClient = {
  getSecretValue: jest.fn(() => mockPromiseReturnValue)
};
const mockAWS = { SecretsManager: jest.fn(() => mockSecretManagerClient) };

// Comment out the line below if you want to see the console.logs in the tests
require('./test-helpers/silence-console-logs');
const authorizationHandler = require('./authorization');

jest.mock('aws-sdk', () => { return mockAWS });

describe('Authorization Handler', () => {
  const event = { queryStringParameters: { authorization: '123' }, methodArn: 'arn' };
  const context = {};
  const callback = jest.fn();

  afterEach(() => {
    mockAWS.SecretsManager.mockClear();
    mockSecretManagerClient.getSecretValue.mockClear();
    callback.mockClear();
  });

  test('secretsManagerClient is initalised with an aws region', async () => {
    await authorizationHandler.handler(event, context, callback);
    expect(mockAWS.SecretsManager).toBeCalledTimes(1);
    expect(mockAWS.SecretsManager).toBeCalledWith({ 'region': awsRegion });
  });

  test('the getSecretValue is called with secretId', async () => {
    const response = await authorizationHandler.handler(event, context, callback);
    const mockedSecretParams = {
      SecretId: autorizationSecretName,
    }
    expect(mockSecretManagerClient.getSecretValue).toBeCalledTimes(1);
    expect(mockSecretManagerClient.getSecretValue).toBeCalledWith(mockedSecretParams);
  });

  test('callback is called with a policy object containing pricipalId and a policyDocument if the authoriaztion value from the user matches the secret in secretManager', async () => {
    const mockPolicy = {
      principalId: 'UserGrantedAccess',
      policyDocument: {
        Version: '2012-10-17',
        Statement: [
          {
            Action: 'execute-api:Invoke',
            Effect: 'Allow',
            Resource: "arn"
          }
        ]
      }
    };

    await authorizationHandler.handler(event, context, callback);
    expect(callback).toBeCalledTimes(1);
    expect(callback).toBeCalledWith(null, mockPolicy);
  });

  test('callback is called with an Unauthorized string if the authoriaztion value from the user DOES NOT match the secret in secretManager', async () => {
    const event = { queryStringParameters: { authorization: '000' } };

    await authorizationHandler.handler(event, context, callback);
    expect(callback).toBeCalledTimes(1);
    expect(callback).toBeCalledWith("Unauthorized");
  });

  test('callback is called with an Unauthorized string if the authoriaztion key from the url DOES NOT match the authoriaztion key specified in the configuration file', async () => {
    const event = { queryStringParameters: { auth: '123' } };

    await authorizationHandler.handler(event, context, callback);
    expect(callback).toBeCalledTimes(1);
    expect(callback).toBeCalledWith("Unauthorized");
  });
});
