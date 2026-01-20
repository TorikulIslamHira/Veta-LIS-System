import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App';
import './index.css';
import { MachineProvider } from './context/MachineContext';

ReactDOM.createRoot(document.getElementById('app')!).render(
  <React.StrictMode>
    <MachineProvider>
      <App />
    </MachineProvider>
  </React.StrictMode>,
);
