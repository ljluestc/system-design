const express = require('express');
const routes = require('./routes');
const config = require('./config');
const logger = require('../shared/logger');
const { sanitizeInput } = require('../shared/utils');
const { STATUS_OK, STATUS_ERROR } = require('../shared/constants');

const app = express();

// Middleware setup
app.use(express.json());

// Request logging
app.use((req, res, next) => {
  logger.log(`Incoming request: ${req.method} ${req.url}`);
  next();
});

// Sanitize inputs
app.use((req, res, next) => {
  if (req.body) {
    Object.keys(req.body).forEach((key) => {
      if (typeof req.body[key] === 'string') {
        req.body[key] = sanitizeInput(req.body[key]);
      }
    });
  }
  next();
});

// Health check
app.get('/health', (req, res) => {
  res.status(STATUS_OK).json({ status: 'healthy', uptime: process.uptime() });
});

// Mount routes
app.use('/docs', routes);

// 404 handler
app.use((req, res) => {
  res.status(404).json({ error: 'Not Found' });
});

// Error handler
app.use((err, req, res, next) => {
  logger.log(`Error: ${err.message}`);
  res.status(STATUS_ERROR).json({ error: 'Internal Server Error' });
});

// Start server
const server = app.listen(config.port, () => {
  logger.log(`Document Service running on port ${config.port}`);
});

// Graceful shutdown
process.on('SIGTERM', () => {
  logger.log('SIGTERM received. Shutting down...');
  server.close(() => {
    logger.log('Document Service shut down.');
  });
});

// Additional logic for monitoring, etc.