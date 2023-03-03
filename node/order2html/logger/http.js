const logger = require('./logger');
const morgan = require('morgan');

module.exports = morgan(
  ':method :url :status :response-time ms - :res[content-length]',
  {
    stream: {
      write: message => logger.info(message.substring(0, message.lastIndexOf('\n')))
    }
  }
);