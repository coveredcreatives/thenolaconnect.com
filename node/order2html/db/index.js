const { Pool } = require('pg');
const util = require('util');
const pool = new Pool();
const { logger } = require("../logger")

module.exports = {
  query: (text, params) => {
    const start = Date.now();
    return new Promise((resolve, reject) => {
      pool.query(text, params, (err, res) => {
        const duration = Date.now() - start;
        if (err) {
          logger.log({
            level: 'error',
            message: 'executed query',
            err: err,
            text: text,
            duration: duration,
          });
          return reject(err);
        }
        logger.log({
          level: 'info',
          message: 'executed query',
          rows: res.rowCount,
          text: text,
          duration: duration,
        });
        resolve(res);
      });
    });
  },
}