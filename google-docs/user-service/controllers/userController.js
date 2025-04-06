const userModel = require('../models/userModel');
const jwt = require('jsonwebtoken');
const config = require('../config');
const logger = require('../../shared/logger');
const { STATUS_OK } = require('../../shared/constants');
const { sanitizeInput } = require('../../shared/utils');

// Register a new user
async function register(req, res) {
  const { username, password } = req.body;
  logger.log(`Register attempt for username: ${username}`);

  try {
    // Input validation
    if (!username || !password) {
      logger.log('Registration failed: Missing username or password');
      return res.status(400).json({ error: 'Username and password are required' });
    }

    if (username.length < 3) {
      logger.log('Registration failed: Username too short');
      return res.status(400).json({ error: 'Username must be at least 3 characters' });
    }

    if (password.length < 6) {
      logger.log('Registration failed: Password too short');
      return res.status(400).json({ error: 'Password must be at least 6 characters' });
    }

    // Sanitize inputs
    const sanitizedUsername = sanitizeInput(username);
    const sanitizedPassword = sanitizeInput(password);

    // Check if user already exists
    const existingUser = await userModel.getUserByUsername(sanitizedUsername);
    if (existingUser) {
      logger.log(`Registration failed: Username ${sanitizedUsername} already exists`);
      return res.status(409).json({ error: 'Username already exists' });
    }

    // Create user
    const userId = await userModel.createUser(sanitizedUsername, sanitizedPassword);
    logger.log(`User created with ID: ${userId}`);

    // Generate JWT token
    const token = jwt.sign({ userId }, config.jwtSecret, { expiresIn: '1h' });
    logger.log(`Token generated for user ID: ${userId}`);

    res.status(201).json({
      userId,
      token,
      message: 'User registered successfully',
      timestamp: new Date().toISOString(),
    });
  } catch (error) {
    logger.log(`Registration error: ${error.message}`);
    res.status(500).json({
      error: 'Failed to register user',
      message: error.message,
      timestamp: new Date().toISOString(),
    });
  }
}

// Login a user
async function login(req, res) {
  const { username, password } = req.body;
  logger.log(`Login attempt for username: ${username}`);

  try {
    // Input validation
    if (!username || !password) {
      logger.log('Login failed: Missing username or password');
      return res.status(400).json({ error: 'Username and password are required' });
    }

    // Sanitize inputs
    const sanitizedUsername = sanitizeInput(username);
    const sanitizedPassword = sanitizeInput(password);

    // Fetch user
    const user = await userModel.getUserByUsername(sanitizedUsername);
    if (!user) {
      logger.log(`Login failed: Username ${sanitizedUsername} not found`);
      return res.status(401).json({ error: 'Invalid credentials' });
    }

    // Verify password (in a real app, use hashing)
    if (user.password !== sanitizedPassword) {
      logger.log(`Login failed: Incorrect password for ${sanitizedUsername}`);
      return res.status(401).json({ error: 'Invalid credentials' });
    }

    // Generate JWT token
    const token = jwt.sign({ userId: user.id }, config.jwtSecret, { expiresIn: '1h' });
    logger.log(`Login successful for user ID: ${user.id}`);

    res.status(STATUS_OK).json({
      userId: user.id,
      token,
      message: 'Login successful',
      timestamp: new Date().toISOString(),
    });
  } catch (error) {
    logger.log(`Login error: ${error.message}`);
    res.status(500).json({
      error: 'Failed to login',
      message: error.message,
      timestamp: new Date().toISOString(),
    });
  }
}

// Placeholder for additional functions
async function updateUser(req, res) {
  res.status(STATUS_OK).json({ message: 'Update user not implemented' });
}

async function deleteUser(req, res) {
  res.status(STATUS_OK).json({ message: 'Delete user not implemented' });
}

module.exports = { register, login, updateUser, deleteUser };