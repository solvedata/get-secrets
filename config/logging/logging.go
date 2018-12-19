package logging

import (
	"time"

	"github.com/evalphobia/logrus_sentry"
	"github.com/mexisme/get-secrets/version"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Configure set-ups the Logrus library -- debug mode, etc
// Currently set-up via Viper
func Configure() {
	// We do this before setting debug mode, to help-out the log aggregators:
	switch loggingFormat := viper.GetString("logging.format"); loggingFormat {
	case "":
		fallthrough

	case "text":
		// This is the default log formatter in logrus, anyway:
		log.SetFormatter(&log.TextFormatter{})

	case "json":
		log.SetFormatter(&log.JSONFormatter{})

	default:
		log.Panicf("Log format %#v not supported.", loggingFormat)
	}

	// TODO: Should this be a Debug message?
	log.Infof("## %#v release %v %v ##", version.Application(), version.Release(), version.BuildDate())

	if viper.GetBool("debug") {
		log.SetLevel(log.DebugLevel)
		log.Debug("Debug mode enabled")
	}

	if sentryDsn := viper.GetString("logging.sentry.dsn"); sentryDsn != "" {
		if err := setupSentry(sentryDsn); err != nil {
			log.Error(err)
		}
	}

	// viper.Debug()
}

func setupSentry(sentryDsn string) error {
	log.WithFields(log.Fields{"sentry.dsn": sentryDsn}).Debug("Configuring connection to Sentry.io")

	// TODO: Meta-tag for environment
	// Some meta-tags
	tags := map[string]string{
		// TODO: Pick a better name, as this maps to "SECRETS_APP":
		"app": viper.GetString("app"),
	}

	// Sentry will only log for messages of the following severity:
	levels := []log.Level{
		log.PanicLevel,
		log.FatalLevel,
		log.ErrorLevel,
	}

	// Hook Sentry into Logrus:
	hook, err := logrus_sentry.NewWithTagsSentryHook(sentryDsn, tags, levels)
	if err != nil {
		return err
	}

	// Set the Sentry "release" version:
	log.WithFields(log.Fields{"release": version.Release()}).Debug("Setting release version in Sentry")
	hook.SetRelease(version.Release())

	//hook.StacktraceConfiguration.Enable = true

	// It seems as if the default 100ms is too short:
	hook.Timeout = 1 * time.Second

	// Now, add it into the Logrus hook-chain
	log.AddHook(hook)

	log.Info("Sentry enabled")

	return nil
}
