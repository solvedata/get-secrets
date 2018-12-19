package s3

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/mexisme/get-secrets/config"
	configAws "github.com/mexisme/get-secrets/config/aws"
	log "github.com/sirupsen/logrus"
)

//nolint:gochecknoinits
func init() {
	config.AddConfigItems([]string{"s3.dotenv_path"})
}

// bucketPrefix is the required interface for "Source" attr
type bucketPrefix interface {
	Bucket() string
	Prefix() string
	PrefixDir() string
	JoinPath(string) string
}

// Details for getting Secret files
type Details struct {
	source    bucketPrefix
	awsConfig *configAws.Details
	s3Session *s3.S3
}

// New object
func New() *Details {
	return (&Details{}).WithAwsConfig(configAws.New())
}

// WithAwsConfig creates new struct with `source` updated
func (s *Details) WithAwsConfig(awsConfig *configAws.Details) *Details {
	clone := *s // This does a shallow clone

	clone.awsConfig = awsConfig

	return &clone
}

// WithSource creates new struct with `source` updated
func (s *Details) WithSource(source bucketPrefix) *Details {
	clone := *s // This does a shallow clone

	clone.source = source

	if clone.s3Session == nil {
		var err error
		if clone.s3Session, err = s.newS3Session(); err != nil {
			log.Panic(err)
		}
	}

	return &clone
}

// WithS3Session creates new struct with `s3Session` updated
func (s *Details) WithS3Session(s3Session *s3.S3) *Details {
	clone := *s // This does a shallow clone

	clone.s3Session = s3Session

	return &clone
}

// S3 returns the S3 Session property
func (s *Details) S3() *s3.S3 {
	return s.s3Session
}

// TODO: 'error' is unused:
//nolint:unparam
func (s *Details) newS3Session() (*s3.S3, error) {
	s3Session := s3.New(s.awsConfig.Session(), s.awsConfig.Config())
	log.WithFields(log.Fields{"s3Session": s3Session}).Debug("Created new S3 Session")

	return s3Session, nil
}
