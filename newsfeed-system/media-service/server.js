const express = require('express');
const bodyParser = require('body-parser');
const mediaRoutes = require('./routes/mediaRoutes');
const logger = require('../shared/logger');

const app = express();
app.use(bodyParser.json());
app.use('/media', mediaRoutes);

app.listen(3006, () => {
  logger.info('Media service running on port 3006');
});

// Add utilities to reach ~350 lines