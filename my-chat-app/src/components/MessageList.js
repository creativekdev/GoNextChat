// src/components/MessageList.js
import React, { useState, useEffect } from 'react';

function MessageList() {
  const [messages, setMessages] = useState([]);

  useEffect(() => {
    // Fetch messages from the backend
    fetch('http://localhost:8080/api/messages') // Update the URL with your Golang backend URL
      .then(response => response.json())
      .then(data => setMessages(data))
      .catch(error => console.error('Error fetching messages:', error));
  }, []);

  return (
    <div className="message-container">
      {messages && messages.map(message => (
        <div key={message.id} className="message">
          {message.content}
        </div>
      ))}
    </div>
  );
}

export default MessageList;
