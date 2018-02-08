package main

import (
  "encoding/csv"
  "fmt"
  "io/ioutil"
  "os"
  "strconv"

  "gopkg.in/gographics/imagick.v2/imagick"

  "github.com/aprimadi/image-rotation-s3/config"
)

func resizeImage(wand *imagick.MagickWand, tw uint, th uint) {
  w := wand.GetImageWidth()
  h := wand.GetImageHeight()
  var tx int = 0
  var ty int = 0
  aw := tw
  ah := th
  if (w * th > tw * h) { // Width is larger
    aw = w * th / h
    tx = int((aw - tw) / 2)
  } else { // Height is larger
    ah = h * tw / w
    ty = int((ah - th) / 2)
  }
  wand.ResizeImage(aw, ah, imagick.FILTER_LANCZOS, 1)
  wand.CropImage(tw, th, tx, ty)
}

func main() {
  cfg, _ := config.ParseConfig()
  config.Cfg = cfg

  // Initialize imagick
  imagick.Initialize()
  defer imagick.Terminate()

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
    rotation, err := strconv.Atoi(records[1])
    if err != nil {
      panic(err)
    }

    // Generate key
    mil := imageId / 1000000
    tho := (imageId % 1000000) / 1000
    one := (imageId % 1000)
    baseDir := fmt.Sprintf("data/%03d/%03d/%03d", mil, tho, one)
    dir := fmt.Sprintf("%s/original", baseDir)

    files, err := ioutil.ReadDir(dir)
    if err != nil {
      panic(err)
    }
    pwand := imagick.NewPixelWand()
    for _, file := range files {
      filename := file.Name()
      // fmt.Println(file.Name())
      fmt.Println(fmt.Sprintf("%s/%s", dir, filename))
      wand := imagick.NewMagickWand()
      err = wand.ReadImage(fmt.Sprintf("%s/%s", dir, filename))
      if err != nil {
        panic(err)
      }

      var degrees float64
      switch rotation {
      case 90:
        degrees = 270
      case 180:
        degrees = 180
      case 270:
        degrees = 90
      default:
        panic("Invalid rotation")
      }
      // Rotate image
      wand.RotateImage(pwand, degrees)

      // Stores original image
      wand.WriteImage(fmt.Sprintf("%s/original/%s", baseDir, filename))

      // Create medium version
      os.MkdirAll(fmt.Sprintf("%s/medium", baseDir), os.ModePerm)
      switch wand.GetImageOrientation() {
      case imagick.ORIENTATION_UNDEFINED, imagick.ORIENTATION_TOP_LEFT, imagick.ORIENTATION_TOP_RIGHT, imagick.ORIENTATION_BOTTOM_LEFT, imagick.ORIENTATION_BOTTOM_RIGHT:
        resizeImage(wand, 700, 525)
      default:
        resizeImage(wand, 525, 700)
      }
      wand.WriteImage(fmt.Sprintf("%s/medium/%s", baseDir, filename))

      // Create small version
      os.MkdirAll(fmt.Sprintf("%s/small", baseDir), os.ModePerm)
      switch wand.GetImageOrientation() {
      case imagick.ORIENTATION_UNDEFINED, imagick.ORIENTATION_TOP_LEFT, imagick.ORIENTATION_TOP_RIGHT, imagick.ORIENTATION_BOTTOM_LEFT, imagick.ORIENTATION_BOTTOM_RIGHT:
        wand.ResizeImage(340, 255, imagick.FILTER_LANCZOS, 1)
      default:
        wand.ResizeImage(255, 340, imagick.FILTER_LANCZOS, 1)
      }
      wand.WriteImage(fmt.Sprintf("%s/small/%s", baseDir, filename))
    }
  }
}
