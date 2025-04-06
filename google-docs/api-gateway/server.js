const express = require('express');
const routes = require('./routes');
const config = require('./config');

const app = express();
app.use(express.json());
app.use('/api', routes);

app.listen(config.port, () => {
  console.log(`API Gateway running on port ${config.port}`);
});

// Add error handling, middleware, and additional routes to expand to 350 lines