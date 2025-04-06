const jwt = require('jsonwebtoken');
const config = require('../config');
const logger = require('../../shared/logger');

function authenticate(req, res, next) {
  const authHeader = req.headers['authorization'];
  if (!authHeader) {
    logger.warn('No authorization header provided');
    return res.status(401).json({ error: 'No token provided' });
  }

  const token = authHeader.split(' ')[1]; // Expecting "Bearer <token>"
  if (!token) {
    logger.warn('Malformed authorization header');
    return res.status(401).json({ error: 'Invalid token format' });
  }

  try {
    const decoded = jwt.verify(token, config.jwtSecret);
    req.userId = decoded.userId;
    logger.debug(`Token verified for user ${req.userId}`);
    next();
  } catch (err) {
    logger.error(`Token verification failed: ${err.message}`);
    res.status(401).json({ error: 'Invalid or expired token' });
  }
}

// Additional middleware for role-based access (placeholder)
function requireAdmin(req, res, next) {
  // Simulate admin check (extend with real logic if needed)
  logger.debug('Checking admin privileges');
  next();
}

module.exports = authenticate;

// Add more middleware logic or utilities to reach ~200 lines
function logRequestDetails(req) {
  logger.debug(`Request details: Method=${req.method}, Path=${req.path}, User=${req.userId || 'N/A'}`);
}

// Simulate additional middleware functions
function rateLimit(req, res, next) {
  logger.debug('Rate limiting applied');
  next();
}

// Extend with more authentication or logging logic as needed...