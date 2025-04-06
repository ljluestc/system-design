const otEngine = require('./otEngine');
const sessionManager = require('./sessionManager');
const logger = require('../shared/logger');

module.exports = (socket, io) => {
  // Log connection details
  logger.log(`Socket connected: ${socket.id}`);

  // Handle user joining a document session
  socket.on('join', (data) => {
    const { docId, userId } = data;
    if (!docId || !userId) {
      logger.log(`Invalid join request from ${socket.id}: Missing docId or userId`);
      socket.emit('error', { message: 'docId and userId are required' });
      return;
    }

    try {
      socket.join(docId);
      sessionManager.joinSession(docId, userId);
      logger.log(`User ${userId} joined doc ${docId} via socket ${socket.id}`);
      socket.emit('joined', { docId, userId, message: 'Successfully joined document' });

      // Notify other users in the room
      socket.to(docId).emit('userJoined', { userId, timestamp: new Date().toISOString() });
    } catch (err) {
      logger.log(`Error in join handler: ${err.message}`);
      socket.emit('error', { message: 'Failed to join document' });
    }
  });

  // Handle document operations (e.g., text edits)
  socket.on('operation', (data) => {
    const { docId, operation, userId } = data;
    if (!docId || !operation || !userId) {
      logger.log(`Invalid operation from ${socket.id}: Missing required fields`);
      socket.emit('error', { message: 'docId, operation, and userId are required' });
      return;
    }

    try {
      logger.log(`Received operation for doc ${docId} from user ${userId}`);
      const transformedOp = otEngine.applyOperation(operation);
      logger.log(`Transformed operation: ${JSON.stringify(transformedOp)}`);

      // Broadcast the transformed operation to all clients in the room except sender
      socket.to(docId).emit('update', {
        docId,
        operation: transformedOp,
        userId,
        timestamp: new Date().toISOString(),
      });

      // Acknowledge to sender
      socket.emit('operationAck', {
        docId,
        operationId: operation.id,
        status: 'applied',
      });
    } catch (error) {
      logger.log(`Error applying operation: ${error.message}`);
      socket.emit('error', { message: 'Failed to apply operation' });
    }
  });

  // Handle user leaving a document
  socket.on('leave', (data) => {
    const { docId, userId } = data;
    if (!docId || !userId) {
      logger.log(`Invalid leave request from ${socket.id}`);
      return;
    }

    try {
      socket.leave(docId);
      sessionManager.leaveSession(docId, userId);
      logger.log(`User ${userId} left doc ${docId}`);
      socket.to(docId).emit('userLeft', { userId, timestamp: new Date().toISOString() });
    } catch (err) {
      logger.log(`Error in leave handler: ${err.message}`);
    }
  });

  // Handle disconnection
  socket.on('disconnect', () => {
    logger.log(`Socket disconnected: ${socket.id}`);
    // Cleanup sessions if necessary
    const rooms = Object.keys(socket.rooms).filter((room) => room !== socket.id);
    rooms.forEach((docId) => {
      const users = sessionManager.getSessionUsers(docId);
      if (users) {
        users.forEach((userId) => {
          io.to(docId).emit('userLeft', { userId, timestamp: new Date().toISOString() });
        });
      }
    });
  });

  // Additional event handlers can be added here (e.g., ping, sync)
  socket.on('ping', () => {
    socket.emit('pong', { timestamp: new Date().toISOString() });
  });
};