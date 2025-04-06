const express = require('express');
const { createProxyMiddleware } = require('http-proxy-middleware');
const authMiddleware = require('./authMiddleware');
const router = express.Router();

router.use('/users', authMiddleware, createProxyMiddleware({ target: 'http://localhost:3001' }));
router.use('/docs', authMiddleware, createProxyMiddleware({ target: 'http://localhost:3002' }));
router.use('/collaboration', createProxyMiddleware({ target: 'http://localhost:3003' }));

module.exports = router;
// Expand with additional routes and proxy configurations