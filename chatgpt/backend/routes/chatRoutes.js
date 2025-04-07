const express = require('express');
const { getMessages, saveMessage } = require('../controllers/chatController');
const auth = require('../middleware/auth');

const router = express.Router();

router.get('/messages', auth, getMessages);
router.post('/messages', auth, saveMessage);

module.exports = router;