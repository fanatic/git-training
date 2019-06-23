package configor

import (
	"os"
	"regexp"
)

type Configor struct {
	*Config
}

type Config struct {
	Environment string
	ENVPrefix   string
}

// New initialize a Configor
func New(config *Config) *Configor {
	if config == nil {
		config = &Config{}
	}
	return &Configor{Config: config}
}

// GetEnvironment get environment
func (configor *Configor) GetEnvironment() string {
	if configor.Environment == "" {
		if env := os.Getenv("ENVIRONMENT"); env != "" {
			return env
		}

		if isTest, _ := regexp.MatchString("/_test/", os.Args[0]); isTest {
			return "test"
		}

		return "development"
	}
	return configor.Environment
}

// LoadFromSingleFile loads a single file with map[string]Config (string is environment)
func (configor *Configor) LoadFromSingleFile(config interface{}, file string) error {
	if err := processFileWithEnvironment(config, file, configor.GetEnvironment()); err != nil {
		return err
	}

	prefix := configor.getENVPrefix(config)
	if prefix == "-" {
		return processTags(config)
	}
	return processTags(config, prefix)
}

// Load will unmarshal configurations to struct from files that you provide
func (configor *Configor) Load(config interface{}, files ...string) error {
	for _, file := range configor.getConfigurationFiles(files...) {
		if err := processFile(config, file); err != nil {
			return err
		}
	}

	if prefix := configor.getENVPrefix(config); prefix == "-" {
		return processTags(config)
	} else {
		return processTags(config, prefix)
	}
}

// ENV return environment
func ENV() string {
	return New(nil).GetEnvironment()
}

// Load will unmarshal configurations to struct from files that you provide
func Load(config interface{}, files ...string) error {
	return New(nil).Load(config, files...)
}

// LoadFromSingleFile loads a single file with map[string]Config (string is environment)
func LoadFromSingleFile(config interface{}, file string) error {
	return New(nil).LoadFromSingleFile(config, file)
}
