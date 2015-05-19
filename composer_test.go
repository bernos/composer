package composer_test

import (
	"bytes"
	"io"
	"strings"

	. "github.com/bernos/composer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Composer", func() {

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

		It("Should replace multiple composer elements with loaded HTML", func() {
			composed := Compose(strings.NewReader("<div><p composer-url=\"http://composer/b\" title=\"captain\">other stuff</p></div><p><p composer-url=\"http://composer/b\"></p></p>"), loader)
			buf := new(bytes.Buffer)
			buf.ReadFrom(composed)

			Expect(buf.String()).To(Equal("<div>Hello, client from b</div><p>Hello, client from b</p>"))
		})

		It("Should recursively process composer elements found in loaded HTML", func() {
			composed := Compose(strings.NewReader("<div><p composer-url=\"http://composer/a\" title=\"captain\">other stuff</p></div>"), loader)
			buf := new(bytes.Buffer)
			buf.ReadFrom(composed)

			Expect(buf.String()).To(Equal("<div>Hello, client from Hello, client from b</div>"))
		})
	})
})
