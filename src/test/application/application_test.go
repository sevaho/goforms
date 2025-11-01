package application_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
	app "github.com/sevaho/goforms/src"
	"github.com/sevaho/goforms/src/config"
	"github.com/tidwall/gjson"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var client = resty.New()
var port = 30001
var testApp = fmt.Sprintf("http://localhost:%d", port)

var formWithMailerSendID = uuid.NewString()
var formWithFakeBackendSendID = uuid.NewString()

var formsconfig = []byte(`
forms:
  - id: "` + formWithMailerSendID + `"
    provider: mailersend
    name: Contact TTC Teneramonda
    subject: Contact formulier website TTC Teneramonda
    sender:
      email: noreply@ttcteneramonda.be
      name: TTC Teneramonda website
    recipients:
    - email: ttcteneramonda@outlook.com
      name: TTC Teneramonda website
  - id: "` + formWithFakeBackendSendID + `"
    provider: fake
    name: Contact TTC Teneramonda
    subject: Contact formulier website TTC Teneramonda
    sender:
      email: noreply@ttcteneramonda.be
      name: TTC Teneramonda website
    recipients:
    - email: ttcteneramonda@outlook.com
      name: TTC Teneramonda website
`)

var _ = Context("Application", func() {
	var ctx context.Context
	var cancel context.CancelFunc
	var wg sync.WaitGroup
	var env *config.Config

	BeforeEach(func() {
		// setup the environment
		env = config.New()
		env.FORMS_CONFIG_BASE64 = base64.StdEncoding.EncodeToString(formsconfig)
		env.LOG_LEVEL = 4
		env.IS_DEVELOPMENT = false // Force use of embedded templates for tests

		// setup application
		ctx, cancel = context.WithCancel(context.Background())
		wg.Add(1)
		go func() { defer wg.Done(); app.Run(ctx, port, env) }()

		// Wait till application is ready to test
		waitForReady(ctx, time.Second*2, testApp+"/healthz")
	})
	AfterEach(func() {
		cancel()
		wg.Wait()
	})

	JustBeforeEach(func() {
		httpmock.Activate()
	})
	JustAfterEach(func() {
		httpmock.DeactivateAndReset()
	})

	When("Fetching the index page", func() {
		It("should render the correct html page", func() {
			// when
			res, err := client.R().Get(testApp)

			// then
			Expect(err).To(BeNil())
			Expect(res.StatusCode()).To(Equal(200), res.String())
		})
	})

	When("Fetching the templates should return 404 as development is off", func() {
		It("should return the correct html template", func() {
			// when
			res, err := client.R().Get(testApp + "/templates/email/contact")

			// then
			Expect(err).To(BeNil())
			Expect(res.StatusCode()).To(Equal(404), res.String())
		})
	})

	When("Fetching mails via API", func() {
		It("should return error when not authenticated", func() {
			// when
			res, err := client.R().Get(testApp + "/api/mails")

			// then
			Expect(err).To(BeNil())
			Expect(res.StatusCode()).To(Equal(400))
		})
	})

	When("Doing a form inquiry with mailersend backend", func() {
		var response *resty.Response
		var captures *[]http.Request

		JustBeforeEach(func() {

			var mocks []http.Request
			mockGoogleRecaptcha(&mocks)
			mockMailersend(&mocks)
			captures = &mocks

			var formData = map[string]string{"name": "John Doe", "age": "30"}
			res, err := client.R().SetFormData(formData).Post(testApp + "/forms/" + formWithMailerSendID)
			Expect(err).To(BeNil())
			response = res
		})

		It("should show a success page", func() {
			// then
			Expect(response.StatusCode()).To(Equal(200), response.String())
		})

		It("should store the mail in database and be retrieveable", func() {
			// when
			res, err := client.R().SetHeader("Authorization", "Bearer "+env.API_KEY).Get(testApp + "/api/mails")
			// then
			Expect(err).To(BeNil())
			Expect(res.StatusCode()).To(Equal(200), res.String())
			Expect(gjson.Get(res.String(), "count").Int()).To(Equal(int64(1)), res.String())

		})

		It("should do an API call to Google captcha and mailersend", func() {
			Expect(captures).ToNot(BeNil())
			Expect(len(*captures)).To(Equal(2))
		})
	})

	When("Doing a form inquiry with fake backend", func() {
		var response *resty.Response
		var captures *[]http.Request

		JustBeforeEach(func() {
			var mocks []http.Request
			mockGoogleRecaptcha(&mocks)
			captures = &mocks

			var formData = map[string]string{"name": "John Doe", "age": "30"}
			res, err := client.R().SetFormData(formData).Post(testApp + "/forms/" + formWithFakeBackendSendID)
			Expect(err).To(BeNil())
			response = res
		})

		It("should show a success page", func() {
			// then
			Expect(response.StatusCode()).To(Equal(200), response.String())
		})

		It("should store the mail in database and be retrieveable", func() {
			// when
			res, err := client.R().SetHeader("Authorization", "Bearer "+env.API_KEY).Get(testApp + "/api/mails")
			// then
			Expect(err).To(BeNil())
			Expect(res.StatusCode()).To(Equal(200), res.String())
			Expect(gjson.Get(res.String(), "count").Int()).To(Equal(int64(1)), res.String())

		})

		It("should do an API call to Google captcha", func() {
			Expect(captures).ToNot(BeNil())
			Expect(len(*captures)).To(Equal(1))
		})
	})

	When("Doing a form inquiry with EN language", func() {
		It("should respond in english", func() {
			// setup
			mockGoogleRecaptcha(nil)

			// given
			var formData = map[string]string{"name": "John Doe", "age": "30"}
			languageParam := "en"

			// when
			res, err := client.R().SetFormData(formData).SetQueryParam("language", languageParam).Post(testApp + "/forms/" + formWithFakeBackendSendID)

			// then
			Expect(err).To(BeNil())
			Expect(res.String()).To(ContainSubstring("html"))
			Expect(res.String()).To(ContainSubstring("Return to website"))
		})
	})

	When("Doing a form inquiry with NL language", func() {
		It("should respond in dutch", func() {
			// setup
			mockGoogleRecaptcha(nil)

			// given
			var formData = map[string]string{"name": "John Doe", "age": "30"}
			languageParam := "nl"

			// when
			res, err := client.R().SetFormData(formData).SetQueryParam("language", languageParam).Post(testApp + "/forms/" + formWithFakeBackendSendID)

			// then
			Expect(err).To(BeNil())
			Expect(res.String()).To(ContainSubstring("html"))
			Expect(res.String()).To(ContainSubstring("Terug naar website"))
		})
	})
})
