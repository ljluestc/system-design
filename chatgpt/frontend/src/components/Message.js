import React from 'react';

function Message({ text, sender }) {
  return <div className={`message ${sender}`}>{text}</div>;
}

export default Message;