const express = require('express');
const bodyParser = require('body-parser');
const newsfeedRoutes = require('./routes/newsfeedRoutes');
const logger = require('../shared/logger');

const app = express();
app.use(bodyParser.json());
app.use('/newsfeed', newsfeedRoutes);

app.listen(3003, () => {
  logger.info('Newsfeed service running on port 3003');
});

// Add similar error handling and utilities as in user-service/server.js