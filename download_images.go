package main

import (
  "encoding/csv"
  "fmt"
  "os"
  "path"
  "strconv"

  "gopkg.in/gographics/imagick.v2/imagick"

  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/service/s3"
  "github.com/aws/aws-sdk-go/service/s3/s3manager"

  irAws "github.com/aprimadi/image-rotation-s3/aws"
  "github.com/aprimadi/image-rotation-s3/config"
)

func main() {
  cfg, _ := config.ParseConfig()
  config.Cfg = cfg

  // Read S3 config
  irAws.S3Cfg = irAws.S3LoadConfig(cfg)

  sess := irAws.GetSession()
  downloader := s3manager.NewDownloader(sess)
  svc := s3.New(sess)

  file, err := os.Open("result.csv")
  if err != nil {
    panic(err)
  }

  reader := csv.NewReader(file)
  for {
    records, err := reader.Read()
    if err != nil {
      break
    }
    imageId, err := strconv.Atoi(records[0])
    if err != nil {
      panic(err)
    }

    // Generate key
    mil := imageId / 1000000
    tho := (imageId % 1000000) / 1000
    one := (imageId % 1000)
    dir := fmt.Sprintf("%03d/%03d/%03d/original", mil, tho, one)
    result, err := svc.ListObjects(&s3.ListObjectsInput{
      Bucket: aws.String(irAws.S3Cfg.Bucket),
      Prefix: aws.String(fmt.Sprintf("appointments/completed_works/%s", dir)),
    })
    if err != nil {
      panic(err)
    }
    _, file := path.Split(*result.Contents[0].Key)
    key := fmt.Sprintf("appointments/completed_works/%s/%s", dir, file)
    localKey := fmt.Sprintf("data/%s/%s", dir, file)
    fmt.Println(key)

    // Download image
    os.MkdirAll(fmt.Sprintf("data/%s", dir), os.ModePerm)
    f, err := os.Create(localKey)
    downloader.Download(
      f,
      &s3.GetObjectInput{
        Bucket: aws.String(irAws.S3Cfg.Bucket),
        Key: aws.String(key),
      },
    )
  }
}
