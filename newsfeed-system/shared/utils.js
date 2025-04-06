const jwt = require('jsonwebtoken');

function generateJWT(userId) {
  return jwt.sign({ userId }, 'secret_key_123', { expiresIn: '1h' });
}

function verifyJWT(token) {
  try {
    return jwt.verify(token, 'secret_key_123');
  } catch (err) {
    return null;
  }
}

module.exports = { generateJWT, verifyJWT };

// Add more utilities to reach ~200 lines