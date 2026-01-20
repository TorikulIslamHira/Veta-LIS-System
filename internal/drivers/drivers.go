package drivers

import "context"

// MachineDriver is the core interface that all lab machines must implement.
// It follows the Adapter Pattern to support multiple machine types (HA3100, Scitek, Mindray).
type MachineDriver interface {
	// Connect establishes a connection to the lab machine (serial or TCP).
	// Returns an error if connection fails.
	Connect(ctx context.Context) error

	// Disconnect gracefully closes the connection to the lab machine.
	Disconnect(ctx context.Context) error

	// IsConnected returns true if the driver is currently connected to the machine.
	IsConnected() bool

	// Listen starts listening for incoming data from the lab machine.
	// Data is sent through the provided channel.
	// This method should run in a goroutine and can be stopped by canceling the context.
	Listen(ctx context.Context, dataChan chan<- *LabResult) error

	// SendCommand sends a command to the lab machine and returns the response.
	SendCommand(ctx context.Context, cmd string) (string, error)

	// GetMachineInfo returns metadata about the connected machine.
	GetMachineInfo(ctx context.Context) (*MachineInfo, error)

	// ParseData parses raw data from the machine into a LabResult struct.
	// This is called internally by Listen() after receiving data.
	ParseData(rawData []byte) (*LabResult, error)

	// GetDriverType returns the type of driver (e.g., "HA3100", "Scitek", "Mindray").
	GetDriverType() string
}

// MachineInfo contains metadata about a connected lab machine.
type MachineInfo struct {
	Type              string // e.g., "HA3100", "Scitek", "Mindray"
	Model             string // e.g., "HA-3100"
	SerialNumber      string
	FirmwareVersion   string
	ConnectionType    string // "SERIAL" or "TCP"
	Port              string // COM port (e.g., "COM3") or IP address with port (e.g., "192.168.1.100:5000")
	BaudRate          int    // For serial connections
	Status            string // "READY", "BUSY", "ERROR"
	LastTestTime      string
	TotalTestsRun     int
}

// LabResult represents a single test result from a lab machine.
type LabResult struct {
	ID                string      // Unique identifier (UUID)
	MachineType       string      // Type of machine that generated this result
	TestID            string      // Test ID from the machine
	PetID             string      // Reference to pet (from DB or input)
	OwnerID           string      // Reference to pet owner
	TestName          string      // Name of the test performed
	TestCode          string      // Code of the test
	ResultValue       string      // The actual result value
	ResultUnit        string      // Unit of measurement (e.g., "mg/dL")
	ReferenceRange    string      // Normal reference range
	Status            string      // "NORMAL", "ABNORMAL", "CRITICAL", "PENDING"
	Timestamp         string      // ISO 8601 timestamp
	RawData           []byte      // Raw data from machine (for debugging)
	Notes             string      // Additional notes
}

// DriverFactory creates driver instances based on machine type.
type DriverFactory struct {
	drivers map[string]func() MachineDriver
}

// NewDriverFactory creates a new driver factory.
func NewDriverFactory() *DriverFactory {
	return &DriverFactory{
		drivers: make(map[string]func() MachineDriver),
	}
}

// Register registers a new driver type with the factory.
func (df *DriverFactory) Register(machineType string, factory func() MachineDriver) {
	df.drivers[machineType] = factory
}

// Create creates a new driver instance of the specified type.
func (df *DriverFactory) Create(machineType string) (MachineDriver, error) {
	if factory, ok := df.drivers[machineType]; ok {
		return factory(), nil
	}
	return nil, ErrUnsupportedMachineType
}

// MachineConfig holds configuration for a specific machine connection.
type MachineConfig struct {
	MachineType   string
	ConnectionType string // "SERIAL" or "TCP"
	Port          string   // COM port or IP:port
	BaudRate      int      // For serial (e.g., 9600, 115200)
	Timeout       int      // Connection timeout in seconds
	RetryAttempts int      // Number of retry attempts
	ParseFormat   string   // e.g., "HL7", "ASTM", "CUSTOM"
}
