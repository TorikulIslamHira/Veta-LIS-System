import React, { createContext, useState, useCallback, ReactNode } from 'react';

export interface Machine {
  id: string;
  type: string;
  connectionType: string;
  port: string;
  enabled: boolean;
  connected: boolean;
}

export interface MachineStatus {
  [machineId: string]: {
    type: string;
    model: string;
    connected: boolean;
    status: string;
    connectionType: string;
    port: string;
  };
}

interface MachineContextType {
  machines: Machine[];
  setMachines: (machines: Machine[]) => void;
  activeMachine: string | null;
  setActiveMachine: (id: string | null) => void;
  machineStatus: MachineStatus;
  updateMachineStatus: (id: string, status: any) => void;
}

export const MachineContext = createContext<MachineContextType | undefined>(undefined);

interface MachineProviderProps {
  children: ReactNode;
}

export const MachineProvider: React.FC<MachineProviderProps> = ({ children }) => {
  const [machines, setMachines] = useState<Machine[]>([]);
  const [activeMachine, setActiveMachine] = useState<string | null>(null);
  const [machineStatus, setMachineStatus] = useState<MachineStatus>({});

  const updateMachineStatus = useCallback((id: string, status: any) => {
    setMachineStatus(prev => ({
      ...prev,
      [id]: status,
    }));
  }, []);

  return (
    <MachineContext.Provider value={{
      machines,
      setMachines,
      activeMachine,
      setActiveMachine,
      machineStatus,
      updateMachineStatus,
    }}>
      {children}
    </MachineContext.Provider>
  );
};
