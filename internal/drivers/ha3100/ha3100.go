package ha3100

import (
	"context"
	"fmt"
	"strings"
	"time"

	"veta-lis/internal/drivers"
)

// HA3100Driver implements the MachineDriver interface for HEAL HA-3100 analyzers.
// It handles serial communication and ASTM parsing.
type HA3100Driver struct {
	config      *drivers.MachineConfig
	isConnected bool
	serialPort  interface{} // Would be *serial.Port in production
	lastContact time.Time
}

// NewHA3100Driver creates a new HA3100 driver instance.
func NewHA3100Driver() drivers.MachineDriver {
	return &HA3100Driver{
		isConnected: false,
	}
}

// Connect establishes a serial connection to the HA3100 analyzer.
func (d *HA3100Driver) Connect(ctx context.Context) error {
	if d.isConnected {
		return drivers.ErrAlreadyConnected
	}

	// In production, this would:
	// 1. Open serial port (d.config.Port, typically COM3/COM4)
	// 2. Set baud rate (typically 9600 or 19200)
	// 3. Configure handshaking (usually RTS/CTS)
	// 4. Send initialization command (e.g., "*IDN?" to query device info)
	// 5. Wait for acknowledgment

	// Mock connection for demonstration
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		d.isConnected = true
		d.lastContact = time.Now()
		return nil
	}
}

// Disconnect closes the serial connection.
func (d *HA3100Driver) Disconnect(ctx context.Context) error {
	if !d.isConnected {
		return drivers.ErrNotConnected
	}

	// In production: Close serial port, clean up resources
	d.isConnected = false
	return nil
}

// IsConnected returns the connection status.
func (d *HA3100Driver) IsConnected() bool {
	return d.isConnected
}

// Listen continuously reads from the HA3100 and sends parsed results through the channel.
func (d *HA3100Driver) Listen(ctx context.Context, dataChan chan<- *drivers.LabResult) error {
	if !d.isConnected {
		return drivers.ErrNotConnected
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			// In production:
			// 1. Read from serial port (blocking or timeout)
			// 2. Accumulate data until complete message (frame boundaries)
			// 3. Call ParseData on complete message
			// 4. Send to channel
			// 5. Handle any parsing errors

			// Mock: simulate receiving data every 5 seconds
			time.Sleep(5 * time.Second)
			
			// This is where parsed results would be sent
			// dataChan <- result
		}
	}
}

// SendCommand sends a command to the HA3100 and waits for response.
func (d *HA3100Driver) SendCommand(ctx context.Context, cmd string) (string, error) {
	if !d.isConnected {
		return "", drivers.ErrNotConnected
	}

	// In production:
	// 1. Write command to serial port
	// 2. Set up timeout timer
	// 3. Read response from serial port
	// 4. Return response or timeout error

	return fmt.Sprintf("Response to: %s", cmd), nil
}

// GetMachineInfo returns information about the connected HA3100.
func (d *HA3100Driver) GetMachineInfo(ctx context.Context) (*drivers.MachineInfo, error) {
	if !d.isConnected {
		return nil, drivers.ErrNotConnected
	}

	return &drivers.MachineInfo{
		Type:            "HA3100",
		Model:           "HA-3100",
		SerialNumber:    "SN123456789", // Would be queried from device
		FirmwareVersion: "2.1.5",
		ConnectionType:  "SERIAL",
		Port:            d.config.Port,
		BaudRate:        d.config.BaudRate,
		Status:          "READY",
		LastTestTime:    d.lastContact.Format(time.RFC3339),
		TotalTestsRun:   0, // Would be queried from device
	}, nil
}

// ParseData parses ASTM-formatted data from the HA3100.
// ASTM format: |^~\&| frame start, multiple records separated by \n
func (d *HA3100Driver) ParseData(rawData []byte) (*drivers.LabResult, error) {
	if len(rawData) == 0 {
		return nil, drivers.ErrInvalidData
	}

	// In production, implement full ASTM parsing:
	// 1. Validate ASTM frame markers
	// 2. Split records (H, P, R, L)
	// 3. Extract test results from R (Result) records
	// 4. Parse field separators and handle escaping
	// 5. Map ASTM codes to lab test names
	// 6. Validate data integrity

	dataStr := string(rawData)
	
	// Mock parsing: extract test value from format like "GLUCOSE|85|mg/dL|60-110"
	parts := strings.Split(dataStr, "|")
	if len(parts) < 4 {
		return nil, drivers.ErrParseError
	}

	return &drivers.LabResult{
		ID:             fmt.Sprintf("HA3100-%d", time.Now().UnixNano()),
		MachineType:    "HA3100",
		TestName:       parts[0],
		ResultValue:    parts[1],
		ResultUnit:     parts[2],
		ReferenceRange: parts[3],
		Status:         "PENDING", // Will be determined by ResultValidator
		Timestamp:      time.Now().Format(time.RFC3339),
		RawData:        rawData,
	}, nil
}

// GetDriverType returns the driver type identifier.
func (d *HA3100Driver) GetDriverType() string {
	return "HA3100"
}

// SetConfig sets the machine configuration for this driver.
func (d *HA3100Driver) SetConfig(config *drivers.MachineConfig) {
	d.config = config
}
