const mongoose = require('mongoose');

const messageSchema = new mongoose.Schema({
  userId: { type: mongoose.Schema.Types.ObjectId, ref: 'User' },
  message: String,
  sender: String,
  timestamp: { type: Date, default: Date.now },
});

module.exports = mongoose.model('Message', messageSchema);