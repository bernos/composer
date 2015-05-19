package composer

import (
	"bytes"
	"io"

	"golang.org/x/net/html"
)

type Parser struct {
	z *html.Tokenizer
}

func NewParser(r io.Reader) *Parser {
	return &Parser{z: html.NewTokenizer(r)}
}

/**
 * extractComposerTags will read html content from an io.Reader and extract any
 * composer tags, creating a ComposerTag for each. It will also append all html
 * content outside of the composer tags to to a buffer, ready to merged with
 * the tags
 */
func (p *Parser) ExtractComposerTags() (*bytes.Buffer, []*ComposerTag) {
	buf := new(bytes.Buffer)
	depth := 0
	tags := make([]*ComposerTag, 0)

	for {
		tt := p.z.Next()

		if tt == html.ErrorToken {
			return buf, tags
		}

		if depth == 0 {
			if tag := p.GetComposerTag(); tag != nil {
				tag.Offset = buf.Len()
				tags = append(tags, tag)
				if tt == html.StartTagToken {
					depth++
				}
			} else {
				buf.Write(p.z.Raw())
			}
		} else if tt == html.EndTagToken {
			depth--
		}
	}
}

func (p *Parser) GetComposerTag() *ComposerTag {

	attributes := p.GetAttributes(make([]*html.Attribute, 0))

	composerTag := new(ComposerTag)

	for i := range attributes {
		if attributes[i].Key == "composer-url" {
			composerTag.Url = attributes[i].Val
		}
	}

	if len(composerTag.Url) > 0 {
		return composerTag
	}

	return nil
}

func (p *Parser) GetAttributes(a []*html.Attribute) []*html.Attribute {
	key, val, more := p.z.TagAttr()

	a = append(a, &html.Attribute{Key: string(key), Val: string(val)})

	if more {
		return p.GetAttributes(a)
	}

	return a
}
