const sqlite = require('sqlite');
const sqlite3 = require('sqlite3');
const logger = require('../../shared/logger');

let db;

async function initDb() {
  db = await sqlite.open({ filename: './post.db', driver: sqlite3.Database });
  await db.exec(`
    CREATE TABLE IF NOT EXISTS posts (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      user_id INTEGER,
      content TEXT NOT NULL,
      media TEXT,
      timestamp INTEGER,
      likes INTEGER DEFAULT 0,
      comments INTEGER DEFAULT 0
    );
  `);
  logger.info('Post database initialized');
}

async function createPost(userId, content, media) {
  const result = await db.run(
    'INSERT INTO posts (user_id, content, media, timestamp) VALUES (?, ?, ?, ?)',
    [userId, content, JSON.stringify(media), Date.now()]
  );
  return result.lastID;
}

async function getPost(id) {
  const post = await db.get('SELECT * FROM posts WHERE id = ?', [id]);
  if (post) post.media = JSON.parse(post.media);
  return post;
}

async function likePost(id, userId) {
  await db.run('UPDATE posts SET likes = likes + 1 WHERE id = ?', [id]);
}

async function commentPost(id, userId, content) {
  await db.run('UPDATE posts SET comments = comments + 1 WHERE id = ?', [id]);
}

module.exports = { initDb, createPost, getPost, likePost, commentPost };

// Extend with more operations to reach ~350 lines