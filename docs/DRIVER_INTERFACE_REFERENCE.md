# MachineDriver Interface - Complete Code Reference

## Core Interface Definition

```go
package drivers

import "context"

// MachineDriver is the core interface that all lab machines must implement.
// It follows the Adapter Pattern to support multiple machine types (HA3100, Scitek, Mindray).
// This enables seamless addition of new machines without modifying existing code.
type MachineDriver interface {
	// Connect establishes a connection to the lab machine (serial or TCP).
	// For HA3100: Opens COM port with specified baud rate
	// For Scitek: Establishes TCP socket to IP:port
	// For Mindray: Initiates TCP handshake with proprietary protocol
	// Returns an error if connection fails.
	Connect(ctx context.Context) error

	// Disconnect gracefully closes the connection to the lab machine.
	// Cleans up resources, stops listening goroutine.
	Disconnect(ctx context.Context) error

	// IsConnected returns true if the driver is currently connected to the machine.
	// Used for status checks and UI indicators.
	IsConnected() bool

	// Listen starts listening for incoming data from the lab machine.
	// This method should run in a goroutine and can be stopped by canceling the context.
	// Data flows through the provided channel.
	// Handles:
	// - HA3100: Reads serial frames, parses ASTM, sends LabResult
	// - Scitek: Reads TCP bytes, detects MLLP frames, parses HL7, sends LabResult
	// - Mindray: Reads TCP binary packets, parses CBC, sends LabResult
	Listen(ctx context.Context, dataChan chan<- *LabResult) error

	// SendCommand sends a command to the lab machine and returns the response.
	// HA3100: "*IDN?" queries device identifier
	// Scitek: Sends HL7 queries wrapped in MLLP frames
	// Mindray: Sends proprietary control commands
	SendCommand(ctx context.Context, cmd string) (string, error)

	// GetMachineInfo returns metadata about the connected machine.
	// Returns device type, model, firmware version, connection status, etc.
	GetMachineInfo(ctx context.Context) (*MachineInfo, error)

	// ParseData parses raw data from the machine into a LabResult struct.
	// Machine-specific parsing happens here:
	// - HA3100: Parses ASTM format (|^~\& delimiters, pipe-separated fields)
	// - Scitek: Parses HL7 format (\r segment delimiters, |-separated fields)
	// - Mindray: Parses binary hematology data (CBC parameters)
	ParseData(rawData []byte) (*LabResult, error)

	// GetDriverType returns the type of driver (e.g., "HA3100", "Scitek", "Mindray").
	// Used for logging and driver identification.
	GetDriverType() string
}

// MachineInfo contains metadata about a connected lab machine.
// Standardized across all drivers for consistent UI/logging.
type MachineInfo struct {
	Type              string // e.g., "HA3100", "Scitek", "Mindray"
	Model             string // e.g., "HA-3100", "Scitek Analyzer"
	SerialNumber      string
	FirmwareVersion   string
	ConnectionType    string // "SERIAL" or "TCP"
	Port              string // COM port (e.g., "COM3") or IP address with port (e.g., "192.168.1.100:5000")
	BaudRate          int    // For serial connections (9600, 19200, 115200, etc.)
	Status            string // "READY", "BUSY", "ERROR"
	LastTestTime      string // ISO 8601 timestamp
	TotalTestsRun     int    // Number of tests since connection
}

// LabResult represents a single test result from a lab machine.
// Standardized format regardless of source machine type.
// This is what flows through the dataChan from Listen().
type LabResult struct {
	ID                string      // Unique identifier (UUID or machine-generated)
	MachineType       string      // Type of machine that generated this result
	TestID            string      // Test ID from the machine (reference to the test run)
	PetID             string      // Reference to pet in database (filled by App layer)
	OwnerID           string      // Reference to pet owner (filled by App layer)
	TestName          string      // Name of the test performed (e.g., "Glucose", "WBC")
	TestCode          string      // Code of the test (e.g., "GLU", "WBC")
	ResultValue       string      // The actual result value (e.g., "85", "7.2")
	ResultUnit        string      // Unit of measurement (e.g., "mg/dL", "cells/µL")
	ReferenceRange    string      // Normal reference range (e.g., "60-110")
	Status            string      // "NORMAL", "ABNORMAL", "CRITICAL", "PENDING"
	Timestamp         string      // ISO 8601 timestamp when result was received
	RawData           []byte      // Raw data from machine (for debugging/re-parsing)
	Notes             string      // Additional notes or flags
}

// MachineConfig holds configuration for a specific machine connection.
// Loaded from config.yaml and validated by ConfigManager.
type MachineConfig struct {
	MachineType    string // Machine type identifier (must match registered driver)
	ConnectionType string // "SERIAL" or "TCP"
	Port           string // COM port (e.g., "COM3") or IP:port (e.g., "192.168.1.100:5000")
	BaudRate       int    // For serial connections (ignored for TCP)
	Timeout        int    // Connection timeout in seconds
	RetryAttempts  int    // Number of retry attempts on connection failure
	ParseFormat    string // Protocol format: "HL7", "ASTM", "CUSTOM"
}

// DriverFactory creates driver instances based on machine type.
// Uses Registry pattern for runtime driver selection.
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
// Example: factory.Register("HA3100", ha3100.NewHA3100Driver)
func (df *DriverFactory) Register(machineType string, factory func() MachineDriver) {
	df.drivers[machineType] = factory
}

// Create creates a new driver instance of the specified type.
// Returns error if machine type is not registered.
func (df *DriverFactory) Create(machineType string) (MachineDriver, error) {
	if factory, ok := df.drivers[machineType]; ok {
		return factory(), nil
	}
	return nil, ErrUnsupportedMachineType
}
```

