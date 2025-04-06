const express = require('express');
const bodyParser = require('body-parser');
const searchRoutes = require('./routes/searchRoutes');
const logger = require('../shared/logger');

const app = express();
app.use(bodyParser.json());
app.use('/search', searchRoutes);

app.listen(3004, () => {
  logger.info('Search service running on port 3004');
});

// Add similar utilities to reach ~350 lines