package model

type AwsS3 struct {
	Bucket    string `desc:"S3 bucket name"`
	FileName  string `desc:"File to upload or download"`
	Operation string `desc:"up|dwn upload or download file" flag:"op"`
	Region    string `desc:"AWS region the bucket is in e.g us-east-1"`
}
type Config struct {
	S3 AwsS3
}
