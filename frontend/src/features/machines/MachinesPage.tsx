import React, { useEffect, useState } from 'react';
import { useMachineControl } from '../../hooks/useBackend';
import Card from '../../components/Card';
import Button from '../../components/Button';
import StatusBadge from '../../components/StatusBadge';

const MachinesPage: React.FC = () => {
  const { loadAllMachines, connectMachine, disconnectMachine } = useMachineControl();
  const [machines, setMachines] = useState<any[]>([]);
  const [selectedMachine, setSelectedMachine] = useState<any>(null);

  useEffect(() => {
    loadMachines();
  }, []);

  const loadMachines = async () => {
    try {
      const data = await loadAllMachines();
      setMachines(data || []);
    } catch (error) {
      console.error('Failed to load machines:', error);
    }
  };

  const handleConnect = async (machineId: string) => {
    try {
      await connectMachine(machineId);
      await loadMachines();
    } catch (error) {
      console.error('Connection failed:', error);
    }
  };

  const handleDisconnect = async (machineId: string) => {
    try {
      await disconnectMachine(machineId);
      await loadMachines();
    } catch (error) {
      console.error('Disconnection failed:', error);
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-2xl font-bold text-gray-900">Lab Machines</h2>
          <p className="text-gray-600 mt-1">Manage and monitor connected laboratory equipment</p>
        </div>
        <Button onClick={() => alert('Add machine feature coming soon')}>
          + Add New Machine
        </Button>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Machines List */}
        <div className="lg:col-span-2">
          <Card>
            <div className="space-y-4">
              {machines.length === 0 ? (
                <div className="text-center py-12">
                  <div className="text-6xl mb-4">🔬</div>
                  <p className="text-gray-500 mb-4">No machines configured yet</p>
                  <Button onClick={() => alert('Add machine feature coming soon')}>Add Your First Machine</Button>
                </div>
              ) : (
                machines.map((machine) => (
                  <div
                    key={machine.id}
                    onClick={() => setSelectedMachine(machine)}
                    className={`p-4 border rounded-lg cursor-pointer transition-all ${
                      selectedMachine?.id === machine.id
                        ? 'border-primary bg-blue-50'
                        : 'border-gray-200 hover:border-gray-300'
                    }`}
                  >
                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-4">
                        <div className="w-12 h-12 bg-gray-100 rounded-lg flex items-center justify-center text-2xl">
                          🔬
                        </div>
                        <div>
                          <h4 className="font-semibold text-gray-900">{machine.type}</h4>
                          <p className="text-sm text-gray-500">{machine.port}</p>
                        </div>
                      </div>

                      <div className="flex items-center gap-3">
                        <StatusBadge status={machine.status || 'disconnected'} />
                        {machine.status === 'connected' ? (
                          <Button
                            onClick={() => handleDisconnect(machine.id)}
                            variant="danger"
                            size="sm"
                          >
                            Disconnect
                          </Button>
                        ) : (
                          <Button
                            onClick={() => handleConnect(machine.id)}
                            variant="success"
                            size="sm"
                          >
                            Connect
                          </Button>
                        )}
                      </div>
                    </div>
                  </div>
                ))
              )}
            </div>
          </Card>
        </div>

        {/* Machine Details */}
        <div>
          <Card title="Machine Details">
            {selectedMachine ? (
              <div className="space-y-4">
                <div>
                  <label className="text-sm font-medium text-gray-600">Machine Type</label>
                  <p className="text-lg font-semibold mt-1">{selectedMachine.type}</p>
                </div>
                <div>
                  <label className="text-sm font-medium text-gray-600">Connection Type</label>
                  <p className="text-gray-900 mt-1">{selectedMachine.connection_type}</p>
                </div>
                <div>
                  <label className="text-sm font-medium text-gray-600">Port</label>
                  <p className="text-gray-900 mt-1">{selectedMachine.port}</p>
                </div>
                <div>
                  <label className="text-sm font-medium text-gray-600">Parse Format</label>
                  <p className="text-gray-900 mt-1">{selectedMachine.parse_format || 'N/A'}</p>
                </div>
                <div>
                  <label className="text-sm font-medium text-gray-600">Status</label>
                  <div className="mt-1">
                    <StatusBadge status={selectedMachine.status || 'disconnected'} />
                  </div>
                </div>

                <div className="pt-4 border-t">
                  <Button variant="secondary" className="w-full mb-2">
                    Edit Configuration
                  </Button>
                  <Button variant="danger" className="w-full">
                    Remove Machine
                  </Button>
                </div>
              </div>
            ) : (
              <div className="text-center py-8 text-gray-500">
                Select a machine to view details
              </div>
            )}
          </Card>
        </div>
      </div>
    </div>
  );
};

export default MachinesPage;
