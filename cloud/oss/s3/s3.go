package s3

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
	"github.com/boss-net/api/cloud/oss"
)

type S3Storage struct {
	bucket string
	client *s3.Client
}

func NewS3Storage(args oss.OSSArgs) (oss.OSS, error) {
	var err error
	if args.S3 == nil {
		return nil, oss.ErrArgumentInvalid.WithDetail("can't find s3 argument in OSSArgs")
	}
	err = args.S3.Validate()
	if err != nil {
		return nil, err
	}
	useAws := args.S3.UseAws
	ak := args.S3.AccessKey
	sk := args.S3.SecretKey
	region := args.S3.Region
	endpoint := args.S3.Endpoint
	usePathStyle := args.S3.UsePathStyle
	bucket := args.S3.Bucket
	useIamRole := args.S3.UseIamRole

	var cfg aws.Config
	var client *s3.Client

	if useAws {
		if (ak == "" && sk == "") || useIamRole {
			cfg, err = config.LoadDefaultConfig(
				context.TODO(),
				config.WithRegion(region),
			)
		} else {
			cfg, err = config.LoadDefaultConfig(
				context.TODO(),
				config.WithRegion(region),
				config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
					ak,
					sk,
					"",
				)),
			)
		}
		if err != nil {
			return nil, oss.ErrProviderInit.WithError(err)
		}

		client = s3.NewFromConfig(cfg, func(options *s3.Options) {
			if endpoint != "" {
				options.BaseEndpoint = aws.String(endpoint)
			}
			options.UsePathStyle = usePathStyle
		})
	} else {
		client = s3.New(s3.Options{
			Credentials:  credentials.NewStaticCredentialsProvider(ak, sk, ""),
			UsePathStyle: usePathStyle,
			Region:       region,
			EndpointResolver: s3.EndpointResolverFunc(
				func(region string, options s3.EndpointResolverOptions) (aws.Endpoint, error) {
					return aws.Endpoint{
						URL:               endpoint,
						HostnameImmutable: false,
						SigningName:       "s3",
						PartitionID:       "aws",
						SigningRegion:     region,
						SigningMethod:     "v4",
						Source:            aws.EndpointSourceCustom,
					}, nil
				}),
		})
	}

	// check bucket
	_, err = client.HeadBucket(context.TODO(), &s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "NotFound" {
			_, err = client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
				Bucket: aws.String(bucket),
			})
			if err != nil {
				return nil, oss.ErrProviderInit.WithError(err)
			}
		}
	}
	return &S3Storage{bucket: bucket, client: client}, nil
}

func (s *S3Storage) Save(key string, data []byte) error {
	_, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	})
	return err
}

func (s *S3Storage) Load(key string) ([]byte, error) {
	resp, err := s.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}

	return io.ReadAll(resp.Body)
}

func (s *S3Storage) Exists(key string) (bool, error) {
	_, err := s.client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	return err == nil, nil
}

func (s *S3Storage) Delete(key string) error {
	_, err := s.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	return err
}

func (s *S3Storage) List(prefix string) ([]oss.OSSPath, error) {
	// append a slash to the prefix if it doesn't end with one
	if !strings.HasSuffix(prefix, "/") {
		prefix = prefix + "/"
	}

	var keys []oss.OSSPath
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(prefix),
	}

	paginator := s3.NewListObjectsV2Paginator(s.client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, err
		}
		for _, obj := range page.Contents {
			// remove prefix
			key := strings.TrimPrefix(*obj.Key, prefix)
			// remove leading slash
			key = strings.TrimPrefix(key, "/")
			keys = append(keys, oss.OSSPath{
				Path:  key,
				IsDir: false,
			})
		}
	}

	return keys, nil
}

func (s *S3Storage) State(key string) (oss.OSSState, error) {
	resp, err := s.client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return oss.OSSState{}, err
	}

	if resp.ContentLength == nil {
		resp.ContentLength = ToPtr[int64](0)
	}
	if resp.LastModified == nil {
		resp.LastModified = ToPtr(time.Time{})
	}

	return oss.OSSState{
		Size:         *resp.ContentLength,
		LastModified: *resp.LastModified,
	}, nil
}

func (s *S3Storage) Type() string {
	return oss.OSS_TYPE_S3
}

func ToPtr[T any](value T) *T {
	return &value
}
