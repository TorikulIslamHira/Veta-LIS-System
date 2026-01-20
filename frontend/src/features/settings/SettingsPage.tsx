import React, { useState } from 'react';
import Card from '../../components/Card';
import Button from '../../components/Button';

const SettingsPage: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'general' | 'database' | 'machines' | 'notifications'>('general');

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-bold text-gray-900">Settings</h2>
        <p className="text-gray-600 mt-1">Configure system preferences and connections</p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
        {/* Settings Sidebar */}
        <div>
          <Card>
            <nav className="space-y-2">
              {[
                { id: 'general', label: 'General', icon: '⚙️' },
                { id: 'database', label: 'Database', icon: '🗄️' },
                { id: 'machines', label: 'Machines', icon: '🔬' },
                { id: 'notifications', label: 'Notifications', icon: '🔔' },
              ].map((tab) => (
                <button
                  key={tab.id}
                  onClick={() => setActiveTab(tab.id as any)}
                  className={`w-full flex items-center gap-3 px-4 py-2 rounded-lg transition-all text-left ${
                    activeTab === tab.id
                      ? 'bg-primary text-white'
                      : 'text-gray-700 hover:bg-gray-100'
                  }`}
                >
                  <span>{tab.icon}</span>
                  <span className="font-medium">{tab.label}</span>
                </button>
              ))}
            </nav>
          </Card>
        </div>

        {/* Settings Content */}
        <div className="lg:col-span-3">
          {activeTab === 'general' && (
            <Card title="General Settings">
              <div className="space-y-6">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Clinic Name
                  </label>
                  <input
                    type="text"
                    defaultValue="Veta Veterinary Clinic"
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Time Zone
                  </label>
                  <select className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent">
                    <option>UTC-5 (Eastern Time)</option>
                    <option>UTC-6 (Central Time)</option>
                    <option>UTC-7 (Mountain Time)</option>
                    <option>UTC-8 (Pacific Time)</option>
                  </select>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Language
                  </label>
                  <select className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent">
                    <option>English</option>
                    <option>Spanish</option>
                    <option>French</option>
                  </select>
                </div>

                <div className="flex items-center justify-between">
                  <div>
                    <p className="font-medium text-gray-900">Auto-connect machines on startup</p>
                    <p className="text-sm text-gray-500">Automatically connect to configured machines when app starts</p>
                  </div>
                  <button className="w-12 h-6 bg-primary rounded-full relative">
                    <span className="absolute right-1 top-1 w-4 h-4 bg-white rounded-full"></span>
                  </button>
                </div>

                <div className="pt-4">
                  <Button variant="primary">Save Changes</Button>
                </div>
              </div>
            </Card>
          )}

          {activeTab === 'database' && (
            <Card title="Database Configuration">
              <div className="space-y-6">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Database Host
                  </label>
                  <input
                    type="text"
                    defaultValue="localhost"
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Port
                  </label>
                  <input
                    type="text"
                    defaultValue="5432"
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Database Name
                  </label>
                  <input
                    type="text"
                    defaultValue="veta_lis"
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Username
                  </label>
                  <input
                    type="text"
                    defaultValue="postgres"
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Password
                  </label>
                  <input
                    type="password"
                    defaultValue="********"
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                  />
                </div>

                <div className="flex gap-3 pt-4">
                  <Button variant="secondary">Test Connection</Button>
                  <Button variant="primary">Save Configuration</Button>
                </div>
              </div>
            </Card>
          )}

          {activeTab === 'machines' && (
            <Card title="Machine Configuration">
              <div className="space-y-6">
                <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
                  <p className="text-sm text-blue-800">
                    ℹ️ Machine configurations are managed in the config.yaml file. 
                    You can add, edit, or remove machines from the Machines page.
                  </p>
                </div>

                <div>
                  <h4 className="font-semibold mb-3">Default Settings</h4>
                  
                  <div className="space-y-4">
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        Default Timeout (seconds)
                      </label>
                      <input
                        type="number"
                        defaultValue="30"
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                      />
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        Retry Attempts
                      </label>
                      <input
                        type="number"
                        defaultValue="3"
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                      />
                    </div>

                    <div className="flex items-center justify-between">
                      <div>
                        <p className="font-medium text-gray-900">Auto-reconnect on failure</p>
                        <p className="text-sm text-gray-500">Automatically attempt to reconnect if connection is lost</p>
                      </div>
                      <button className="w-12 h-6 bg-primary rounded-full relative">
                        <span className="absolute right-1 top-1 w-4 h-4 bg-white rounded-full"></span>
                      </button>
                    </div>
                  </div>
                </div>

                <div className="pt-4">
                  <Button variant="primary">Save Changes</Button>
                </div>
              </div>
            </Card>
          )}

          {activeTab === 'notifications' && (
            <Card title="Notification Settings">
              <div className="space-y-4">
                {[
                  { title: 'Critical Results', desc: 'Notify when critical test results are detected' },
                  { title: 'Machine Disconnection', desc: 'Alert when a machine loses connection' },
                  { title: 'New Test Results', desc: 'Notify when new test results arrive' },
                  { title: 'System Errors', desc: 'Alert on system or parsing errors' },
                ].map((notification, index) => (
                  <div key={index} className="flex items-center justify-between p-4 bg-gray-50 rounded-lg">
                    <div>
                      <p className="font-medium text-gray-900">{notification.title}</p>
                      <p className="text-sm text-gray-500">{notification.desc}</p>
                    </div>
                    <button className="w-12 h-6 bg-primary rounded-full relative">
                      <span className="absolute right-1 top-1 w-4 h-4 bg-white rounded-full"></span>
                    </button>
                  </div>
                ))}

                <div className="pt-4">
                  <Button variant="primary">Save Changes</Button>
                </div>
              </div>
            </Card>
          )}
        </div>
      </div>
    </div>
  );
};

export default SettingsPage;
