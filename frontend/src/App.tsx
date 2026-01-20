import React from 'react';

export const App: React.FC = () => {
  return (
    <div className="min-h-screen bg-gray-50">
      <nav className="bg-blue-600 text-white shadow-lg">
        <div className="max-w-7xl mx-auto px-6 py-4">
          <h1 className="text-2xl font-bold">Veta LIS</h1>
          <p className="text-blue-100">Veterinary Laboratory Information System</p>
        </div>
      </nav>
      
      <main className="max-w-7xl mx-auto">
        <div className="p-6 text-center">
          <h2 className="text-3xl font-bold mb-4">Welcome to Veta LIS</h2>
          <p className="text-gray-600 mb-8">
            A production-grade laboratory information system for veterinary clinics.
          </p>
          <div className="grid grid-cols-3 gap-4">
            <div className="bg-white p-6 rounded-lg shadow">
              <h3 className="font-semibold text-lg mb-2">Dashboard</h3>
              <p className="text-gray-600">Monitor lab machines in real-time</p>
            </div>
            <div className="bg-white p-6 rounded-lg shadow">
              <h3 className="font-semibold text-lg mb-2">Settings</h3>
              <p className="text-gray-600">Configure machine connections</p>
            </div>
            <div className="bg-white p-6 rounded-lg shadow">
              <h3 className="font-semibold text-lg mb-2">Results</h3>
              <p className="text-gray-600">View and manage test results</p>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
};

export default App;
