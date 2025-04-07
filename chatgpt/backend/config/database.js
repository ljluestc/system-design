const mongoose = require('mongoose');

async function connectToDatabase() {
  try {
    await mongoose.connect('mongodb://localhost/chatgpt', {
      useNewUrlParser: true,
      useUnifiedTopology: true,
    });
    console.log('Connected to MongoDB');
  } catch (error) {
    console.error('Database connection error:', error);
  }
}

module.exports = { connectToDatabase };