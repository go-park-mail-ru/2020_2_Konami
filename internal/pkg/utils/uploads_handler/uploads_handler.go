package uploads_handler

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"strings"
)

type UploadsHandler struct {
	uploadsDir string
}

func NewUploadsHandler(uploadsDir string) UploadsHandler {
	return UploadsHandler{uploadsDir: uploadsDir}
}

func (h UploadsHandler) SavePng(imgPath string, img *image.Image) (string, error) {
	imgPath = h.uploadsDir + "/" + imgPath + ".png"
	f, err := os.OpenFile(imgPath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return "", err
	}
	defer f.Close()
	err = png.Encode(f, *img)
	if err != nil {
		return "", err
	}
	return imgPath, nil
}

func (h UploadsHandler) DecodeBase64Image(encoded *string) (*image.Image, error) {
	jpegPrefix := "data:image/jpeg;base64,"
	pngPrefix := "data:image/png;base64,"
	var img image.Image
	var err error
	rawImage := *encoded
	if strings.HasPrefix(rawImage, jpegPrefix) {
		rawImage = rawImage[len(jpegPrefix):]
		decoded, _ := base64.StdEncoding.DecodeString(rawImage)
		img, err = jpeg.Decode(bytes.NewReader(decoded))
	} else if strings.HasPrefix(rawImage, pngPrefix) {
		rawImage := rawImage[len(pngPrefix):]
		decoded, _ := base64.StdEncoding.DecodeString(rawImage)
		img, err = png.Decode(bytes.NewReader(decoded))
	} else {
		return nil, errors.New("invalid image encoding")
	}
	return &img, err
}

func (h UploadsHandler) UploadBase64Image(imgPath string, encoded *string) (string, error) {
	img, err := h.DecodeBase64Image(encoded)
	if err != nil {
		return "", err
	}
	return h.SavePng(imgPath, img)
}

func (h UploadsHandler) UploadImage(imgPath string, img io.Reader) (string, error) {
	fnameSlice := strings.Split(imgPath, ".")
	ext := fnameSlice[len(fnameSlice)-1]
	if ext != "jpg" && ext != "jpeg" && ext != "png" {
		return "", errors.New("invalid file format")
	}
	imgPath = h.uploadsDir + "/" + imgPath
	f, err := os.OpenFile(imgPath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return "", err
	}
	defer f.Close()
	var written int64 = 0
	written, err = io.Copy(f, img)
	if err != nil {
		return "", err
	}
	if written == 0 {
		return "", errors.New("unable to write file")
	}
	return imgPath, nil
}