---

## How Each Machine Type Implements the Interface

### HA3100 (Serial, ASTM Protocol)

```go
package ha3100

import (
	"context"
	"veta-lis/internal/drivers"
)

type HA3100Driver struct {
	config      *drivers.MachineConfig
	isConnected bool
	serialPort  interface{} // *serial.Port in production
}

// Connect opens COM port (e.g., COM3 @ 9600 baud)
func (d *HA3100Driver) Connect(ctx context.Context) error {
	// Open serial.Port
	// Set baud rate (d.config.BaudRate)
	// Configure handshaking
	// Send init command
	d.isConnected = true
	return nil
}

// Listen reads serial frames, parses ASTM, sends results
func (d *HA3100Driver) Listen(ctx context.Context, dataChan chan<- *drivers.LabResult) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			// Read from serial port (blocking with timeout)
			rawData := readSerialFrame()
			
			// Parse ASTM format
			result, err := d.ParseData(rawData)
			if err != nil {
				continue // Skip malformed data
			}
			
			// Send result through channel (non-blocking with buffer)
			dataChan <- result
		}
	}
}

// ParseData converts ASTM format to LabResult
// ASTM format: |^~\&| header, fields separated by |
// Example: "GLUCOSE|85|mg/dL|60-110"
func (d *HA3100Driver) ParseData(rawData []byte) (*drivers.LabResult, error) {
	// Validate ASTM frame markers
	// Extract test code, value, unit, reference range
	// Return standardized LabResult
	return &drivers.LabResult{
		MachineType: "HA3100",
		TestCode:    "GLU",
		ResultValue: "85",
		ResultUnit:  "mg/dL",
		Status:      "PENDING", // Will be validated by ResultValidator
	}, nil
}

func (d *HA3100Driver) GetDriverType() string {
	return "HA3100"
}
```

### Scitek (TCP, HL7 + MLLP)

```go
package scitek

import (
	"context"
	"net"
	"veta-lis/internal/drivers"
	"veta-lis/internal/parser"
)

type ScitekDriver struct {
	config      *drivers.MachineConfig
	isConnected bool
	tcpConn     net.Conn
	hl7Parser   *parser.HL7Parser
}

// Connect establishes TCP connection to IP:port
func (d *ScitekDriver) Connect(ctx context.Context) error {
	// Establish TCP connection to d.config.Port (e.g., "192.168.1.100:5000")
	// Implement reconnection with exponential backoff
	// Send Scitek handshake
	d.isConnected = true
	return nil
}

// Listen reads TCP data, detects MLLP frames, parses HL7, sends results
func (d *ScitekDriver) Listen(ctx context.Context, dataChan chan<- *drivers.LabResult) error {
	reader := bufio.NewReader(d.tcpConn)
	
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			// Read until MLLP start marker (0x0B)
			startByte, _ := reader.ReadByte()
			if startByte != 0x0B {
				continue
			}
			
			// Read until MLLP end marker (0x1C 0x0D)
			hl7Message := readUntilMLLPEnd(reader)
			
			// Parse HL7
			result, err := d.ParseData(hl7Message)
			if err != nil {
				continue
			}
			
			dataChan <- result
		}
	}
}

// ParseData converts HL7 format to LabResult
// HL7 segments separated by \r: MSH (header), PID (patient), OBX (result)
// Example HL7 segment: "OBX|3|NM|GLU^Glucose||85|mg/dL|60-110|N"
func (d *ScitekDriver) ParseData(rawData []byte) (*drivers.LabResult, error) {
	hl7Msg, err := d.hl7Parser.Parse(rawData)
	if err != nil {
		return nil, err
	}
	
	return &drivers.LabResult{
		MachineType: "Scitek",
		TestCode:    hl7Msg.TestCode,
		TestName:    hl7Msg.TestName,
		ResultValue: hl7Msg.ResultValue,
		ResultUnit:  hl7Msg.ResultUnit,
		Status:      "PENDING",
	}, nil
}

func (d *ScitekDriver) GetDriverType() string {
	return "Scitek"
}
```

