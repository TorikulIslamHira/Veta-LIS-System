package main

import (
	"fmt"
	"strings"
	"time"
)

const (
	// MLLP framing bytes
	MLLP_START_BYTE = 0x0B // Vertical Tab (Start of Block)
	MLLP_END_BYTE_1 = 0x1C // File Separator
	MLLP_END_BYTE_2 = 0x0D // Carriage Return

	// Serial port configuration
	PORT_NAME = "COM3"
	BAUD_RATE = 9600
)

func main() {
	fmt.Println("🔬 Scitek Biochemistry Analyzer Simulator (Loopback Test)")
	fmt.Println("=========================================================")
	fmt.Printf("Port: %s | Baud: %d\n", PORT_NAME, BAUD_RATE)
	fmt.Println()
	fmt.Println("⚠️  LOOPBACK MODE: Testing message construction only")
	fmt.Println("    (No actual serial transmission - just validates MLLP framing)")
	fmt.Println()

	// Construct sample HL7 message for Scitek analyzer
	hl7Message := buildHL7Message()

	fmt.Println("📤 HL7 Message Content:")
	fmt.Println("------------------------")
	fmt.Println(hl7Message)
	fmt.Println("------------------------")

	// Wrap message in MLLP framing
	mllpMessage := wrapInMLLP(hl7Message)

	fmt.Printf("\n📦 Message Statistics:\n")
	fmt.Printf("   - Raw HL7 size: %d bytes\n", len(hl7Message))
	fmt.Printf("   - MLLP wrapped: %d bytes\n", len(mllpMessage))
	fmt.Printf("   - MLLP overhead: 3 bytes (start + 2 end bytes)\n")

	// Display MLLP structure
	fmt.Println("\n🔍 MLLP Frame Analysis:")
	fmt.Printf("   - Start byte: 0x%02X (VT)\n", mllpMessage[0])
	fmt.Printf("   - Message: %d bytes\n", len(hl7Message))
	fmt.Printf("   - End byte 1: 0x%02X (FS)\n", mllpMessage[len(mllpMessage)-2])
	fmt.Printf("   - End byte 2: 0x%02X (CR)\n", mllpMessage[len(mllpMessage)-1])

	// Validate framing
	if validateMLLPFrame(mllpMessage) {
		fmt.Println("\n✅ MLLP framing is VALID")
	} else {
		fmt.Println("\n❌ MLLP framing is INVALID")
		return
	}

	// Show first and last few bytes
	fmt.Println("\n📊 Raw Bytes Preview:")
	fmt.Printf("   First 20 bytes: %v\n", mllpMessage[:min(20, len(mllpMessage))])
	fmt.Printf("   Last 20 bytes:  %v\n", mllpMessage[max(0, len(mllpMessage)-20):])

	// Simulate ACK response
	ackMessage := buildMLLPACK()
	fmt.Println("\n📥 Expected ACK Response:")
	fmt.Printf("   ACK bytes: %v\n", ackMessage)
	fmt.Printf("   ACK string: %q\n", string(ackMessage))

	if validateMLLPFrame(ackMessage) {
		fmt.Println("   ✅ ACK framing is VALID")
	}

	// Parse the HL7 message
	fmt.Println("\n📋 Parsed HL7 Segments:")
	parseHL7(hl7Message)

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("✅ LOOPBACK TEST PASSED")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("\n💡 Next Steps:")
	fmt.Println("   1. Install virtual COM port software (com0com)")
	fmt.Println("   2. Create COM3 ↔ COM4 pair")
	fmt.Println("   3. Run scitek_sim (sends to COM4)")
	fmt.Println("   4. Run your driver (listens on COM3)")
}

// buildHL7Message constructs a sample HL7 ORU^R01 message for lab results
func buildHL7Message() string {
	timestamp := time.Now().Format("20060102150405")
	messageID := fmt.Sprintf("%d", time.Now().Unix())

	segments := []string{
		// MSH - Message Header
		fmt.Sprintf("MSH|^~\\&|Scitek||||%s||ORU^R01|%s|P|2.3.1", timestamp, messageID),

		// PID - Patient Identification
		"PID|1||101||Test^Patient^^^||19900101|M|||123 Main St^^Cityville^^12345",

		// OBR - Observation Request
		fmt.Sprintf("OBR|1||%s|CHEM^Chemistry Panel|||%s", messageID, timestamp),

		// OBX - Observation/Result (Glucose)
		"OBX|1|NM|GLU^Glucose||105|mg/dL|70-110|N|||F",

		// OBX - Observation/Result (Urea)
		"OBX|2|NM|UREA^Blood Urea Nitrogen||40|mg/dL|7-20|H|||F",

		// OBX - Observation/Result (Creatinine)
		"OBX|3|NM|CREA^Creatinine||1.2|mg/dL|0.6-1.2|N|||F",
	}

	message := ""
	for _, segment := range segments {
		message += segment + "\r"
	}

	return message
}

// wrapInMLLP wraps the HL7 message with MLLP framing bytes
func wrapInMLLP(message string) []byte {
	messageBytes := []byte(message)
	mllpMessage := make([]byte, 0, len(messageBytes)+3)

	mllpMessage = append(mllpMessage, MLLP_START_BYTE)
	mllpMessage = append(mllpMessage, messageBytes...)
	mllpMessage = append(mllpMessage, MLLP_END_BYTE_1, MLLP_END_BYTE_2)

	return mllpMessage
}

// buildMLLPACK creates a proper MLLP-wrapped ACK message
func buildMLLPACK() []byte {
	ack := "ACK"
	return wrapInMLLP(ack)
}

// validateMLLPFrame checks if the message has valid MLLP framing
func validateMLLPFrame(data []byte) bool {
	if len(data) < 3 {
		return false
	}

	if data[0] != MLLP_START_BYTE {
		return false
	}

	if data[len(data)-2] != MLLP_END_BYTE_1 || data[len(data)-1] != MLLP_END_BYTE_2 {
		return false
	}

	return true
}

// parseHL7 displays the parsed segments of an HL7 message
func parseHL7(message string) {
	segments := splitSegments(message)
	for i, segment := range segments {
		if segment == "" {
			continue
		}
		fields := splitFields(segment)
		segmentType := fields[0]
		fmt.Printf("   %d. %s (%d fields)\n", i+1, segmentType, len(fields))

		// Show key fields based on segment type
		switch segmentType {
		case "MSH":
			if len(fields) > 9 {
				fmt.Printf("      - Message Type: %s\n", fields[9])
				fmt.Printf("      - Message ID: %s\n", fields[10])
			}
		case "PID":
			if len(fields) > 5 {
				fmt.Printf("      - Patient ID: %s\n", fields[3])
				fmt.Printf("      - Patient Name: %s\n", fields[5])
			}
		case "OBX":
			if len(fields) > 5 {
				fmt.Printf("      - Test: %s\n", fields[3])
				fmt.Printf("      - Value: %s %s\n", fields[5], fields[6])
			}
		}
	}
}

// splitSegments splits HL7 message by carriage return
func splitSegments(message string) []string {
	segments := []string{}
	current := ""
	for _, ch := range message {
		if ch == '\r' {
			if current != "" {
				segments = append(segments, current)
				current = ""
			}
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		segments = append(segments, current)
	}
	return segments
}

// splitFields splits an HL7 segment by field separator
func splitFields(segment string) []string {
	fields := []string{}
	current := ""
	for _, ch := range segment {
		if ch == '|' {
			fields = append(fields, current)
			current = ""
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		fields = append(fields, current)
	}
	return fields
}

// Helper functions for Go versions without min/max built-ins
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
