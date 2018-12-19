package dotenv

import (
	"github.com/mexisme/get-secrets/errors"
	s3ish "github.com/mexisme/get-secrets/files/s3"
	urlish "github.com/mexisme/get-secrets/files/s3/s3url"

	"github.com/mexisme/multiconfig"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// ReadFromS3 TODO
func ReadFromS3(envs *multiconfig.Map) {
	// TODO: Extract file-reading
	s3Path := viper.GetString("s3.dotenv_path")
	s3url := urlish.New().WithURL(s3Path)
	log.Infof("S3 .env Base path = %#v (%#v)", s3Path, s3url)

	s3 := s3ish.New().WithSource(s3url)

	{
		defer log.Warn("NB: Did you perhaps intend to enable 'dotenv.skip' ($SKIP_SECRETS)?")

		s3lists, err := s3.List()
		if err != nil {
			errors.PanicOnErrors(err)
		}

		if err := s3.ReadListToCallback(s3lists, EnvAddConfig(envs)); err != nil {
			errors.PanicOnErrors(err)
		}
	}
}