### Mindray (TCP, Proprietary Binary)

```go
package mindray

import (
	"context"
	"net"
	"veta-lis/internal/drivers"
)

type MindrayDriver struct {
	config      *drivers.MachineConfig
	isConnected bool
	tcpConn     net.Conn
}

// Connect establishes TCP connection to Mindray analyzer
func (d *MindrayDriver) Connect(ctx context.Context) error {
	// Establish TCP connection to d.config.Port
	// Send Mindray-specific handshake
	// Initialize communication parameters
	d.isConnected = true
	return nil
}

// Listen reads TCP binary packets, parses hematology data
func (d *MindrayDriver) Listen(ctx context.Context, dataChan chan<- *drivers.LabResult) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			// Read Mindray binary frame
			frame := readMindrayFrame(d.tcpConn)
			
			// Parse CBC (Complete Blood Count) data
			result, err := d.ParseData(frame)
			if err != nil {
				continue
			}
			
			dataChan <- result
		}
	}
}

// ParseData converts Mindray binary format to LabResult
// Extracts: WBC, RBC, HGB, HCT, MCV, MCH, MCHC, PLT with flags
func (d *MindrayDriver) ParseData(rawData []byte) (*drivers.LabResult, error) {
	// Decode binary parameter blocks
	// Extract test values and flags (L=low, H=high, N=normal)
	// Return standardized LabResult
	
	return &drivers.LabResult{
		MachineType: "Mindray",
		TestName:    "Complete Blood Count",
		// CBC parameters extracted from binary...
		Status:      "PENDING",
	}, nil
}

func (d *MindrayDriver) GetDriverType() string {
	return "Mindray"
}
```

---

## Usage in Application Layer

```go
package app

import (
	"context"
	"veta-lis/internal/drivers"
	"veta-lis/internal/drivers/ha3100"
	"veta-lis/internal/drivers/scitek"
)

type App struct {
	driverFactory *drivers.DriverFactory
	activeDrivers map[string]drivers.MachineDriver
	resultChannel chan *drivers.LabResult
}

func (a *App) Startup(ctx context.Context) error {
	// Register all available drivers
	factory := drivers.NewDriverFactory()
	factory.Register("HA3100", ha3100.NewHA3100Driver)
	factory.Register("Scitek", scitek.NewScitekDriver)
	// ... register others
	
	a.driverFactory = factory
	return nil
}

// ConnectMachine called from frontend via Wails binding
func (a *App) ConnectMachine(machineID string) (string, error) {
	// Get config from ConfigManager
	machineConfig, _ := a.configManager.GetMachineByID(machineID)
	
	// Create appropriate driver (no hard-coded dependencies!)
	driver, err := a.driverFactory.Create(machineConfig.Type)
	if err != nil {
		return "", err
	}
	
	// Connect
	if err := driver.Connect(a.ctx); err != nil {
		return "", err
	}
	
	// Store driver
	a.activeDrivers[machineID] = driver
	
	// Start listening (runs in separate goroutine)
	go driver.Listen(a.ctx, a.resultChannel)
	
	// Get status
	info, _ := driver.GetMachineInfo(a.ctx)
	return info.Status, nil
}
```

---

## Key Insights

### ✅ Adapter Pattern Benefits

1. **No Hard-Coded Dependencies**
   - App layer doesn't know if it's HA3100 or Scitek
   - Factory determines at runtime

2. **Easy to Extend**
   - Add new machine? Just create driver + register
   - No changes to App layer

3. **Testable**
   - Mock driver can be injected
   - Unit tests don't depend on real hardware

### ✅ Interface Design

1. **Unified Data Format**
   - All machines return standardized `LabResult`
   - App layer always receives same structure

2. **Channel-Based Communication**
   - Non-blocking (buffered channel)
   - Multiple goroutines can safely send data

3. **Context-Based Control**
   - Graceful cancellation
   - Respects application lifecycle

---

**This design enables Veta LIS to support unlimited lab machines without modifying a single line of existing code!** 🎯
