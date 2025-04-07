import React, { useState, useEffect } from 'react';
import Message from './Message';
import { sendMessage, onMessage } from '../services/socket';

function ChatWindow() {
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState('');

  useEffect(() => {
    onMessage((msg) => setMessages((prev) => [...prev, { text: msg, sender: 'bot' }]));
  }, []);

  const handleSend = () => {
    if (input.trim()) {
      setMessages((prev) => [...prev, { text: input, sender: 'user' }]);
      sendMessage({ userId: localStorage.getItem('userId'), message: input });
      setInput('');
    }
  };

  return (
    <div className="chat-window">
      <div className="messages">
        {messages.map((msg, idx) => (
          <Message key={idx} text={msg.text} sender={msg.sender} />
        ))}
      </div>
      <input value={input} onChange={(e) => setInput(e.target.value)} />
      <button onClick={handleSend}>Send</button>
    </div>
  );
}

export default ChatWindow;