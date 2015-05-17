package composer

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/net/html"
)

var _ = Describe("Composer", func() {

	Context("Tokenizing HTML", func() {
		htmlWithComposerAttribute := "<div composer-url=\"http://example.com\" title=\"captain\"></div>"
		htmlWithoutComposerAttribute := "<div name=\"brendan\" title=\"captain\"></div>"

		It("Should retrieve all html attributes", func() {

			z := html.NewTokenizer(strings.NewReader(htmlWithoutComposerAttribute))
			z.Next()

			attributes := getAttributes(z, make([]*html.Attribute, 0))

			Expect(len(attributes)).To(Equal(2))
			Expect(attributes[0].Key).To(Equal("name"))
		})

		It("Should retrieve properly initialized ComposerTag from HTML", func() {
			z := html.NewTokenizer(strings.NewReader(htmlWithComposerAttribute))
			z.Next()

			t := getComposerTag(z)

			Expect(t).NotTo(BeNil())
			Expect(t.Url).To(Equal("http://example.com"))
		})

		It("Should not return a ComposerTag if composer attributes are not present", func() {
			z := html.NewTokenizer(strings.NewReader(htmlWithoutComposerAttribute))
			z.Next()

			t := getComposerTag(z)

			Expect(t).To(BeNil())
		})

	})

	Context("Composing HTML", func() {
		var ts *httptest.Server
		var serverUrl string

		BeforeEach(func() {
			ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/a" {
					fmt.Fprintf(w, "Hello, client from <div composer-url=\"%s/b\"></div>", serverUrl)
				} else if r.URL.Path == "/b" {
					fmt.Fprint(w, "Hello, client from b")
				} else {
					fmt.Fprint(w, "Hello, "+r.URL.Path)
				}
			}))

			serverUrl = ts.URL
		})

		AfterEach(func() {
			ts.Close()
		})

		It("Should replace composer elements with loaded HTML", func() {
			composed := Compose(strings.NewReader(fmt.Sprintf("<div><p composer-url=\"%s/b\" title=\"captain\">other stuff</p></div>", ts.URL)))
			buf := new(bytes.Buffer)
			buf.ReadFrom(composed)

			Expect(buf.String()).To(Equal("<div>Hello, client from b</div>"))
		})

		It("Should recursively process composer elements found in loaded HTML", func() {
			composed := Compose(strings.NewReader(fmt.Sprintf("<div><p composer-url=\"%s/a\" title=\"captain\">other stuff</p></div>", ts.URL)))
			buf := new(bytes.Buffer)
			buf.ReadFrom(composed)

			Expect(buf.String()).To(Equal("<div>Hello, client from Hello, client from b</div>"))
		})
	})
})
