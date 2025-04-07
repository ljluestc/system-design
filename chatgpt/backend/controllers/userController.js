const bcrypt = require('bcrypt');
const jwt = require('jsonwebtoken');
const User = require('../models/User');
const config = require('../config/app');

async function register(req, res) {
  const { username, password } = req.body;
  try {
    const hashedPassword = await bcrypt.hash(password, 10);
    const user = new User({ username, password: hashedPassword });
    await user.save();
    res.status(201).json({ message: 'User registered' });
  } catch (error) {
    res.status(500).json({ error: 'Failed to register user' });
  }
}

async function login(req, res) {
  const { username, password } = req.body;
  try {
    const user = await User.findOne({ username });
    if (!user || !(await bcrypt.compare(password, user.password))) {
      return res.status(401).json({ error: 'Invalid credentials' });
    }
    const token = jwt.sign({ userId: user._id }, config.jwtSecret, { expiresIn: '1h' });
    res.json({ token, userId: user._id });
  } catch (error) {
    res.status(500).json({ error: 'Failed to login' });
  }
}

module.exports = { register, login };