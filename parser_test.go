package composer_test

import (
	"strings"

	. "github.com/bernos/composer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Parser", func() {
	Context("When Tokenizing HTML", func() {
		var (
			parser *Parser
		)

		html := "<p><h1>Hello</h1><div composer-url=\"http://example.com\">inner content</div><h2>blah</h2><p composer-url=\"http://example.com/b\"></p>"

		BeforeEach(func() {
			parser = NewParser(strings.NewReader(html))
		})

		It("Should extract composer tags and buffer containing html", func() {
			buf, tags := parser.ExtractComposerTags()

			checkTag := func(t *ComposerTag, url string, offset int) {
				Expect(t.Url).To(Equal(url))
				Expect(t.Offset).To(Equal(offset))
			}

			Expect(buf.String()).To(Equal("<p><h1>Hello</h1><h2>blah</h2>"))
			Expect(len(tags)).To(Equal(2))

			checkTag(tags[0], "http://example.com", 17)
			checkTag(tags[1], "http://example.com/b", 30)
		})
	})
})
