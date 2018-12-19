package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/mexisme/get-secrets/config"
	log "github.com/sirupsen/logrus"
)

//nolint:gochecknoinits
func init() {
	config.AddConfigItems([]string{"aws.region"})
}

// AwsRegion is the default AWS Region
// TODO: `Region` should be in a config file (or ~/.aws/config) or in the s3/url package?
const AwsRegion = "ap-southeast-2"

// Details for getting Secret files
type Details struct {
	config  *aws.Config
	session *session.Session
}

// New object
func New() *Details {
	return (&Details{}).WithRegion(AwsRegion)
}

// WithRegion creates new struct with `config` updated with the AWS Region
func (s *Details) WithRegion(region string) *Details {
	clone := *s // This does a shallow clone

	clone.config = &aws.Config{Region: aws.String(region)}

	return &clone
}

// WithSession creates new struct with `session` updated
func (s *Details) WithSession(session *session.Session) *Details {
	clone := *s // This does a shallow clone

	clone.session = session

	return &clone
}

// Config returns the configured `config` field
func (s *Details) Config() *aws.Config {
	return s.config
}

// Session returns the contents of the `session` field
func (s *Details) Session() *session.Session {
	if s.session == nil {
		s.session = s.newSession()
	}

	return s.session
}

func (s *Details) newSession() *session.Session {
	// TODO: Enable AWS_SDK_LOAD_CONFIG env-var, somehow!
	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
		SharedConfigState:       session.SharedConfigEnable,
	}))
	log.WithFields(log.Fields{"session": awsSession}).Debug("Created new AWS Session")

	return awsSession
}
