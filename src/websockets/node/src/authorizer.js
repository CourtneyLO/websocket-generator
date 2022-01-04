const AWS = require("aws-sdk");

const lambdaClient = new AWS.Lambda({ region: process.env.AWS_REGION_VALUE });
const secretsManagerClient = new AWS.SecretsManager({ region: process.env.AWS_REGION_VALUE });

const getSecretValue = async (event, context, callback) => {
  console.log('SECRETS HANDLER');

  const secretData = await secretsManagerClient.getSecretValue({
    SecretId: `${process.env.AUTHORIZATION_SECRET_NAME}`
  }).promise();

  console.log('secret returned');
  return JSON.parse(secretData.SecretString);
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
  console.log('AUTHORIZER HANDLER');

  const queryStringParameters = event.queryStringParameters;
  const secret = await getSecretValue();

  if (queryStringParameters[`${process.env.AUTHORIZATION_KEY}`] === secret) {
    console.log("Access Granted");
    callback(null, generateAllowPolicy(event.methodArn));
  } else {
    console.log("Unauthorized");
    callback("Unauthorized");
  }
};
