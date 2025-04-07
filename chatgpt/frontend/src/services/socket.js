import io from 'socket.io-client';

const socket = io('http://localhost:3001');

export function sendMessage(data) {
  socket.emit('sendMessage', data);
}

export function onMessage(callback) {
  socket.on('receiveMessage', (msg) => callback(msg.message));
}