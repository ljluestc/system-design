const db = require('../../shared/db');
const logger = require('../../shared/logger');

// Create a new user
async function createUser(username, password) {
  try {
    logger.log(`Creating user with username: ${username}`);
    const [id] = await db('users').insert({
      username,
      password, // In production, hash this!
      created_at: new Date(),
      updated_at: new Date(),
    });
    logger.log(`User created with ID: ${id}`);
    return id;
  } catch (error) {
    logger.log(`Error creating user: ${error.message}`);
    throw new Error(`Database error: ${error.message}`);
  }
}

// Get user by username
async function getUserByUsername(username) {
  try {
    logger.log(`Fetching user with username: ${username}`);
    const user = await db('users')
      .where({ username })
      .first();
    if (!user) {
      logger.log(`User not found: ${username}`);
    } else {
      logger.log(`User found: ${username}, ID: ${user.id}`);
    }
    return user;
  } catch (error) {
    logger.log(`Error fetching user: ${error.message}`);
    throw new Error(`Database error: ${error.message}`);
  }
}

// Placeholder for additional database operations
async function updateUser(userId, updates) {
  try {
    logger.log(`Updating user ID: ${userId}`);
    const updatedCount = await db('users')
      .where({ id: userId })
      .update({ ...updates, updated_at: new Date() });
    logger.log(`Updated ${updatedCount} user(s)`);
    return updatedCount;
  } catch (error) {
    logger.log(`Error updating user: ${error.message}`);
    throw error;
  }
}

async function deleteUser(userId) {
  try {
    logger.log(`Deleting user ID: ${userId}`);
    const deletedCount = await db('users')
      .where({ id: userId })
      .del();
    logger.log(`Deleted ${deletedCount} user(s)`);
    return deletedCount;
  } catch (error) {
    logger.log(`Error deleting user: ${error.message}`);
    throw error;
  }
}

module.exports = { createUser, getUserByUsername, updateUser, deleteUser };