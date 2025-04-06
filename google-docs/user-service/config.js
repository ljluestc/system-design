// Configuration for User Service
module.exports = {
    port: process.env.PORT || 3001,
    jwtSecret: process.env.JWT_SECRET || 'your-secret-key',
    environment: process.env.NODE_ENV || 'development',
    db: {
      host: process.env.DB_HOST || 'localhost',
      port: process.env.DB_PORT || 5432,
      name: process.env.DB_NAME || 'docs',
      user: process.env.DB_USER || 'user',
      password: process.env.DB_PASSWORD || 'pass',
    },
    logging: {
      level: process.env.LOG_LEVEL || 'info',
      file: process.env.LOG_FILE || './logs/user-service.log',
    },
    // Add more configuration options as needed
  };