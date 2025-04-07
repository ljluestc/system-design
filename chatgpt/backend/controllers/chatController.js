const Message = require('../models/Message');

async function getMessages(req, res) {
  try {
    const messages = await Message.find({ userId: req.user.userId });
    res.json(messages);
  } catch (error) {
    res.status(500).json({ error: 'Failed to retrieve messages' });
  }
}

async function saveMessage(req, res) {
  const { message, sender } = req.body;
  try {
    const newMessage = new Message({ userId: req.user.userId, message, sender });
    await newMessage.save();
    res.status(201).json({ message: 'Message saved' });
  } catch (error) {
    res.status(500).json({ error: 'Failed to save message' });
  }
}

module.exports = { getMessages, saveMessage };