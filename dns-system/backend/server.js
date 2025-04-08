const express = require('express');
const { connectDB } = require('./config/database');
const { connectRedis } = require('./config/redis');
const logger = require('./utils/logger');
const app = require('./app');

const startServer = async () => {
  try {
    await connectDB();
    await connectRedis();
    app.listen(3001, () => logger.info('Server running on port 3001'));
  } catch (error) {
    logger.error(`Server startup failed: ${error.message}`);
    process.exit(1);
  }
};

startServer();