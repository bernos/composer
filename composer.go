package composer

import (
	"bytes"
	"io"
)

type ComposerTag struct {
	Content io.Reader
	Offset  int
	Url     string
}

func Compose(r io.Reader, l Loader) io.Reader {
	parser := NewParser(r)
	buf, tags := parser.ExtractComposerTags()
	output := new(bytes.Buffer)

	if len(tags) > 0 {
		offset := 0

		for tag := range BuildTagPipeline(tags, l) {
			output.Write(buf.Next(tag.Offset - offset))

			if tag.Content != nil {
				output.ReadFrom(Compose(tag.Content, l))
			} else {
				output.WriteString("<!-- composer failed to load content from " + tag.Url + " -->")
			}

			offset += tag.Offset
		}
	}

	output.Write(buf.Next(buf.Len()))

	return output
}
