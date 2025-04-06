const express = require('express');
const bodyParser = require('body-parser');
const notificationRoutes = require('./routes/notificationRoutes');
const logger = require('../shared/logger');

const app = express();
app.use(bodyParser.json());
app.use('/notifications', notificationRoutes);

app.listen(3005, () => {
  logger.info('Notification service running on port 3005');
});

// Add utilities to reach ~350 lines