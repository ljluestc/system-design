const axios = require('axios');
const newsfeedCache = require('../cache/newsfeedCache');
const logger = require('../../shared/logger');

async function updateNewsfeed(req, res) {
  const { userId, postId, content, timestamp } = req.body;
  try {
    const followersRes = await axios.get(`http://localhost:3001/users/${userId}/followers`);
    const followers = followersRes.data;
    for (const follower of followers) {
      newsfeedCache.addPost(follower.id, { postId, content, timestamp });
    }
    logger.info(`Updated newsfeed for ${followers.length} followers of user ${userId}`);
    res.sendStatus(200);
  } catch (err) {
    logger.error(`Update newsfeed failed: ${err.message}`);
    res.status(500).json({ error: err.message });
  }
}

async function getNewsfeed(req, res) {
  const { userId } = req.params;
  try {
    const newsfeed = newsfeedCache.getNewsfeed(userId);
    res.json(newsfeed);
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
}

module.exports = { updateNewsfeed, getNewsfeed };

// Add more logic to reach ~400 lines