package problem

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/meghashyamc/hackattic/pkg/httpclient"
)

func Get(problemString string, accessToken string) ([]byte, error) {
	ctx := context.Background()
	client := httpclient.GetDefaultClient()

	response, err := client.Get(ctx, fmt.Sprintf("/challenges/%s/problem?access_token=%s", problemString, accessToken))
	if err != nil {
		slog.Error("got an unexpected error when getting problem", "err", err)
		return nil, err
	}

	return response.Body, nil
}

func Submit(problemString string, accessToken string, solution any) ([]byte, error) {
	ctx := context.Background()
	client := httpclient.GetDefaultClient()

	response, err := client.Post(ctx, fmt.Sprintf("/challenges/%s/solve?access_token=%s", problemString, accessToken), solution, httpclient.WithHeaders(map[string]string{"Content-Type": "application/json"}))
	if err != nil {
		slog.Error("got an unexpected error when submitting solution", "err", err)
		return nil, err
	}

	slog.Info("submitted solution", "status_code", response.StatusCode)

	return response.Body, nil
}

func DownloadFile(filepath string, url string) error {

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	ctx := context.Background()
	client := httpclient.GetDefaultClient()
	endpoint := strings.TrimPrefix(url, os.Getenv("BASE_URL"))
	response, err := client.Get(ctx, endpoint)
	if err != nil {
		slog.Error("got an unexpected error when downloading file", "err", err)
		return err
	}
	responseReader := bytes.NewReader(response.Body)
	_, err = io.Copy(out, responseReader)
	if err != nil {
		slog.Error("got an unexpected error when downloading file", "err", err)
		return err
	}

	return nil
}
