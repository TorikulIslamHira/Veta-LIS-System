package scitek

import (
	"context"
	"fmt"
	"time"

	"veta-lis/internal/drivers"
	"veta-lis/internal/parser"
)

// ScitekDriver implements the MachineDriver interface for Scitek analyzers.
// It handles TCP/IP communication and HL7 parsing with MLLP framing.
type ScitekDriver struct {
	config      *drivers.MachineConfig
	isConnected bool
	tcpConn     interface{} // Would be *net.Conn in production
	lastContact time.Time
	hl7Parser   *parser.HL7Parser
}

// NewScitekDriver creates a new Scitek driver instance.
func NewScitekDriver() drivers.MachineDriver {
	return &ScitekDriver{
		isConnected: false,
		hl7Parser:   parser.NewHL7Parser(),
	}
}

// Connect establishes a TCP connection to the Scitek analyzer.
func (d *ScitekDriver) Connect(ctx context.Context) error {
	if d.isConnected {
		return drivers.ErrAlreadyConnected
	}

	// In production:
	// 1. Establish TCP connection to d.config.Port (typically IP:5000 or similar)
	// 2. Implement reconnection with exponential backoff
	// 3. Send initial handshake (Scitek-specific protocol)
	// 4. Start keep-alive mechanism

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		d.isConnected = true
		d.lastContact = time.Now()
		return nil
	}
}

// Disconnect closes the TCP connection.
func (d *ScitekDriver) Disconnect(ctx context.Context) error {
	if !d.isConnected {
		return drivers.ErrNotConnected
	}

	// In production: Close TCP connection, clean up resources
	d.isConnected = false
	return nil
}

// IsConnected returns the connection status.
func (d *ScitekDriver) IsConnected() bool {
	return d.isConnected
}

// Listen continuously reads HL7 messages from Scitek with MLLP framing.
// MLLP format: 0x0B (start) + HL7 message + 0x1C + 0x0D (end)
func (d *ScitekDriver) Listen(ctx context.Context, dataChan chan<- *drivers.LabResult) error {
	if !d.isConnected {
		return drivers.ErrNotConnected
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			// In production:
			// 1. Read from TCP connection (buffered reader with timeout)
			// 2. Look for MLLP frame boundaries (0x0B at start, 0x1C 0x0D at end)
			// 3. Extract HL7 message
			// 4. Call ParseData
			// 5. Send results to channel
			// 6. Handle partial messages and reconnection

			time.Sleep(5 * time.Second)
			
			// Mock: dataChan <- result
		}
	}
}

// SendCommand sends a command/query to the Scitek analyzer.
func (d *ScitekDriver) SendCommand(ctx context.Context, cmd string) (string, error) {
	if !d.isConnected {
		return "", drivers.ErrNotConnected
	}

	// In production:
	// 1. Wrap command in MLLP frame
	// 2. Send over TCP
	// 3. Wait for acknowledgment

	return fmt.Sprintf("Scitek ACK: %s", cmd), nil
}

// GetMachineInfo returns information about the connected Scitek analyzer.
func (d *ScitekDriver) GetMachineInfo(ctx context.Context) (*drivers.MachineInfo, error) {
	if !d.isConnected {
		return nil, drivers.ErrNotConnected
	}

	return &drivers.MachineInfo{
		Type:            "Scitek",
		Model:           "Scitek Analyzer",
		SerialNumber:    "SK987654321",
		FirmwareVersion: "3.0.2",
		ConnectionType:  "TCP",
		Port:            d.config.Port,
		Status:          "READY",
		LastTestTime:    d.lastContact.Format(time.RFC3339),
		TotalTestsRun:   0,
	}, nil
}

// ParseData parses HL7-formatted data from the Scitek analyzer.
// HL7 format with segments separated by \r
func (d *ScitekDriver) ParseData(rawData []byte) (*drivers.LabResult, error) {
	if len(rawData) == 0 {
		return nil, drivers.ErrInvalidData
	}

	// Use the HL7 parser from the parser package
	hl7Message, err := d.hl7Parser.Parse(rawData)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", drivers.ErrParseError, err)
	}

	// Convert HL7 parsed data to LabResult
	result := &drivers.LabResult{
		ID:          fmt.Sprintf("Scitek-%d", time.Now().UnixNano()),
		MachineType: "Scitek",
		TestID:      hl7Message.TestID,
		TestName:    hl7Message.TestName,
		TestCode:    hl7Message.TestCode,
		ResultValue: hl7Message.ResultValue,
		ResultUnit:  hl7Message.ResultUnit,
		Status:      "PENDING",
		Timestamp:   time.Now().Format(time.RFC3339),
		RawData:     rawData,
	}

	return result, nil
}

// GetDriverType returns the driver type identifier.
func (d *ScitekDriver) GetDriverType() string {
	return "Scitek"
}

// SetConfig sets the machine configuration for this driver.
func (d *ScitekDriver) SetConfig(config *drivers.MachineConfig) {
	d.config = config
}
