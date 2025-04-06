const notificationStore = require('../store/notificationStore');
const logger = require('../../shared/logger');

async function notify(req, res) {
  const { userId, type, fromUserId, postId } = req.body;
  try {
    notificationStore.addNotification(userId, { type, fromUserId, postId, timestamp: Date.now() });
    res.sendStatus(200);
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
}

async function getNotifications(req, res) {
  const { userId } = req.params;
  try {
    const notifications = notificationStore.getNotifications(userId);
    res.json(notifications);
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
}

module.exports = { notify, getNotifications };

// Add more logic to reach ~400 lines