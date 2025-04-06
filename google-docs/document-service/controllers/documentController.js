const documentModel = require('../models/documentModel');
const logger = require('../../shared/logger');

async function createDocument(req, res) {
  const { title, ownerId } = req.body;
  try {
    if (!title || !ownerId) {
      throw new Error('Title and ownerId are required');
    }
    const docId = await documentModel.createDocument(title, ownerId);
    res.status(201).json({ docId });
  } catch (error) {
    logger.log(`Error creating document: ${error.message}`);
    res.status(500).json({ error: error.message });
  }
}

module.exports = { createDocument };