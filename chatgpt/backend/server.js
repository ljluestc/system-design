const express = require('express');
const http = require('http');
const socketIo = require('socket.io');
const cors = require('cors');
const Queue = require('bull');
const userRoutes = require('./routes/userRoutes');
const chatRoutes = require('./routes/chatRoutes');
const { connectToDatabase } = require('./config/database');
const logger = require('./utils/logger');
const aiService = require('./services/aiService');

const app = express();
const server = http.createServer(app);
const io = socketIo(server, { cors: { origin: '*' } });
const aiQueue = new Queue('aiQueue');

app.use(cors());
app.use(express.json());
app.use('/api/users', userRoutes);
app.use('/api/chat', chatRoutes);

io.on('connection', (socket) => {
  logger.log('New client connected');
  socket.on('sendMessage', (data) => {
    aiQueue.add(data);
  });
  socket.on('disconnect', () => logger.log('Client disconnected'));
});

aiQueue.process(async (job) => {
  const { userId, message } = job.data;
  const response = await aiService.generateResponse(message);
  io.to(userId).emit('receiveMessage', { message: response });
});

connectToDatabase().then(() => {
  server.listen(3001, () => logger.log('Server running on port 3001'));
});