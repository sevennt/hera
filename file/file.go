package file

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/fsnotify/fsnotify"
	"github.com/go-yaml/yaml"
)

// Provider file provider.
type Provider struct {
	path  string
	watch bool
}

// NewProvider returns new FileProvider.
func NewProvider(path string, watch bool) *Provider {
	return &Provider{path: path, watch: watch}
}

func (fp *Provider) Read() (conf map[string]interface{}, err error) {
	var content []byte
	content, err = ioutil.ReadFile(fp.path)
	if err != nil {
		return
	}

	fType := filepath.Ext(fp.path)
	conf = make(map[string]interface{})
	switch fType {
	case ".yaml", ".yml":
		err = yaml.Unmarshal(content, &conf)
	case ".json":
		err = json.Unmarshal(content, &conf)
	case ".toml":
		err = toml.Unmarshal(content, &conf)
	case ".ini":
		// TODO add ini configuration parse
		fallthrough
	default:
		err = fmt.Errorf("file type %v unsupported", fType)
	}
	return
}

// Watch file and automate update.
func (fp *Provider) Watch(watcher func(map[string]interface{})) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	defer w.Close()
	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-w.Events:
				log.Println("event: ", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file: ", event.Name)
				}
			case err := <-w.Errors:
				log.Println("error: ", err)
			}
		}
	}()

	err = w.Add(fp.path)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
