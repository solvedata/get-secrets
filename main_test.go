package main_test

import (
	// TODO: Add this back when some tests are written:
	// . "github.com/mexisme/get-secrets"

	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var _ = Describe("The main package", func() {
	var s3PathEnvVarName = "SECRETS_S3_DOTENV_PATH"

	It("setting env-var gets read by Viper", func() {
		if err := viper.BindEnv("testmain"); err != nil {
			log.WithField("Error", err).Panic(`When viper.BindEnv("testmain")`)
		}
		os.Setenv("SECRETS_TESTMAIN", "test-secrets")
		Expect(viper.GetString("testmain")).To(Equal("test-secrets"))
	})

	It("setting s3.dotenv_path passes-into Viper", func() {
		envVal, envExists := os.LookupEnv(s3PathEnvVarName)
		os.Setenv(s3PathEnvVarName, "test-secrets")

		Expect(viper.GetString("s3.dotenv_path")).To(Equal("test-secrets"), "Env Var name = %#v", s3PathEnvVarName)
		if envExists {
			os.Setenv(s3PathEnvVarName, envVal)
		} else {
			os.Unsetenv(s3PathEnvVarName)
		}
	})
})
