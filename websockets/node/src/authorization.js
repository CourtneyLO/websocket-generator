const AWS = require('aws-sdk');

const { AUTHORIZATION_SECRET_NAME, AWS_REGION_VALUE, AUTHORIZATION_KEY } = process.env

const secretsManagerClient = new AWS.SecretsManager({ region: AWS_REGION_VALUE });

const getSecretValue = async () => {
  const secretData = await secretsManagerClient.getSecretValue({
    SecretId: `${AUTHORIZATION_SECRET_NAME}`
  }).promise();

  return secretData.SecretString;
};

const generateAllowPolicy = function(resource) {
  const policyDocument = {
    Version: '2012-10-17',
    Statement: [
      {
        Action: 'execute-api:Invoke',
        Effect: 'Allow',
        Resource: resource
      }
    ]
  };

  return {
    principalId: 'UserGrantedAccess',
    policyDocument
  };
};

exports.handler = async (event, context, callback) => {
  console.log('Authoriaztion Handler');

  const queryStringParameters = event.queryStringParameters;
  const secret = await getSecretValue();

  if (queryStringParameters[`${AUTHORIZATION_KEY}`] === secret) {
    console.log("Success: Access Granted");
    callback(null, generateAllowPolicy(event.methodArn));
  } else {
    console.error("Error: Unauthorized");
    callback("Unauthorized");
  }
};
