/*
Package main is the main loop.

Note the one addiional config field:

- s3.dotenv_path ($SECRETS_S3_DOTENV_PATH) -- defines the dir where to read all .env files from
*/
package main

import (
	"fmt"
	"os"

	"github.com/mexisme/get-secrets/config"
	"github.com/mexisme/get-secrets/dotenv"
	"github.com/mexisme/get-secrets/errors"
	execish "github.com/mexisme/get-secrets/exec"
	"github.com/mexisme/multiconfig"
	"github.com/mexisme/multiconfig/env"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

//nolint:gochecknoinits
func init() {
	config.ImportMe()
}

func main() {
	// When any other part of the app panics, we'd prefer to give them a "friendlier" face
	defer errors.Recovery()

	appName, appEnv := viper.GetString("application.name"), viper.GetString("application.environment")
	log.Infof("Preparing to run app %#v in env %#v", appName, appEnv)

	secrets := multiconfig.New()

	if viper.GetBool("dotenv.skip") {
		log.Info("Not getting .env secrets due to configuration")

	} else {
		dotenv.ReadFromS3(secrets)
	}

	// Add the SSM / KMS stuff in here?

	secrets.AddItem(env.New().FromOsEnviron())
	envs := dotenv.EnvMerge(secrets)

	if len(os.Args) > 1 {
		runner := execish.New().WithEnvs(envs)
		errors.PanicOnErrors(runner.WithCommand(os.Args[1:]).Exec())
	}

	fmt.Println("")
	fmt.Println("# No command provided to execute")
	osEnviron, err := envs.ToOsEnviron()
	if err != nil {
		errors.PanicOnErrors(err)
	}
	for _, envLine := range osEnviron {
		fmt.Printf("export %s\n", envLine)
	}
}
