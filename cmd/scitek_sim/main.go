package main

import (
	"fmt"
	"log"
	"time"

	"go.bug.st/serial"
)

const (
	// MLLP framing bytes
	MLLP_START_BYTE = 0x0B // Vertical Tab (Start of Block)
	MLLP_END_BYTE_1 = 0x1C // File Separator
	MLLP_END_BYTE_2 = 0x0D // Carriage Return

	// Serial port configuration
	PORT_NAME = "COM4"
	BAUD_RATE = 9600
)

func main() {
	fmt.Println("🔬 Scitek Biochemistry Analyzer Simulator")
	fmt.Println("==========================================")
	fmt.Printf("Port: %s | Baud: %d\n\n", PORT_NAME, BAUD_RATE)

	// Configure serial port
	mode := &serial.Mode{
		BaudRate: BAUD_RATE,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}

	// Open serial port
	port, err := serial.Open(PORT_NAME, mode)
	if err != nil {
		log.Fatalf("❌ Failed to open port %s: %v", PORT_NAME, err)
	}
	defer port.Close()

	fmt.Printf("✅ Serial port %s opened successfully\n", PORT_NAME)

	// Set read timeout
	err = port.SetReadTimeout(5 * time.Second)
	if err != nil {
		log.Fatalf("❌ Failed to set read timeout: %v", err)
	}

	// Construct sample HL7 message for Scitek analyzer
	hl7Message := buildHL7Message()

	fmt.Println("\n📤 Sending HL7 Message:")
	fmt.Println("------------------------")
	fmt.Println(hl7Message)
	fmt.Println("------------------------")

	// Wrap message in MLLP framing
	mllpMessage := wrapInMLLP(hl7Message)

	fmt.Printf("\n📦 Total bytes to send: %d (including MLLP framing)\n", len(mllpMessage))

	// Send the MLLP-wrapped message
	n, err := port.Write(mllpMessage)
	if err != nil {
		log.Fatalf("❌ Failed to write to port: %v", err)
	}

	fmt.Printf("✅ Sent %d bytes to %s\n", n, PORT_NAME)

	// Wait for ACK response
	fmt.Println("\n⏳ Waiting for ACK from driver...")

	ackBuffer := make([]byte, 100)
	n, err = port.Read(ackBuffer)
	if err != nil {
		log.Printf("❌ Timeout or error reading ACK: %v", err)
		return
	}

	if n > 0 {
		fmt.Printf("✅ ACK Received (%d bytes)\n", n)
		fmt.Printf("📥 ACK Data: %v\n", ackBuffer[:n])
		fmt.Printf("📥 ACK String: %q\n", string(ackBuffer[:n]))

		// Validate ACK format
		if isValidACK(ackBuffer[:n]) {
			fmt.Println("✅ Valid MLLP ACK received - transmission successful!")
		} else {
			fmt.Println("⚠️  Received data but not a valid MLLP ACK")
		}
	} else {
		fmt.Println("❌ No ACK received - timeout")
	}

	fmt.Println("\n🏁 Simulator finished")
}

// buildHL7Message constructs a sample HL7 ORU^R01 message for lab results
func buildHL7Message() string {
	// Using \r as segment separator (carriage return)
	// Real HL7 uses CR (0x0D) as segment terminator
	timestamp := time.Now().Format("20060102150405")
	messageID := fmt.Sprintf("%d", time.Now().Unix())

	segments := []string{
		// MSH - Message Header
		fmt.Sprintf("MSH|^~\\&|Scitek||||%s||ORU^R01|%s|P|2.3.1", timestamp, messageID),

		// PID - Patient Identification
		"PID|1||101||Test^Patient^^^||19900101|M|||123 Main St^^Cityville^^12345",

		// OBR - Observation Request (optional but recommended)
		fmt.Sprintf("OBR|1||%s|CHEM^Chemistry Panel|||%s", messageID, timestamp),

		// OBX - Observation/Result (Glucose)
		"OBX|1|NM|GLU^Glucose||105|mg/dL|70-110|N|||F",

		// OBX - Observation/Result (Urea)
		"OBX|2|NM|UREA^Blood Urea Nitrogen||40|mg/dL|7-20|H|||F",

		// OBX - Observation/Result (Creatinine) - Extra test
		"OBX|3|NM|CREA^Creatinine||1.2|mg/dL|0.6-1.2|N|||F",
	}

	// Join segments with carriage return
	message := ""
	for _, segment := range segments {
		message += segment + "\r"
	}

	return message
}

// wrapInMLLP wraps the HL7 message with MLLP framing bytes
func wrapInMLLP(message string) []byte {
	// MLLP Format: <VT>message<FS><CR>
	// VT = 0x0B (start)
	// FS = 0x1C (end byte 1)
	// CR = 0x0D (end byte 2)

	messageBytes := []byte(message)
	mllpMessage := make([]byte, 0, len(messageBytes)+3)

	// Add start byte
	mllpMessage = append(mllpMessage, MLLP_START_BYTE)

	// Add message content
	mllpMessage = append(mllpMessage, messageBytes...)

	// Add end bytes
	mllpMessage = append(mllpMessage, MLLP_END_BYTE_1, MLLP_END_BYTE_2)

	return mllpMessage
}

// isValidACK checks if the received data is a valid MLLP ACK
func isValidACK(data []byte) bool {
	// Expected ACK format: <VT>ACK<FS><CR> or similar
	// Minimum length: 3 bytes (start + something + end)
	if len(data) < 3 {
		return false
	}

	// Check MLLP framing
	if data[0] != MLLP_START_BYTE {
		return false
	}

	// Check for end bytes (should be at the end)
	if len(data) >= 2 {
		if data[len(data)-2] == MLLP_END_BYTE_1 && data[len(data)-1] == MLLP_END_BYTE_2 {
			return true
		}
	}

	return false
}
