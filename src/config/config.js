// require('dotenv').config();
module.exports = {
    development: {
      username: 'wishlist',
      password: 'changeme',
      database: 'wishlist',
      host: '127.0.0.1',
      dialect: 'mysql'
    },
    test: {
        username: process.env.DB_USERNAME,
        password: process.env.DB_PASSWORD,
        database: process.env.DB_NAME,
        host: process.env.DB_HOSTNAME,
        dialect: 'mysql'
    },
    production: {
      username: process.env.DB_USERNAME,
      password: process.env.DB_PASSWORD,
      database: process.env.DB_NAME,
      host: process.env.DB_HOSTNAME,
      dialect: 'mysql'
    }
  };
