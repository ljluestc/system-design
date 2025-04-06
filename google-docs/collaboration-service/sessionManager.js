const sessions = new Map();
const logger = require('../shared/logger');

function joinSession(docId, userId) {
  if (!docId || !userId) {
    throw new Error('docId and userId are required');
  }

  if (!sessions.has(docId)) {
    sessions.set(docId, new Set());
  }
  sessions.get(docId).add(userId);
  logger.log(`User ${userId} joined session for doc ${docId}`);
}

function leaveSession(docId, userId) {
  if (!sessions.has(docId)) {
    logger.log(`No session found for doc ${docId}`);
    return;
  }
  sessions.get(docId).delete(userId);
  logger.log(`User ${userId} left session for doc ${docId}`);
  if (sessions.get(docId).size === 0) {
    sessions.delete(docId);
    logger.log(`Session for doc ${docId} closed`);
  }
}

function getSessionUsers(docId) {
  if (!sessions.has(docId)) {
    return [];
  }
  return Array.from(sessions.get(docId));
}

module.exports = { joinSession, leaveSession, getSessionUsers };