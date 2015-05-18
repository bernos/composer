package composer

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Loader func(string) io.Reader

func NewLoader(cache Cache) Loader {
	return func(url string) io.Reader {
		ch := make(chan io.Reader, 1)

		go func(url string) {

			cached, ok := cache.Get(url)

			if ok {
				ch <- strings.NewReader(cached)
			} else {

				res, err := http.Get(url)

				if err != nil {
					ch <- strings.NewReader(fmt.Sprintf("<-- composer encountered an error while loading content from %s - %v", url, err))
				} else {
					buf := new(bytes.Buffer)
					buf.ReadFrom(res.Body)
					content := buf.String()

					res.Body.Close()

					cache.Set(url, content, time.Now().Add(time.Minute*10))
					ch <- strings.NewReader(content)
				}
			}
		}(url)

		select {
		case r := <-ch:
			return r
		case <-time.After(time.Second * 1):
			return nil
		}
	}
}
