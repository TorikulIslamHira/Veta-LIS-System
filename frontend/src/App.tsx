import React, { useState } from 'react';
import { MachineProvider } from './context/MachineContext';
import Sidebar from './components/Sidebar';
import Header from './components/Header';
import Dashboard from './features/dashboard/Dashboard';
import MachinesPage from './features/machines/MachinesPage';
import PetsPage from './features/pets/PetsPage';
import ResultsPage from './features/results/ResultsPage';
import SettingsPage from './features/settings/SettingsPage';

export const App: React.FC = () => {
  const [currentPage, setCurrentPage] = useState('dashboard');

  const renderPage = () => {
    switch (currentPage) {
      case 'dashboard':
        return <Dashboard />;
      case 'machines':
        return <MachinesPage />;
      case 'pets':
        return <PetsPage />;
      case 'results':
        return <ResultsPage />;
      case 'settings':
        return <SettingsPage />;
      default:
        return <Dashboard />;
    }
  };

  return (
    <MachineProvider>
      <div className="flex h-screen bg-gray-50">
        <Sidebar currentPage={currentPage} onNavigate={setCurrentPage} />
        <div className="flex-1 flex flex-col overflow-hidden">
          <Header currentPage={currentPage} />
          <main className="flex-1 overflow-y-auto bg-gray-50 p-6">
            {renderPage()}
          </main>
        </div>
      </div>
    </MachineProvider>
  );
};

export default App;
