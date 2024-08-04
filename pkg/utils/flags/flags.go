package flags

import (
	"errors"
	"flag"
)

const (
	configPathFlag = "config-path"
)

type CMDFlags struct {
	ConfigPath string
}

func ParseFlags() (*CMDFlags, error) {
	configPath := flag.String(configPathFlag, "config/local.yaml", "Configuration file path")

	flag.Parse()

	if *configPath == "" {
		return nil, errors.New("Configuration file path was not found in application flags")
	}

	return &CMDFlags{ConfigPath: *configPath}, nil
}

func MustParseFlags() *CMDFlags {
	flags, err := ParseFlags()
	if err != nil {
		panic(err)
	}

	return flags
}
