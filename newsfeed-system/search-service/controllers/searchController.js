const express = require('express');
const router = express.Router();
const searchController = require('../controllers/searchController');

router.post('/index-post', searchController.indexPost);
router.get('/', searchController.searchPosts);

module.exports = router;

// Add validation to reach ~300 lines