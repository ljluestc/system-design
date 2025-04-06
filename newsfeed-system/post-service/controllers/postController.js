const postModel = require('../models/postModel');
const axios = require('axios');
const logger = require('../../shared/logger');

async function createPost(req, res) {
  const { content, media } = req.body;
  const userId = req.userId;
  try {
    if (!content || content.length > 280) {
      throw new Error('Content must be non-empty and less than 280 characters');
    }
    const postId = await postModel.createPost(userId, content, media || []);
    const timestamp = Date.now();
    await axios.post('http://localhost:3003/update-newsfeed', { userId, postId, content, timestamp });
    await axios.post('http://localhost:3004/index-post', { postId, content });
    logger.info(`Post ${postId} created by user ${userId}`);
    res.status(201).json({ postId });
  } catch (err) {
    logger.error(`Create post failed: ${err.message}`);
    res.status(400).json({ error: err.message });
  }
}

async function getPost(req, res) {
  const { id } = req.params;
  try {
    const post = await postModel.getPost(id);
    if (!post) throw new Error('Post not found');
    res.json(post);
  } catch (err) {
    res.status(404).json({ error: err.message });
  }
}

async function likePost(req, res) {
  const { id } = req.params;
  const userId = req.userId;
  try {
    await postModel.likePost(id, userId);
    res.status(200).json({ message: 'Post liked' });
  } catch (err) {
    res.status(400).json({ error: err.message });
  }
}

async function commentPost(req, res) {
  const { id } = req.params;
  const { content } = req.body;
  const userId = req.userId;
  try {
    await postModel.commentPost(id, userId, content);
    res.status(200).json({ message: 'Comment added' });
  } catch (err) {
    res.status(400).json({ error: err.message });
  }
}

module.exports = { createPost, getPost, likePost, commentPost };

// Add more functions or detailed logic to reach ~400 lines