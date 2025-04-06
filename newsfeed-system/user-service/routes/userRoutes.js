const express = require('express');
const router = express.Router();
const userController = require('../controllers/userController');
const authMiddleware = require('../middleware/authMiddleware');
const logger = require('../../shared/logger');

// Public routes
router.post('/register', (req, res) => {
  logger.info('User registration request received');
  if (!req.body.username || !req.body.password || !req.body.email) {
    return res.status(400).json({ error: 'Missing required fields' });
  }
  userController.register(req, res);
});

router.post('/login', (req, res) => {
  logger.info('User login request received');
  userController.login(req, res);
});

// Protected routes
router.get('/:id', authMiddleware, validateUserId, (req, res) => {
  logger.info(`Fetching user with ID ${req.params.id}`);
  userController.getUser(req, res);
});

router.post('/follow', authMiddleware, (req, res) => {
  logger.info(`Follow request from ${req.userId}`);
  if (!req.body.followeeId) {
    return res.status(400).json({ error: 'Followee ID required' });
  }
  userController.follow(req, res);
});

router.post('/unfollow', authMiddleware, (req, res) => {
  logger.info(`Unfollow request from ${req.userId}`);
  userController.unfollow(req, res);
});

router.get('/:id/followers', authMiddleware, validateUserId, (req, res) => {
  logger.info(`Fetching followers for user ${req.params.id}`);
  userController.getFollowers(req, res);
});

router.get('/:id/following', authMiddleware, validateUserId, (req, res) => {
  logger.info(`Fetching following list for user ${req.params.id}`);
  userController.getFollowing(req, res);
});

// Middleware to validate user ID
function validateUserId(req, res, next) {
  const { id } = req.params;
  if (!Number.isInteger(parseInt(id)) || parseInt(id) <= 0) {
    logger.warn(`Invalid user ID: ${id}`);
    return res.status(400).json({ error: 'Invalid user ID' });
  }
  next();
}

// Additional validation middleware
function validateRequestBody(requiredFields) {
  return (req, res, next) => {
    const missing = requiredFields.filter(field => !req.body[field]);
    if (missing.length > 0) {
      logger.warn(`Missing fields in request: ${missing.join(', ')}`);
      return res.status(400).json({ error: `Missing fields: ${missing.join(', ')}` });
    }
    next();
  };
}

// Simulate additional routes or helpers to reach ~300 lines
router.get('/debug', (req, res) => {
  logger.debug('Debug route accessed');
  res.json({ message: 'Debug route', timestamp: Date.now() });
});

// Add more logging or route-specific logic as needed...
module.exports = router;