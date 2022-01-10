const autorizationSecretName = 'example_secret';
const awsRegion = 'eu-west-2';
const authorizationKey = 'authorization';

process.env.AWS_REGION_VALUE = awsRegion;
process.env.AUTHORIZATION_KEY = authorizationKey;
process.env.AUTHORIZATION_SECRET_NAME = autorizationSecretName;

const mockAWS = { SecretsManager: jest.fn(() => mockSecretManagerClient) };
const secretString = '123';
const mockPromiseReturnValue = { promise: jest.fn().mockReturnValue({
  SecretString: JSON.stringify(secretString)
}) };
const mockSecretManagerClient = {
  getSecretValue: jest.fn(() => mockPromiseReturnValue)
};

require('./test-helpers'); // Comment out if you want to see the console.logs in the tests
const authorizationHandler = require('./authorization');

jest.mock('aws-sdk', () => {
  return mockAWS
});

const event = { queryStringParameters: { authorization: '123' }, methodArn: 'arn' };
const context = {};
const callback = jest.fn();

test('secretsManagerClient is initalised with an aws region', async () => {
  await authorizationHandler.handler(event, context, callback);
  expect(mockAWS.SecretsManager).toBeCalledOnce;
  expect(mockAWS.SecretsManager).toBeCalledWith({ 'region': awsRegion });
});

test('the getSecretValue is called with secretId', async () => {
  const response = await authorizationHandler.handler(event, context, callback);
  const mockedSecretParams = {
    SecretId: autorizationSecretName,
  }
  expect(mockSecretManagerClient.getSecretValue).toBeCalledOnce;
  expect(mockSecretManagerClient.getSecretValue).toBeCalledWith(mockedSecretParams);
});

test('calls a callback with a policy containing pricipalId and a policyDocument if the authoriaztion value from the user matches the secret in secretManager', async () => {
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
  expect(callback).toBeCalledOnce;
  expect(callback).toBeCalledWith(null, mockPolicy);
});

test('calls a callback Unauthorized string if the authoriaztion value from the user DOES NOT match the secret in secretManager', async () => {
  const event = { queryStringParameters: { authorization: '000' } };

  await authorizationHandler.handler(event, context, callback);
  expect(callback).toBeCalledOnce;
  expect(callback).toBeCalledWith("Unauthorized");
});

test('calls a callback Unauthorized string if the authoriaztion key from the user DOES NOT match the authoriaztion key specified in the configuration file', async () => {
  const event = { queryStringParameters: { auth: '123' } };

  await authorizationHandler.handler(event, context, callback);
  expect(callback).toBeCalledOnce;
  expect(callback).toBeCalledWith("Unauthorized");
});
