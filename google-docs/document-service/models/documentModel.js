const db = require('../../shared/db');
const logger = require('../../shared/logger');

async function createDocument(title, ownerId) {
  try {
    const [id] = await db('documents').insert({ title, ownerId });
    return id;
  } catch (error) {
    logger.log(`Error creating document: ${error.message}`);
    throw error;
  }
}

module.exports = { createDocument };