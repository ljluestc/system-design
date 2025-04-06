const express = require('express');
const router = express.Router();
const notificationController = require('../controllers/notificationController');

router.post('/', notificationController.notify);
router.get('/:userId', notificationController.getNotifications);

module.exports = router;

// Add validation to reach ~300 lines