package composer_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	. "github.com/bernos/composer"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Loader", func() {
	Context("When loading content", func() {

		var (
			server *httptest.Server
			url    string
		)

		BeforeEach(func() {
			server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, "From Server")
			}))

			url = server.URL
		})

		AfterEach(func() {
			defer server.Close()
		})

		It("Should load content from the server", func() {
			loader := NewLoader(NewMemoryCache())
			content := loader(url)
			buf := new(bytes.Buffer)
			buf.ReadFrom(content)

			Expect(buf.String()).To(Equal("From Server"))
		})

		It("Should add loaded content to the cache", func() {
			cache := NewMemoryCache()
			loader := NewLoader(cache)

			_ = loader(url)

			cached, ok := cache.Get(url)

			Expect(ok).To(BeTrue())
			Expect(cached).To(Equal("From Server"))
		})

		It("Should use content from the cache when available", func() {
			cache := NewMemoryCache()
			loader := NewLoader(cache)

			cache.Set(url, "From Cache", time.Now().Add(time.Second*10))

			content := loader(url)
			buf := new(bytes.Buffer)
			buf.ReadFrom(content)

			Expect(buf.String()).To(Equal("From Cache"))
		})
	})
})
