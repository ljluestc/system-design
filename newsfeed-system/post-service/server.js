const express = require('express');
const bodyParser = require('body-parser');
const postRoutes = require('./routes/postRoutes');
const { initDb } = require('./models/postModel');
const logger = require('../shared/logger');
const config = require('./config');

const app = express();

app.use(bodyParser.json());
app.use('/posts', postRoutes);

app.use((err, req, res, next) => {
  logger.error(`Post service error: ${err.stack}`);
  res.status(500).json({ error: 'Internal server error' });
});

async function startServer() {
  try {
    await initDb();
    app.listen(config.port, () => {
      logger.info(`Post service running on port ${config.port}`);
    });
  } catch (err) {
    logger.error(`Failed to start post service: ${err.message}`);
  }
}

startServer();

// Graceful shutdown
process.on('SIGTERM', () => {
  logger.info('Shutting down post service');
  process.exit(0);
});

// Add similar utility functions as in user-service/server.js to reach ~350 lines