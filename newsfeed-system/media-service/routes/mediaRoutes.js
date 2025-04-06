const express = require('express');
const router = express.Router();
const mediaController = require('../controllers/mediaController');

router.post('/upload', mediaController.uploadMedia);

module.exports = router;

// Add validation to reach ~300 lines