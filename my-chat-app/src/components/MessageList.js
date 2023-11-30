// src/components/MessageList.js
import React, { useState, useEffect } from 'react';
import { w3cwebsocket as W3CWebSocket } from 'websocket';

const ws = new W3CWebSocket('ws://localhost:8080/ws');

function MessageList() {
  const [messages, setMessages] = useState([]);

  useEffect(() => {
    // Fetch initial messages from the server
    fetchMessages();

    // WebSocket event listeners
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

  const fetchMessages = async () => {
    try {
      const response = await fetch('http://localhost:8080/api/messages');
      const data = await response.json();
      setMessages(data);
    } catch (error) {
      console.error('Error fetching messages:', error);
    }
  };

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
