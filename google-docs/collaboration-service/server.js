const express = require('express');
const http = require('http');
const socketIo = require('socket.io');
const socketHandler = require('./socketHandler');
const logger = require('../shared/logger');
const config = require('./config');
const { STATUS_OK } = require('../shared/constants');

const app = express();
const server = http.createServer(app);
const io = socketIo(server, {
  cors: {
    origin: '*',
    methods: ['GET', 'POST'],
  },
});

// Middleware for basic API endpoints
app.use(express.json());

// Health check endpoint
app.get('/health', (req, res) => {
  const connectedClients = io.engine.clientsCount;
  res.status(STATUS_OK).json({
    status: 'healthy',
    uptime: process.uptime(),
    connectedClients,
    timestamp: new Date().toISOString(),
  });
});

// Log all incoming HTTP requests
app.use((req, res, next) => {
  logger.log(`HTTP request: ${req.method} ${req.url}`);
  next();
});

// WebSocket connection handler
io.on('connection', (socket) => {
  logger.log(`New WebSocket connection: ${socket.id}`);
  socketHandler(socket, io);

  socket.on('disconnect', () => {
    logger.log(`WebSocket disconnected: ${socket.id}`);
  });

  socket.on('error', (err) => {
    logger.log(`WebSocket error: ${err.message}`);
  });
});

// Error handling for HTTP server
app.use((err, req, res, next) => {
  logger.log(`Server error: ${err.message}`);
  res.status(500).json({ error: 'Internal Server Error', message: err.message });
});

// Start the server
server.listen(config.port, () => {
  logger.log(`Collaboration Service running on port ${config.port}`);
  logger.log(`Environment: ${process.env.NODE_ENV || 'development'}`);
});

// Graceful shutdown
process.on('SIGTERM', () => {
  logger.log('SIGTERM received. Shutting down gracefully...');
  io.close(() => {
    logger.log('WebSocket server closed.');
  });
  server.close(() => {
    logger.log('HTTP server closed.');
    process.exit(0);
  });
});

process.on('SIGINT', () => {
  logger.log('SIGINT received. Shutting down gracefully...');
  io.close(() => {
    logger.log('WebSocket server closed.');
  });
  server.close(() => {
    logger.log('HTTP server closed.');
    process.exit(0);
  });
});

// Monitor server health
function monitorHealth() {
  setInterval(() => {
    const memoryUsage = process.memoryUsage();
    const clients = io.engine.clientsCount;
    logger.log(`Health check: ${clients} clients connected`);
    logger.log(`Memory: RSS=${memoryUsage.rss}, Heap Total=${memoryUsage.heapTotal}`);
  }, 30000); // Every 30 seconds
}

monitorHealth();

// Additional logic can be added here (e.g., connection throttling, metrics)