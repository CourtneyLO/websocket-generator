require('./test-helpers'); // Comment out if you want to see the console.logs in the tests
const defaultHandler = require('./default');

const event = { body: '__ping__' };

test('the response statusCode returned is 200 when __ping__ is sent in the request body', async () => {
  const response = await defaultHandler.handler(event);
  expect(response.statusCode).toBe(200);
});

test('the response body returned is __pong__ when __ping__ is sent in the request body', async () => {
  const response = await defaultHandler.handler(event);
  expect(response.body).toBe("__pong__");
});

test('Response is undefined when __ping__ is not given in the request body', async () => {
  const event = { body: 'something else' };
  const response = await defaultHandler.handler(event);
  expect(response).toBeUndefined();
});
