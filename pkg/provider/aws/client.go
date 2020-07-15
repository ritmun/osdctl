package aws

// Generate client mocks for testing
//go:generate mockgen -source=client.go -package=mock -destination=mock/client.go

import (
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"

	"github.com/pkg/errors"
)

// AwsClientInput input for new aws client
type AwsClientInput struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
	Region          string
}

// TODO: Add more methods when needed
type Client interface {
	// sts
	AssumeRole(*sts.AssumeRoleInput) (*sts.AssumeRoleOutput, error)
	GetCallerIdentity(*sts.GetCallerIdentityInput) (*sts.GetCallerIdentityOutput, error)
	GetFederationToken(*sts.GetFederationTokenInput) (*sts.GetFederationTokenOutput, error)

	// S3
	ListBuckets(*s3.ListBucketsInput) (*s3.ListBucketsOutput, error)
	DeleteBucket(*s3.DeleteBucketInput) (*s3.DeleteBucketOutput, error)
	ListObjects(*s3.ListObjectsInput) (*s3.ListObjectsOutput, error)
	DeleteObjects(*s3.DeleteObjectsInput) (*s3.DeleteObjectsOutput, error)

	//iam
	CreateAccessKey(*iam.CreateAccessKeyInput) (*iam.CreateAccessKeyOutput, error)
	DeleteAccessKey(*iam.DeleteAccessKeyInput) (*iam.DeleteAccessKeyOutput, error)
	ListAccessKeys(*iam.ListAccessKeysInput) (*iam.ListAccessKeysOutput, error)
	GetUser(*iam.GetUserInput) (*iam.GetUserOutput, error)
	CreateUser(*iam.CreateUserInput) (*iam.CreateUserOutput, error)
	ListUsers(*iam.ListUsersInput) (*iam.ListUsersOutput, error)
	AttachUserPolicy(*iam.AttachUserPolicyInput) (*iam.AttachUserPolicyOutput, error)
}

type AwsClient struct {
	iamClient iamiface.IAMAPI
	stsClient stsiface.STSAPI
	s3Client  s3iface.S3API
}

// NewAwsClient creates an AWS client with credentials in the environment
func NewAwsClient(profile, region, configFile string) (Client, error) {
	opt := session.Options{
		Config: aws.Config{
			Region: aws.String(region),
		},
		Profile: profile,
	}

	// only set config file if it is not empty
	if configFile != "" {
		absCfgPath, err := filepath.Abs(configFile)
		if err != nil {
			return nil, err
		}
		opt.SharedConfigFiles = []string{absCfgPath}
	}

	sess := session.Must(session.NewSessionWithOptions(opt))
	_, err := sess.Config.Credentials.Get()

	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case "NoCredentialProviders":
			return nil, errors.Wrap(err, "Could not create AWS session")
		default:
			return nil, errors.Wrap(err, "Could not create AWS session")
		}
	}

	return &AwsClient{
		iamClient: iam.New(sess),
		stsClient: sts.New(sess),
		s3Client:  s3.New(sess),
	}, nil
}

// NewAwsClientWithInput creates an AWS client with input credentials
func NewAwsClientWithInput(input *AwsClientInput) (Client, error) {
	config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(input.AccessKeyID, input.SecretAccessKey, input.SessionToken),
		Region:      aws.String(input.Region),
	}

	s, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}

	return &AwsClient{
		iamClient: iam.New(s),
		stsClient: sts.New(s),
		s3Client:  s3.New(s),
	}, nil
}

func (c *AwsClient) AssumeRole(input *sts.AssumeRoleInput) (*sts.AssumeRoleOutput, error) {
	return c.stsClient.AssumeRole(input)
}

func (c *AwsClient) GetCallerIdentity(input *sts.GetCallerIdentityInput) (*sts.GetCallerIdentityOutput, error) {
	return c.stsClient.GetCallerIdentity(input)
}

func (c *AwsClient) GetFederationToken(input *sts.GetFederationTokenInput) (*sts.GetFederationTokenOutput, error) {
	return c.stsClient.GetFederationToken(input)
}

func (c *AwsClient) ListBuckets(input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error) {
	return c.s3Client.ListBuckets(input)
}

func (c *AwsClient) DeleteBucket(input *s3.DeleteBucketInput) (*s3.DeleteBucketOutput, error) {
	return c.s3Client.DeleteBucket(input)
}

func (c *AwsClient) ListObjects(input *s3.ListObjectsInput) (*s3.ListObjectsOutput, error) {
	return c.s3Client.ListObjects(input)
}

func (c *AwsClient) DeleteObjects(input *s3.DeleteObjectsInput) (*s3.DeleteObjectsOutput, error) {
	return c.s3Client.DeleteObjects(input)
}

func (c *AwsClient) CreateAccessKey(input *iam.CreateAccessKeyInput) (*iam.CreateAccessKeyOutput, error) {
	return c.iamClient.CreateAccessKey(input)
}

func (c *AwsClient) DeleteAccessKey(input *iam.DeleteAccessKeyInput) (*iam.DeleteAccessKeyOutput, error) {
	return c.iamClient.DeleteAccessKey(input)
}

func (c *AwsClient) ListAccessKeys(input *iam.ListAccessKeysInput) (*iam.ListAccessKeysOutput, error) {
	return c.iamClient.ListAccessKeys(input)
}

func (c *AwsClient) GetUser(input *iam.GetUserInput) (*iam.GetUserOutput, error) {
	return c.iamClient.GetUser(input)
}

func (c *AwsClient) CreateUser(input *iam.CreateUserInput) (*iam.CreateUserOutput, error) {
	return c.iamClient.CreateUser(input)
}

func (c *AwsClient) ListUsers(input *iam.ListUsersInput) (*iam.ListUsersOutput, error) {
	return c.iamClient.ListUsers(input)
}

func (c *AwsClient) AttachUserPolicy(input *iam.AttachUserPolicyInput) (*iam.AttachUserPolicyOutput, error) {
	return c.iamClient.AttachUserPolicy(input)
}
