import React, { useEffect, useState } from 'react';
import { useMachineControl } from '../../hooks/useBackend';
import Card from '../../components/Card';
import StatusBadge from '../../components/StatusBadge';
import Button from '../../components/Button';

export const Dashboard: React.FC = () => {
  const { loadAllMachines, connectMachine, disconnectMachine, getMachineStatus } = useMachineControl();
  const [machines, setMachines] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [stats, setStats] = useState({
    totalMachines: 0,
    activeMachines: 0,
    totalTests: 247,
    todayTests: 18,
  });

  useEffect(() => {
    const initDashboard = async () => {
      try {
        const loadedMachines = await loadAllMachines();
        setMachines(loadedMachines);
        setStats(prev => ({
          ...prev,
          totalMachines: loadedMachines?.length || 0,
          activeMachines: loadedMachines?.filter((m: any) => m.connected).length || 0,
        }));
      } catch (error) {
        console.error('Failed to initialize dashboard:', error);
      } finally {
        setLoading(false);
      }
    };

    initDashboard();
  }, [loadAllMachines]);

  const handleConnect = async (machineId: string) => {
    try {
      await connectMachine(machineId);
      const loadedMachines = await loadAllMachines();
      setMachines(loadedMachines);
    } catch (error) {
      console.error('Connection failed:', error);
    }
  };

  const handleDisconnect = async (machineId: string) => {
    try {
      await disconnectMachine(machineId);
      const loadedMachines = await loadAllMachines();
      setMachines(loadedMachines);
    } catch (error) {
      console.error('Disconnection failed:', error);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-screen">
        <div className="text-xl text-gray-600">Loading machines...</div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <Card>
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Total Machines</p>
              <p className="text-3xl font-bold text-gray-900 mt-2">{stats.totalMachines}</p>
            </div>
            <div className="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center text-2xl">
              🔬
            </div>
          </div>
        </Card>

        <Card>
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Active Machines</p>
              <p className="text-3xl font-bold text-green-600 mt-2">{stats.activeMachines}</p>
            </div>
            <div className="w-12 h-12 bg-green-100 rounded-lg flex items-center justify-center text-2xl">
              ✓
            </div>
          </div>
        </Card>

        <Card>
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Total Tests</p>
              <p className="text-3xl font-bold text-gray-900 mt-2">{stats.totalTests}</p>
            </div>
            <div className="w-12 h-12 bg-purple-100 rounded-lg flex items-center justify-center text-2xl">
              📊
            </div>
          </div>
        </Card>

        <Card>
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Today's Tests</p>
              <p className="text-3xl font-bold text-blue-600 mt-2">{stats.todayTests}</p>
            </div>
            <div className="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center text-2xl">
              📅
            </div>
          </div>
        </Card>
      </div>

      {/* Lab Machines */}
      <Card title="Lab Machines" actions={
        <Button onClick={() => loadAllMachines()} variant="secondary" size="sm">
          🔄 Refresh
        </Button>
      }>
        {machines.length === 0 ? (
          <div className="text-center py-12">
            <p className="text-gray-500 mb-4">No machines configured</p>
            <Button variant="primary">+ Add Machine</Button>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {machines.map((machine) => (
              <div
                key={machine.id}
                className="border border-gray-200 rounded-lg p-4 hover:shadow-md transition-shadow"
              >
                <div className="flex items-start justify-between mb-3">
                  <div>
                    <h4 className="font-semibold text-gray-900">{machine.type}</h4>
                    <p className="text-sm text-gray-500">{machine.model || 'Model N/A'}</p>
                  </div>
                  <StatusBadge status={machine.connected ? 'connected' : 'disconnected'} />
                </div>

                <div className="space-y-2 mb-4">
                  <div className="flex justify-between text-sm">
                    <span className="text-gray-600">Port:</span>
                    <span className="font-medium">{machine.port}</span>
                  </div>
                  <div className="flex justify-between text-sm">
                    <span className="text-gray-600">Type:</span>
                    <span className="font-medium">{machine.connection_type}</span>
                  </div>
                </div>

                <div className="flex gap-2">
                  {machine.connected ? (
                    <Button
                      onClick={() => handleDisconnect(machine.id)}
                      variant="danger"
                      size="sm"
                      className="flex-1"
                    >
                      Disconnect
                    </Button>
                  ) : (
                    <Button
                      onClick={() => handleConnect(machine.id)}
                      variant="success"
                      size="sm"
                      className="flex-1"
                    >
                      Connect
                    </Button>
                  )}
                  <Button variant="secondary" size="sm">
                    ⚙️
                  </Button>
                </div>
              </div>
            ))}
          </div>
        )}
      </Card>

      {/* Recent Activity */}
      <Card title="Recent Activity">
        <div className="space-y-3">
          {[
            { time: '10 min ago', action: 'HA3100 completed CBC test for Pet #1234', type: 'success' },
            { time: '25 min ago', action: 'Scitek connected successfully', type: 'info' },
            { time: '1 hour ago', action: 'Critical result detected for Pet #1230', type: 'warning' },
            { time: '2 hours ago', action: 'Mindray disconnected', type: 'error' },
          ].map((activity, index) => (
            <div key={index} className="flex items-start gap-3 p-3 bg-gray-50 rounded-lg">
              <div className={`w-2 h-2 rounded-full mt-2 ${
                activity.type === 'success' ? 'bg-green-500' :
                activity.type === 'info' ? 'bg-blue-500' :
                activity.type === 'warning' ? 'bg-yellow-500' : 'bg-red-500'
              }`}></div>
              <div className="flex-1">
                <p className="text-sm text-gray-900">{activity.action}</p>
                <p className="text-xs text-gray-500 mt-1">{activity.time}</p>
              </div>
            </div>
          ))}
        </div>
      </Card>
    </div>
  );
};

export default Dashboard;

