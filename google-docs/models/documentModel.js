const db = require('../../shared/db');

async function createDocument(title, ownerId) {
  const [id] = await db('documents').insert({ title, ownerId });
  return id;
}

module.exports = { createDocument };