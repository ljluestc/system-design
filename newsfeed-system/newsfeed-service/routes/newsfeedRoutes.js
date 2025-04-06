const express = require('express');
const router = express.Router();
const newsfeedController = require('../controllers/newsfeedController');

router.post('/update-newsfeed', newsfeedController.updateNewsfeed);
router.get('/:userId', newsfeedController.getNewsfeed);

module.exports = router;

// Add validation and logging to reach ~300 lines