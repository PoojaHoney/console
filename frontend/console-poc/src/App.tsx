import React from 'react';
import logo from './logo.svg';
import './App.css';
import CreateServiceAccount from './cloudProviders/gcp/iam';

function App() {
  return (
    <div>
      <CreateServiceAccount />
    </div>
  );
}

export default App;
