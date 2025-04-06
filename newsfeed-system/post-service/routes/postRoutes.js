const express = require('express');
const router = express.Router();
const postController = require('../controllers/postController');
const authMiddleware = require('../middleware/authMiddleware');

router.post('/', authMiddleware, postController.createPost);
router.get('/:id', postController.getPost);
router.post('/:id/like', authMiddleware, postController.likePost);
router.post('/:id/comment', authMiddleware, postController.commentPost);

module.exports = router;

// Add validation and logging similar to userRoutes.js to reach ~300 lines