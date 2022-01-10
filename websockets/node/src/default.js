exports.handler = async function(event) {
  console.log('Default Handler');

	if (event.body === '__ping__') {
    return { statusCode: 200, body: '__pong__' };
  }
};
