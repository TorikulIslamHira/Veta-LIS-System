package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
	"veta-lis/internal/drivers"
)

// SystemConfig holds the entire system configuration.
type SystemConfig struct {
	Database DatabaseConfig `yaml:"database"`
	Machines []MachineConfig `yaml:"machines"`
	Server   ServerConfig   `yaml:"server"`
}

// DatabaseConfig holds PostgreSQL connection settings.
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	SSLMode  string `yaml:"ssl_mode"` // "disable", "require", etc.
}

// MachineConfig holds configuration for a specific lab machine.
type MachineConfig struct {
	ID             string `yaml:"id"`              // Unique identifier (e.g., "machine_1")
	Type           string `yaml:"type"`            // "HA3100", "Scitek", "Mindray"
	Enabled        bool   `yaml:"enabled"`
	ConnectionType string `yaml:"connection_type"` // "SERIAL" or "TCP"
	Port           string `yaml:"port"`            // COM port or IP:port
	BaudRate       int    `yaml:"baud_rate"`       // For serial (0 for TCP)
	Timeout        int    `yaml:"timeout"`         // Seconds
	RetryAttempts  int    `yaml:"retry_attempts"`
	ParseFormat    string `yaml:"parse_format"` // "HL7", "ASTM", "CUSTOM"
}

// ServerConfig holds Wails server settings.
type ServerConfig struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

// ConfigManager handles loading and managing system configuration.
type ConfigManager struct {
	config *SystemConfig
	path   string
}

// NewConfigManager creates a new configuration manager.
func NewConfigManager(configPath string) *ConfigManager {
	return &ConfigManager{
		path: configPath,
	}
}

// Load loads configuration from YAML file.
func (cm *ConfigManager) Load() error {
	data, err := os.ReadFile(cm.path)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	cm.config = &SystemConfig{}
	if err := yaml.Unmarshal(data, cm.config); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	return cm.Validate()
}

// Save saves current configuration to YAML file.
func (cm *ConfigManager) Save() error {
	data, err := yaml.Marshal(cm.config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(cm.path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Validate validates the loaded configuration.
func (cm *ConfigManager) Validate() error {
	if cm.config == nil {
		return fmt.Errorf("configuration not loaded")
	}

	if cm.config.Database.Host == "" {
		return fmt.Errorf("database host not configured")
	}

	if cm.config.Database.User == "" {
		return fmt.Errorf("database user not configured")
	}

	if cm.config.Database.Database == "" {
		return fmt.Errorf("database name not configured")
	}

	for _, machine := range cm.config.Machines {
		if machine.Type == "" {
			return fmt.Errorf("machine type not specified: %s", machine.ID)
		}
		if machine.Port == "" {
			return fmt.Errorf("machine port not specified: %s", machine.ID)
		}
	}

	return nil
}

// GetDatabase returns the database configuration.
func (cm *ConfigManager) GetDatabase() DatabaseConfig {
	return cm.config.Database
}

// GetMachines returns all enabled machine configurations.
func (cm *ConfigManager) GetMachines() []MachineConfig {
	var enabled []MachineConfig
	for _, m := range cm.config.Machines {
		if m.Enabled {
			enabled = append(enabled, m)
		}
	}
	return enabled
}

// GetMachineByID returns a specific machine configuration by ID.
func (cm *ConfigManager) GetMachineByID(id string) (*MachineConfig, error) {
	for _, m := range cm.config.Machines {
		if m.ID == id {
			return &m, nil
		}
	}
	return nil, fmt.Errorf("machine not found: %s", id)
}

// AddMachine adds a new machine configuration.
func (cm *ConfigManager) AddMachine(machine MachineConfig) {
	cm.config.Machines = append(cm.config.Machines, machine)
}

// UpdateMachine updates an existing machine configuration.
func (cm *ConfigManager) UpdateMachine(id string, machine MachineConfig) error {
	for i, m := range cm.config.Machines {
		if m.ID == id {
			cm.config.Machines[i] = machine
			return nil
		}
	}
	return fmt.Errorf("machine not found: %s", id)
}

// RemoveMachine removes a machine configuration by ID.
func (cm *ConfigManager) RemoveMachine(id string) error {
	for i, m := range cm.config.Machines {
		if m.ID == id {
			cm.config.Machines = append(cm.config.Machines[:i], cm.config.Machines[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("machine not found: %s", id)
}

// ToDriverConfig converts MachineConfig to driver.MachineConfig.
func (mc *MachineConfig) ToDriverConfig() *drivers.MachineConfig {
	return &drivers.MachineConfig{
		MachineType:   mc.Type,
		ConnectionType: mc.ConnectionType,
		Port:          mc.Port,
		BaudRate:      mc.BaudRate,
		Timeout:       mc.Timeout,
		RetryAttempts: mc.RetryAttempts,
		ParseFormat:   mc.ParseFormat,
	}
}

// GetServer returns the server configuration.
func (cm *ConfigManager) GetServer() ServerConfig {
	if cm.config == nil {
		return ServerConfig{Port: 3000, Host: "localhost"}
	}
	return cm.config.Server
}
