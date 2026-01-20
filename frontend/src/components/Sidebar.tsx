import React from 'react';

interface SidebarProps {
  currentPage: string;
  onNavigate: (page: string) => void;
}

const Sidebar: React.FC<SidebarProps> = ({ currentPage, onNavigate }) => {
  const menuItems = [
    { id: 'dashboard', label: 'Dashboard', icon: '📊' },
    { id: 'machines', label: 'Machines', icon: '🔬' },
    { id: 'pets', label: 'Pets', icon: '🐕' },
    { id: 'results', label: 'Results', icon: '📋' },
    { id: 'settings', label: 'Settings', icon: '⚙️' },
  ];

  return (
    <div className="w-64 bg-primary text-white flex flex-col">
      <div className="p-6 border-b border-blue-800">
        <h1 className="text-2xl font-bold">🏥 Veta LIS</h1>
        <p className="text-sm text-blue-200 mt-1">Veterinary Lab System</p>
      </div>
      
      <nav className="flex-1 p-4">
        <ul className="space-y-2">
          {menuItems.map((item) => (
            <li key={item.id}>
              <button
                onClick={() => onNavigate(item.id)}
                className={`w-full flex items-center gap-3 px-4 py-3 rounded-lg transition-all ${
                  currentPage === item.id
                    ? 'bg-white text-primary shadow-md'
                    : 'text-blue-100 hover:bg-blue-800'
                }`}
              >
                <span className="text-xl">{item.icon}</span>
                <span className="font-medium">{item.label}</span>
              </button>
            </li>
          ))}
        </ul>
      </nav>
      
      <div className="p-4 border-t border-blue-800">
        <div className="flex items-center gap-3 px-4 py-2">
          <div className="w-10 h-10 bg-blue-800 rounded-full flex items-center justify-center">
            👤
          </div>
          <div>
            <p className="font-medium">Dr. Admin</p>
            <p className="text-xs text-blue-200">Veterinarian</p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Sidebar;