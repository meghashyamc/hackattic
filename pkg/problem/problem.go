package problem

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/meghashyamc/hackattic/pkg/httpclient"
)

func Get(problemString string, accessToken string, problemDetails any) error {
	ctx := context.Background()
	client := httpclient.GetDefaultClient()

	response, err := client.Get(ctx, fmt.Sprintf("/challenges/%s/problem?access_token=%s", problemString, accessToken))
	if err != nil {
		slog.Error("got an unexpected error when getting problem", "err", err)
		return err
	}

	err = json.Unmarshal(response.Body, &problemDetails)
	if err != nil {
		slog.Error("got an unexpected error when unmarshalling problem", "err", err)
		os.Exit(1)
	}

	slog.Info("fetched problem successfully")
	return nil
}

func Submit(problemString string, accessToken string, solution any) error {
	ctx := context.Background()
	client := httpclient.GetDefaultClient()

	response, err := client.Post(ctx, fmt.Sprintf("/challenges/%s/solve?access_token=%s", problemString, accessToken), solution, httpclient.WithHeaders(map[string]string{"Content-Type": "application/json"}))
	if err != nil {
		slog.Error("got an unexpected error when submitting solution", "err", err)
		return err
	}

	slog.Info("submitted solution", "status_code", response.StatusCode, "response", string(response.Body))

	return nil
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

	slog.Info("downloaded file successfully")

	return nil
}
