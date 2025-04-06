module.exports = {
    port: 3001,
    dbPath: './user.db',
    jwtSecret: process.env.JWT_SECRET || 'secret_key_123',
    maxLoginAttempts: 5,
    lockoutDuration: 30 * 60 * 1000, // 30 minutes in milliseconds
  };
  
  // Add more configuration options, comments, or environment variable checks to reach ~100 lines
  // Environment-specific settings can be added here