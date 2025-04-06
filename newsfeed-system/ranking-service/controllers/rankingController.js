const postModel = require('../models/postModel');
const logger = require('../../shared/logger');

async function rankPosts() {
  try {
    const posts = await postModel.getAllPosts();
    posts.sort((a, b) => (b.likes + b.comments) - (a.likes + a.comments));
    return posts.slice(0, 10);
  } catch (err) {
    logger.error(`Rank posts failed: ${err.message}`);
    throw err;
  }
}

module.exports = { rankPosts };

// Add more logic to reach ~400 lines