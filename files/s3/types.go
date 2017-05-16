package s3

import (
	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// AwsRegion is the default AWS Region
const AwsRegion = "ap-southeast-2"

var (
	// AwsConfig is the default AWS Config
	AwsConfig = &aws.Config{Region: aws.String(AwsRegion)}
)

// bucketPrefix is the required interface for "Source" attr
type bucketPrefix interface {
	Bucket() string
	Prefix() string
}

// Details for getting Secret files
type Details struct {
	source    bucketPrefix
	awsConfig *aws.Config
	s3Session *s3.S3
	//session *session.Session
}

// New object
func New() *Details {
	return &Details{}
}

// WithSource creates new struct with `source` updated
func (s *Details) WithSource(source bucketPrefix) *Details {
	clone := *s // This does a shallow clone

	clone.source = source

	var err error
	if clone.s3Session, err = s.newS3Session(); err != nil {
		log.Panic(err)
	}

	return &clone
}

// WithS3Session creates new struct with `s3Session` updated
func (s *Details) WithS3Session(s3Session *s3.S3) *Details {
	clone := *s // This does a shallow clone

	clone.s3Session = s3Session

	return &clone
}

// S3 returns the
func (s *Details) S3() *s3.S3 {
	return s.s3Session
}

func (s *Details) newS3Session() (*s3.S3, error) {
	// TODO: Enable AWS_SDK_LOAD_CONFIG env-var, somehow!
	session, err := session.NewSession()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// TODO: `Region` should be in a config file (or ~/.aws/config) or in the s3/url package?
	s3Session := s3.New(session, s.awsConfig)
	log.WithFields(log.Fields{"s3Session": s3Session}).Debug("Created new S3 Session")

	return s3Session, nil
}
