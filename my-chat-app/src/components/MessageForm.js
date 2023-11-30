// src/components/MessageForm.js
import React, { useState } from 'react';

function MessageForm() {
  const [newMessage, setNewMessage] = useState('');

  const handleInputChange = (e) => {
    setNewMessage(e.target.value);
  };

  const handleSubmit = (e) => {
    e.preventDefault();

    // Post new message to the backend
    fetch('http://localhost:8080/api/messages', { // Update the URL with your Golang backend URL
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ content: newMessage }),
    })
      .then(response => response.json())
      .then(data => {
        // Handle the response if needed
        console.log(data);
      })
      .catch(error => console.error('Error posting message:', error));

    setNewMessage('');
  };

  return (
    <form onSubmit={handleSubmit}>
      <input
        type="text"
        value={newMessage}
        onChange={handleInputChange}
        placeholder="Type your message..."
      />
      <button type="submit">Send</button>
    </form>
  );
}

export default MessageForm;
