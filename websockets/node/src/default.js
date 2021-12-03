exports.handler = async function(event, context, callback) {
  console.log('DEFAULT HANDLER');

	if (event.body === '__ping__') {
    return { statusCode: 200, body: '__pong__' };
  }
};
