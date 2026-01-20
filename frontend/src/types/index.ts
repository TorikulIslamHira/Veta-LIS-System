// Wails runtime interface
export interface WailsRuntime {
  EventsOn(eventName: string, callback: (data: any) => void): void;
  EventsOff(eventName: string): void;
  EventsEmit(eventName: string, data?: any): void;
}

// Frontend types for machines
export interface Machine {
  id: string;
  type: 'HA3100' | 'Scitek' | 'Mindray' | string;
  enabled: boolean;
  connectionType: 'SERIAL' | 'TCP';
  port: string;
  connected: boolean;
}

export interface MachineStatus {
  type: string;
  model: string;
  connected: boolean;
  status: 'READY' | 'BUSY' | 'ERROR' | string;
  connectionType: string;
  port: string;
  firmwareVersion?: string;
  lastTestTime?: string;
  totalTestsRun?: number;
}

export interface LabResult {
  id: string;
  machineType: string;
  testID: string;
  petID?: string;
  ownerID?: string;
  testName: string;
  testCode: string;
  resultValue: string;
  resultUnit: string;
  referenceRange: string;
  status: 'NORMAL' | 'ABNORMAL' | 'CRITICAL' | 'PENDING';
  timestamp: string;
  notes?: string;
}

export interface Pet {
  id: string;
  name: string;
  species: string;
  breed: string;
  age: number;
  weight: number;
  ownerID: string;
}

export interface Owner {
  id: string;
  name: string;
  email: string;
  phone: string;
  address: string;
}
