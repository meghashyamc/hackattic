package main

import (
	"encoding/json"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/meghashyamc/hackattic/problem"
	"github.com/meghashyamc/hackattic/validation"
)

type Problem struct {
	ImageURL string `json:"image_url"`
}

type Solution struct {
	Code string `json:"code"`
}

func readRotatedQR(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return "", err
	}

	qrReader := qrcode.NewQRCodeReader()
	hints := map[gozxing.DecodeHintType]interface{}{
		gozxing.DecodeHintType_TRY_HARDER: true,
	}

	// Decode QR code
	result, err := qrReader.Decode(bmp, hints)
	if err != nil {
		return "", err
	}

	return result.GetText(), nil
}

func main() {
	godotenv.Load("../.env")
	problemName := "reading_qr"
	qrFileName := "rotated_qr.png"
	accessToken := os.Getenv("ACCESS_TOKEN")
	if err := validation.ValidateAccessToken(accessToken); err != nil {
		os.Exit(1)
	}

	response, err := problem.Get(problemName, accessToken)
	if err != nil {
		os.Exit(1)
	}
	var problemDetails Problem
	err = json.Unmarshal(response, &problemDetails)
	if err != nil {
		slog.Error("got an unexpected error when unmarshalling problem", "err", err)
		os.Exit(1)
	}
	err = problem.DownloadFile(qrFileName, problemDetails.ImageURL)
	if err != nil {
		os.Exit(1)
	}

	defer os.Remove(qrFileName)

	code, err := readRotatedQR(qrFileName)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("found code", "code", code)
}
