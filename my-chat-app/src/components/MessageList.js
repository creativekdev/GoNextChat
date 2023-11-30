// src/components/MessageList.js
import React, { useState, useEffect } from 'react';
import { w3cwebsocket as W3CWebSocket } from 'websocket';

const ws = new W3CWebSocket('ws://localhost:8080/ws');

function MessageList() {
  const [messages, setMessages] = useState([]);

  useEffect(() => {
    ws.onmessage = (event) => {
      const message = JSON.parse(event.data);
      setMessages((prevMessages) => [...prevMessages, message]);
    };

    ws.onclose = () => {
      console.log('WebSocket connection closed');
    };

    return () => {
      ws.close();
    };
  }, []);

  return (
    <div className="message-container">
      {messages.map((message) => (
        <div key={message.id} className="message">
          {message.content}
        </div>
      ))}
    </div>
  );
}

export default MessageList;
