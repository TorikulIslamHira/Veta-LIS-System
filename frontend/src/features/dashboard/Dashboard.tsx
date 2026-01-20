import React, { useEffect, useState } from 'react';
import { useMachineControl } from '../../hooks/useBackend';

export const Dashboard: React.FC = () => {
  const { loadAllMachines, getMachineStatus } = useMachineControl();
  const [machines, setMachines] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const initDashboard = async () => {
      try {
        const loadedMachines = await loadAllMachines();
        setMachines(loadedMachines);
      } catch (error) {
        console.error('Failed to initialize dashboard:', error);
      } finally {
        setLoading(false);
      }
    };

    initDashboard();
  }, [loadAllMachines]);

  if (loading) {
    return (
      <div className="flex items-center justify-center h-screen">
        <div className="text-xl text-gray-600">Loading machines...</div>
      </div>
    );
  }

  return (
    <div className="p-6">
      <h1 className="text-3xl font-bold mb-6">Veta LIS Dashboard</h1>
      
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {machines.map((machine) => (
          <div key={machine.id} className="bg-white rounded-lg shadow p-4">
            <h2 className="text-lg font-semibold mb-2">{machine.type}</h2>
            <div className="space-y-2 text-sm">
              <p><strong>ID:</strong> {machine.id}</p>
              <p><strong>Port:</strong> {machine.port}</p>
              <p>
                <strong>Status:</strong>{' '}
                <span className={machine.connected ? 'text-green-600' : 'text-red-600'}>
                  {machine.connected ? 'Connected' : 'Disconnected'}
                </span>
              </p>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};
