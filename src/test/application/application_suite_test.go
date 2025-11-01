package application_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestApplication(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Application Suite")
}

// Mock Mailersend Email API
func mockMailersend(captures *[]http.Request) {
	httpmock.RegisterResponder(
		"POST",
		"https://api.mailersend.com/v1/email",
		func(req *http.Request) (*http.Response, error) {
			if captures != nil {
				*captures = append(*captures, *req)
			}
			return httpmock.NewStringResponse(200, "OK"), nil
		},
	)
}

func mockGoogleRecaptcha(captures *[]http.Request) {
	httpmock.RegisterResponder(
		"POST",
		"https://www.google.com/recaptcha/api/siteverify",
		func(req *http.Request) (*http.Response, error) {
			if captures != nil {
				*captures = append(*captures, *req)
			}
			return httpmock.NewStringResponse(200, `{"success": true}`), nil
		},
	)
}

func waitForReady(
	ctx context.Context,
	timeout time.Duration,
	endpoint string,
) error {
	client := http.Client{}
	startTime := time.Now()
	for {
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			endpoint,
			nil,
		)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			// fmt.Printf("Error making request: %s\n", err.Error())
			continue
		}
		if resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return nil
		}
		resp.Body.Close()

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if time.Since(startTime) >= timeout {
				return fmt.Errorf("timeout reached while waiting for endpoint")
			}
			// wait a little while between checks
			time.Sleep(250 * time.Millisecond)
		}
	}
}
