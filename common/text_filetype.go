package common

import (
	"engo.io/engo"
	"io"
	"io/ioutil"
	"fmt"
)

// TextResource contains a level created from a Tile Map XML
type TextResource struct {
	// Text holds the text read from the file
	Text string
	url   string
}

// URL retrieves the url to the .text file
func (r TextResource) URL() string {
	return r.url
}


// textLoader is responsible for managing readable text files within 'engo.Files'.
type textLoader struct {
	files map[string]TextResource
}

// Load will load the text file and any other image resources that are needed
func (t *textLoader) Load(url string, data io.Reader) error {
	textBytes, err := ioutil.ReadAll(data)
	if err != nil {
		return err
	}

	t.files[url] = TextResource{
		Text: string(textBytes),
		url: url,
	}

	return nil
}

// Unload removes the preloaded level from the cache
func (t *textLoader) Unload(url string) error {
	delete(t.files, url)
	return nil
}

// Resource retrieves and returns the preloaded level of type 'TextResource'
func (t *textLoader) Resource(url string) (engo.Resource, error) {
	text, ok := t.files[url]
	if !ok {
		return nil, fmt.Errorf("resource not loaded by `FileLoader`: %q", url)
	}

	return text, nil
}

func init() {
	engo.Files.Register(".txt", &textLoader{files: make(map[string]TextResource)})
	engo.Files.Register(".json", &textLoader{files: make(map[string]TextResource)})
	engo.Files.Register(".text", &textLoader{files: make(map[string]TextResource)})
}
