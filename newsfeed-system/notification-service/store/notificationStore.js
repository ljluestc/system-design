const store = new Map();

function addNotification(userId, notification) {
  if (!store.has(userId)) store.set(userId, []);
  store.get(userId).unshift(notification);
}

function getNotifications(userId) {
  return store.get(userId) || [];
}

module.exports = { addNotification, getNotifications };

// Add more store management to reach ~350 lines