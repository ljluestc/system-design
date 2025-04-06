const express = require('express');
const documentController = require('./controllers/documentController');
const logger = require('../shared/logger');

const router = express.Router();

router.post('/create', documentController.createDocument);
// Add more routes as needed