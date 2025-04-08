import io from "socket.io-client";

const socket = io("http://localhost:3001");

export function onUpdate(callback) {
  socket.on("update", callback);
}