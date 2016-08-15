package uaago_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http/httptest"
	"net/http"
	"github.com/cloudfoundry-incubator/uaago"
	"io/ioutil"
)

var _ = Describe("Client", func() {
	Context("CheckTokenIsValid", func() {
		var (
			uaaTestServer *httptest.Server
			uaaRequests = make(chan *http.Request, 10)
			uaaResponseBodies = make(chan string, 10)
			client *uaago.Client
		)

		BeforeEach(func() {
			uaaTestServer = httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				uaaRequests <- request

				writer.WriteHeader(http.StatusOK)
				writer.Write([]byte(<-uaaResponseBodies))
			}))
			client, _ = uaago.NewClient(uaaTestServer.URL)
		})

		AfterEach(func() {
			uaaTestServer.Close()
		})

		FIt("talks to UAA", func() {
			uaaResponseBodies <- ""

			client.CheckTokenIsValid("some token", "some client_id")

			var request *http.Request
			Eventually(uaaRequests).Should(Receive(&request))
			Expect(request.Method).To(Equal("POST"))
			Expect(request.URL).To(ContainSubstring("/check_token"))

			requestBody, _ := ioutil.ReadAll(request.Body)
			Expect(requestBody).To(ContainSubstring("some token"))

			var response *http.Response
			Eventually(uaaResponseBodies).Should(Receive(&response))
			responseBody, _ := ioutil.ReadAll(response.Body)
			Expect(responseBody).To(ContainSubstring("some client_id"))

		})

		Context("valid: client_id=ingestor", func() {
			It("returns true", func() {
				// spin up a testserver
				// testserver returns 200 and body with client_id="foo"
				// send in token, client_id="foo"
				// check return == true
			})
		})

		Context("invalid: client_id=foo", func() {
			It("returns false", func() {
				// spin up a testserver
				// testserver returns 200 and body with client_id="bar"
				// send in token, client_id="foo"
				// check return == false
			})
		})
	})
})

/*
	UAA request:

	POST /check_token HTTP/1.1
	Host: server.example.com
	Authorization: Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==
	Content-Type: application/x-www-form-urlencoded

	token=eyJ0eXAiOiJKV1QiL


	UAA response:

	HTTP/1.1 200 OK
	Content-Type: application/json

	{
		"jti":"4657c1a8-b2d0-4304-b1fe-7bdc203d944f",
		"aud":["openid","cloud_controller"],
		"scope":["read"],
		"email":"marissa@test.org",
		"exp":138943173,
		"user_id":"41750ae1-b2d0-4304-b1fe-7bdc24256387",
		"user_name":"marissa",
		"client_id":"cf"
	}
*/
