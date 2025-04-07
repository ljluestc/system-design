const express = require('express');
const { register, login } = require('../controllers/userController');
const auth = require('../middleware/auth');

const router = express.Router();

router.post('/register', register);
router.post('/login', login);
router.get('/auth', auth, (req, res) => res.json({ authenticated: true }));

module.exports = router;