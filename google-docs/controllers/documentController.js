const documentModel = require('../models/documentModel');

async function createDocument(req, res) {
  const { title, ownerId } = req.body;
  try {
    const docId = await documentModel.createDocument(title, ownerId);
    res.status(201).json({ docId });
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
}

module.exports = { createDocument };