package composer

import (
	"io"
	"sync"
)

type Loader func(string) io.Reader

type PipelineStep func(<-chan *ComposerTag) <-chan *ComposerTag

func BuildTagPipeline(tags []*ComposerTag, loader Loader) <-chan *ComposerTag {
	source := generateComposerTagsChannel(tags)
	load := fetchContentForComposerTags(loader)

	return fanIn(fanOut(source, load, len(tags)))
}

func fanOut(c <-chan *ComposerTag, f PipelineStep, n int) []<-chan *ComposerTag {
	workers := make([]<-chan *ComposerTag, 0)

	for i := 0; i < n; i++ {
		workers = append(workers, f(c))
	}

	return workers
}

func fanIn(cs []<-chan *ComposerTag) <-chan *ComposerTag {
	var wg sync.WaitGroup
	out := make(chan *ComposerTag)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan *ComposerTag) {
		for tag := range c {
			out <- tag
		}
		wg.Done()
	}

	wg.Add(len(cs))

	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func fetchContentForComposerTags(loader Loader) PipelineStep {
	return func(in <-chan *ComposerTag) <-chan *ComposerTag {
		out := make(chan *ComposerTag)
		go func() {
			for tag := range in {
				tag.Content = loader(tag.Url)
				out <- tag
			}
			close(out)
		}()
		return out
	}
}

func generateComposerTagsChannel(tags []*ComposerTag) <-chan *ComposerTag {
	out := make(chan *ComposerTag)
	go func() {
		for _, tag := range tags {
			out <- tag
		}
		close(out)
	}()
	return out
}
