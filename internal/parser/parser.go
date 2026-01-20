package parser

import (
	"fmt"
	"strings"
)

// HL7Message represents a parsed HL7 message.
type HL7Message struct {
	MessageID   string
	PatientID   string
	TestID      string
	TestName    string
	TestCode    string
	ResultValue string
	ResultUnit  string
	ReferenceRange string
	Status      string // "F" (Final), "P" (Preliminary), etc.
	Raw         []byte
}

// HL7Parser handles parsing of HL7 v2.5+ messages.
type HL7Parser struct {
	fieldSeparator      byte
	componentSeparator  byte
	repetitionSeparator byte
	escapeCharacter     byte
	subcomponentSeparator byte
}

// NewHL7Parser creates a new HL7 parser instance.
func NewHL7Parser() *HL7Parser {
	return &HL7Parser{
		fieldSeparator:        '|',
		componentSeparator:    '^',
		repetitionSeparator:   '~',
		escapeCharacter:       '\\',
		subcomponentSeparator: '&',
	}
}

// Parse parses raw HL7 data into an HL7Message struct.
// Handles MLLP framing (0x0B at start, 0x1C 0x0D at end) if present.
func (p *HL7Parser) Parse(rawData []byte) (*HL7Message, error) {
	if len(rawData) == 0 {
		return nil, fmt.Errorf("empty data")
	}

	data := rawData

	// Remove MLLP framing if present
	if len(data) > 2 && data[0] == 0x0B {
		// Remove start byte (0x0B)
		data = data[1:]
		// Remove end bytes (0x1C 0x0D)
		if len(data) >= 2 && data[len(data)-2] == 0x1C && data[len(data)-1] == 0x0D {
			data = data[:len(data)-2]
		}
	}

	dataStr := string(data)
	segments := strings.Split(dataStr, "\r")

	msg := &HL7Message{
		Raw: rawData,
	}

	// Parse each segment
	for _, segment := range segments {
		if len(segment) == 0 {
			continue
		}

		fields := strings.Split(segment, string(p.fieldSeparator))
		if len(fields) == 0 {
			continue
		}

		segmentID := fields[0]

		switch segmentID {
		case "MSH": // Message header
			if len(fields) > 3 {
				msg.MessageID = fields[3]
			}

		case "PID": // Patient identification
			if len(fields) > 3 {
				msg.PatientID = fields[3]
			}

		case "OBX": // Observation result
			if len(fields) > 3 {
				// OBX[3] = Observation identifier (code)
				msg.TestCode = fields[3]
			}
			if len(fields) > 4 {
				// OBX[4] = Value type
				msg.TestName = fields[4]
			}
			if len(fields) > 5 {
				// OBX[5] = Observation value
				msg.ResultValue = fields[5]
			}
			if len(fields) > 6 {
				// OBX[6] = Units of measure
				msg.ResultUnit = fields[6]
			}
			if len(fields) > 7 {
				// OBX[7] = Reference range
				msg.ReferenceRange = fields[7]
			}
			if len(fields) > 11 {
				// OBX[11] = Result status
				msg.Status = fields[11]
			}

		case "ORC": // Order record
			if len(fields) > 2 {
				msg.TestID = fields[2]
			}
		}
	}

	return msg, nil
}

// ASTMParser handles parsing of ASTM E1381-04 standard messages.
type ASTMParser struct {
	recordSeparator  byte
	componentSeparator byte
	repetitionSeparator byte
	escapeCharacter  byte
}

// NewASTMParser creates a new ASTM parser instance.
func NewASTMParser() *ASTMParser {
	return &ASTMParser{
		recordSeparator:  '\n',
		componentSeparator: '|',
		repetitionSeparator: '\\',
		escapeCharacter:  '&',
	}
}

// ASTMMessage represents a parsed ASTM message.
type ASTMMessage struct {
	MessageID    string
	PatientID    string
	SpecimenID   string
	TestCode     string
	TestName     string
	ResultValue  string
	ResultUnit   string
	ReferenceMin string
	ReferenceMax string
	Status       string
	Raw          []byte
}

// Parse parses raw ASTM data into an ASTMMessage struct.
func (p *ASTMParser) Parse(rawData []byte) (*ASTMMessage, error) {
	if len(rawData) == 0 {
		return nil, fmt.Errorf("empty data")
	}

	dataStr := string(rawData)
	records := strings.Split(dataStr, "\n")

	msg := &ASTMMessage{
		Raw: rawData,
	}

	for _, record := range records {
		if len(record) == 0 {
			continue
		}

		fields := strings.Split(record, string(p.componentSeparator))
		if len(fields) == 0 {
			continue
		}

		recordType := fields[0]

		switch recordType {
		case "H": // Header record
			// H|\\^&|LaboratoryName|...
			if len(fields) > 2 {
				msg.MessageID = fields[2]
			}

		case "P": // Patient record
			// P|1|PatientID|...
			if len(fields) > 2 {
				msg.PatientID = fields[2]
			}

		case "R": // Result record
			// R|1|TestCode|TestValue|Unit|ReferenceRange|Status|...
			if len(fields) > 2 {
				msg.TestCode = fields[2]
			}
			if len(fields) > 3 {
				msg.ResultValue = fields[3]
			}
			if len(fields) > 4 {
				msg.ResultUnit = fields[4]
			}
			if len(fields) > 5 {
				// Parse reference range (e.g., "60-110")
				rangeStr := fields[5]
				parts := strings.Split(rangeStr, "-")
				if len(parts) >= 1 {
					msg.ReferenceMin = parts[0]
				}
				if len(parts) >= 2 {
					msg.ReferenceMax = parts[1]
				}
			}
			if len(fields) > 6 {
				msg.Status = fields[6]
			}

		case "L": // Terminator record
			// End of message, no action needed
		}
	}

	return msg, nil
}
