import { useContext, useCallback } from 'react';
import { MachineContext } from '../context/MachineContext';
import * as app from '../../wailsjs/go/app/App';

/**
 * Custom hook for calling Wails-bound Go methods
 * Simplifies frontend-to-backend communication
 */

export const useMachineControl = () => {
  const context = useContext(MachineContext);
  if (!context) {
    throw new Error('useMachineControl must be used within MachineProvider');
  }

  const connectMachine = useCallback(async (machineId: string) => {
    try {
      const status = await app.ConnectMachine(machineId);
      // Update machine status in context
      context.updateMachineStatus(machineId, { connected: true, status });
      return status;
    } catch (error) {
      console.error('Failed to connect machine:', error);
      throw error;
    }
  }, [context]);

  const disconnectMachine = useCallback(async (machineId: string) => {
    try {
      await app.DisconnectMachine(machineId);
      context.updateMachineStatus(machineId, { connected: false });
    } catch (error) {
      console.error('Failed to disconnect machine:', error);
      throw error;
    }
  }, [context]);

  const getMachineStatus = useCallback(async (machineId: string) => {
    try {
      const status = await app.GetMachineStatus(machineId);
      context.updateMachineStatus(machineId, status);
      return status;
    } catch (error) {
      console.error('Failed to get machine status:', error);
      throw error;
    }
  }, [context]);

  const loadAllMachines = useCallback(async () => {
    try {
      const machines = await app.GetAllMachines();
      if (context.setMachines && Array.isArray(machines)) {
        context.setMachines(machines as any);
      }
      return machines;
    } catch (error) {
      console.error('Failed to load machines:', error);
      throw error;
    }
  }, [context]);

  return {
    connectMachine,
    disconnectMachine,
    getMachineStatus,
    loadAllMachines,
  };
};

export const useMachineConfig = () => {
  const saveMachineConfig = useCallback(async (
    id: string,
    connectionType: string,
    port: string,
    baudRate: number
  ) => {
    try {
      await app.SaveMachineConfig(id, connectionType, port, baudRate);
    } catch (error) {
      console.error('Failed to save machine config:', error);
      throw error;
    }
  }, []);

  return {
    saveMachineConfig,
  };
};
