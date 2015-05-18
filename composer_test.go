package composer

import (
	"bytes"
	"io"
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

		loader := func(url string) io.Reader {
			if url == "http://composer/a" {
				return strings.NewReader("Hello, client from <div composer-url=\"http://composer/b\"></div>")
			} else if url == "http://composer/b" {
				return strings.NewReader("Hello, client from b")
			}
			return strings.NewReader("Hello, " + url)
		}

		It("Should replace composer elements with loaded HTML", func() {
			composed := Compose(strings.NewReader("<div><p composer-url=\"http://composer/b\" title=\"captain\">other stuff</p></div>"), loader)
			buf := new(bytes.Buffer)
			buf.ReadFrom(composed)

			Expect(buf.String()).To(Equal("<div>Hello, client from b</div>"))
		})

		It("Should recursively process composer elements found in loaded HTML", func() {
			composed := Compose(strings.NewReader("<div><p composer-url=\"http://composer/a\" title=\"captain\">other stuff</p></div>"), loader)
			buf := new(bytes.Buffer)
			buf.ReadFrom(composed)

			Expect(buf.String()).To(Equal("<div>Hello, client from Hello, client from b</div>"))
		})
	})
})
