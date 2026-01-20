import React from 'react';

interface HeaderProps {
  currentPage: string;
}

const Header: React.FC<HeaderProps> = ({ currentPage }) => {
  const getPageTitle = () => {
    const titles: Record<string, string> = {
      dashboard: 'Dashboard',
      machines: 'Lab Machines',
      pets: 'Pets & Patients',
      results: 'Test Results',
      settings: 'Settings',
    };
    return titles[currentPage] || 'Dashboard';
  };

  return (
    <header className="bg-white border-b border-gray-200 px-6 py-4">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold text-gray-800">{getPageTitle()}</h2>
          <p className="text-sm text-gray-500 mt-1">
            {new Date().toLocaleDateString('en-US', { 
              weekday: 'long', 
              year: 'numeric', 
              month: 'long', 
              day: 'numeric' 
            })}
          </p>
        </div>
        
        <div className="flex items-center gap-4">
          <button className="relative p-2 text-gray-600 hover:bg-gray-100 rounded-lg">
            🔔
            <span className="absolute top-1 right-1 w-2 h-2 bg-red-500 rounded-full"></span>
          </button>
          
          <div className="flex items-center gap-3 px-4 py-2 bg-gray-50 rounded-lg">
            <div className="w-9 h-9 bg-primary text-white rounded-full flex items-center justify-center">
              👤
            </div>
            <div>
              <p className="text-sm font-medium text-gray-800">Dr. Admin</p>
              <p className="text-xs text-gray-500">Online</p>
            </div>
          </div>
        </div>
      </div>
    </header>
  );
};

export default Header;
