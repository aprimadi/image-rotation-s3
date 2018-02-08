package aws

import (
  "bytes"
  "fmt"
  "io/ioutil"
  "path"

  "github.com/spf13/viper"
  awsSession "github.com/aws/aws-sdk-go/aws/session"

  "github.com/aprimadi/image-rotation-s3/config"
)

type S3Config struct {
  Bucket string
}

var S3Cfg *S3Config
var AWS *awsSession.Session

func GetSession() *awsSession.Session {
  if AWS == nil {
    AWS = awsSession.Must(awsSession.NewSession())
  }
  return AWS
}

func S3LoadConfig(cfg *config.Config) *S3Config {
  viper.SetConfigType("YAML")

  dat, err := ioutil.ReadFile(path.Join(cfg.BaseDir, "config/s3.yml"))
  if err != nil {
    panic(err)
  }

  viper.ReadConfig(bytes.NewReader(dat))
  bucket := viper.GetString(fmt.Sprintf("%s.%s", cfg.Environment, "bucket"))

  s3Cfg := S3Config{
    Bucket: bucket,
  }

  return &s3Cfg
}
