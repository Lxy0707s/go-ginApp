package flags

import "flag"

// Flags is a simple flag interrupter to print value and load correct config file
func Flags(defaultConfigFile string) string {
	cfg := flag.String("c", defaultConfigFile, "configuration file")
	flag.Parse()
	return *cfg
}
