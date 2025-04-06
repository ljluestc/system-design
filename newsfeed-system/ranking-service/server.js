const express = require('express');
const logger = require('../shared/logger');

const app = express();

app.listen(3007, () => {
  logger.info('Ranking service running on port 3007');
});

// Add utilities to reach ~350 lines