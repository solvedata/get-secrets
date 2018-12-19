package config

import (
	"github.com/mexisme/get-secrets/config/logging"
	"github.com/mexisme/get-secrets/config/settings"
)

//nolint:gochecknoglobals
var (
	// We don't want to try to reinitialise the config more than once
	initConfigDone = false
	// We don't want to try to reinitialise the logging more than once
	logConfigDone = false
)

//nolint:gochecknoinits
func init() {
	readConfig()
	configLogging()
}

// ImportMe is to allow other packages to easily depend on this one,
// since most of the important logic is in init()
func ImportMe() {
}

// AddConfigItems passes the configItems through to config.AddConfigItems()
func AddConfigItems(configItems []string) {
	// Need to ensure the system has been configured at least once!
	readConfig() // TODO: Viper dynamically reads -- this may not be needed.
	settings.AddConfigItems(configItems)
}

func readConfig() {
	// This should make it safe to rerun a few times
	if !initConfigDone {
		settings.ReadConfig()
		initConfigDone = true
	}
}

func configLogging() {
	readConfig()

	// This should make it safe to rerun a few times
	if !logConfigDone {
		logging.Configure()
		logConfigDone = true
	}
}
