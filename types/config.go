package types

type AppConfig struct {
	Port      string      `yaml:"port"`
	CydresUrl string      `yaml:"cydresUrl"`
	AwsS3     AwsS3Config `yaml:"aws"`
}
type AwsS3Config struct {
	BucketName string `yaml:"bucket_name"`
	S3Url      string `yaml:"s3Url"`
	AccessKey  string `yaml:"access_key"`
	SecretKey  string `yaml:"secret_key"`
	Region     string `yaml:"region"`
	Folder     string `yaml:"folder_name"`
}
