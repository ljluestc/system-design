const express = require('express');
const routes = require('./routes');
const config = require('./config');
const logger = require('../shared/logger');
const { sanitizeInput } = require('../shared/utils');
const { STATUS_OK, STATUS_ERROR } = require('../shared/constants');

// Initialize Express app
const app = express();

// Middleware for parsing JSON requests
app.use(express.json());

// Request logging middleware
app.use((req, res, next) => {
  const timestamp = new Date().toISOString();
  logger.log(`[${timestamp}] Incoming request: ${req.method} ${req.url}`);
  next();
});

// Sanitize all incoming request bodies
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

// Health check endpoint
app.get('/health', (req, res) => {
  res.status(STATUS_OK).json({
    status: 'healthy',
    uptime: process.uptime(),
    timestamp: new Date().toISOString(),
  });
});

// Mount user routes
app.use('/users', routes);

// 404 handler
app.use((req, res, next) => {
  res.status(404).json({ error: `Route ${req.method} ${req.url} not found` });
});

// Global error handling middleware
app.use((err, req, res, next) => {
  logger.log(`[${new Date().toISOString()}] Error: ${err.message}`);
  logger.log(`Stack trace: ${err.stack}`);
  res.status(STATUS_ERROR).json({
    error: 'Internal Server Error',
    message: err.message,
    timestamp: new Date().toISOString(),
  });
});

// Start the server
const server = app.listen(config.port, () => {
  logger.log(`User Service running on port ${config.port}`);
  logger.log(`Environment: ${process.env.NODE_ENV || 'development'}`);
});

// Graceful shutdown handling
process.on('SIGTERM', () => {
  logger.log('SIGTERM received. Initiating graceful shutdown...');
  server.close(() => {
    logger.log('User Service shut down successfully.');
    process.exit(0);
  });
});

process.on('SIGINT', () => {
  logger.log('SIGINT received. Initiating graceful shutdown...');
  server.close(() => {
    logger.log('User Service shut down successfully.');
    process.exit(0);
  });
});

// Unhandled promise rejection handler
process.on('unhandledRejection', (reason, promise) => {
  logger.log(`Unhandled Rejection at: ${promise}, reason: ${reason}`);
  process.exit(1);
});

// Uncaught exception handler
process.on('uncaughtException', (err) => {
  logger.log(`Uncaught Exception: ${err.message}`);
  logger.log(`Stack trace: ${err.stack}`);
  process.exit(1);
});

// Additional utility function to simulate load
function simulateLoad() {
  logger.log('Simulating server load...');
  setInterval(() => {
    const memoryUsage = process.memoryUsage();
    logger.log(`Memory usage: RSS=${memoryUsage.rss}, Heap Total=${memoryUsage.heapTotal}`);
  }, 60000); // Log memory usage every minute
}

simulateLoad();

// Placeholder for additional middleware or logic
// Example: Rate limiting, CORS, etc., can be added here to expand the file further