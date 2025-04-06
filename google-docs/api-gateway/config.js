module.exports = {
    port: process.env.PORT || 3000,
    jwtSecret: process.env.JWT_SECRET || 'your-secret-key',
    serviceUrls: {
      user: 'http://localhost:3001',
      document: 'http://localhost:3002',
    },
  };
  // Add more environment variables and configurations