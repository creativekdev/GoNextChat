// src/pages/index.js
import React from 'react';
import MessageList from '../components/MessageList';
import MessageForm from '../components/MessageForm';

function Home() {
  return (
    <div className="App">
      <h1>Real-time Chat</h1>
      <MessageList />
      <MessageForm />
    </div>
  );
}

export default Home;
