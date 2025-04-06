const jwt = require('jsonwebtoken');
const config = require('../config');

function authenticate(req, res, next) {
  // Extract token from Authorization header
  const token = req.headers['authorization'];
  if (!token) {
    return res.status(401).json({ error: 'No token provided' });
  }

  try {
    // Verify token using the secret from config
    const decoded = jwt.verify(token, config.jwtSecret);
    req.userId = decoded.userId; // Attach userId to request object
    next(); // Proceed to the next middleware or route handler
  } catch (err) {
    return res.status(401).json({ error: 'Invalid token' });
  }
}

module.exports = authenticate;