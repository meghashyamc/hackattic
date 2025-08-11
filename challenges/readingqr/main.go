package main

// Problem: https://hackattic.com/challenges/reading_qr
import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/meghashyamc/hackattic/pkg/auth"
	"github.com/meghashyamc/hackattic/pkg/problem"
)

type Problem struct {
	ImageURL string `json:"image_url"`
}

type Solution struct {
	Code string `json:"code"`
}

func main() {

	accessToken, err := auth.GetAccessToken()
	if err != nil {
		os.Exit(1)
	}
	problemName := "reading_qr"
	qrFileName := "rotated_qr.png"
	var problemDetails Problem

	if err := problem.Get(problemName, accessToken, problemDetails); err != nil {
		os.Exit(1)
	}

	err = problem.DownloadFile(qrFileName, problemDetails.ImageURL)
	if err != nil {
		os.Exit(1)
	}

	defer os.Remove(qrFileName)

	code, err := readRotatedQR(qrFileName)
	if err != nil {
		os.Exit(1)
	}
	solution := Solution{Code: code}
	if err := problem.Submit(problemName, accessToken, solution); err != nil {
		os.Exit(1)
	}

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
