package fs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"

	"github.com/taxibeat/hypatia/scrape"

	b64 "encoding/base64"

	"github.com/beatlabs/patron/log"
)

const (
	syncFile  = "swagger.json"
	asyncFile = "async.json"
)

type Scraper struct {
	path    string
	watcher *fsnotify.Watcher
	ntf     chan bool
}

func New(path string) (*Scraper, error) {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return NewChild(path, watcher)
}

func NewChild(path string, watcher *fsnotify.Watcher) (*Scraper, error) {
	fl, err := getFile(path)
	if err != nil {
		return nil, err

	}

	err = fl.Sync()
	if err != nil {
		return nil, err

	}

	err = fl.Close()
	if err != nil {
		log.Infof("CLosing:", err)
	}

	scr := &Scraper{path: path, watcher: watcher, ntf: make(chan bool)}

	return scr, nil
}

func (s *Scraper) Updates() chan bool {
	return s.ntf
}

func (s *Scraper) Scrape() []scrape.DocDef {
	s.newWatcher()
	return s.ScrapeFolder()
}

func (s *Scraper) ScrapeFolder() []scrape.DocDef {

	result := []scrape.DocDef{}

	fl, err := getFile(s.path)
	if err != nil {
		log.Debugf("cannot create the scraper", err)
		fl.Close()
		return result
	}
	fls, err := fl.Readdir(0)
	if err != nil {
		log.Debugf("cannot read the directory", err)
		fl.Close()
		return result
	}

	err = fl.Close()
	if err != nil {
		log.Infof("cannot close the directory scrapeFolder", err, fl, s.path)
		return result

	}

	var foldersPool []os.FileInfo

	for _, f := range fls {
		if f.IsDir() {
			switch f.Name() {
			case "doc", "docs":
				scr, err := NewChild(fl.Name()+"/"+f.Name(), s.watcher)
				if err != nil {
					log.Debug(err)
					continue
				}

				dcs, err := scr.scrapeDocs()
				if err != nil {
					log.Debug(err)
					continue
				}
				result = append(result, dcs...)
				return result
			case ".", "..", ".git", "vendor":
				// ignore
			default:
				foldersPool = append(foldersPool, f)
			}

		}
	}

	for _, fld := range foldersPool {
		path, err := filepath.Abs(fl.Name())
		if err != nil {
			log.Infof("Invalid path", err)
		}
		scr, err := NewChild(filepath.Join(path, fld.Name()), s.watcher)
		if err != nil {

			log.Fatalf("WHat?", fl.Name()+"|||||||||||"+fld.Name(), err)
		}

		rss := scr.ScrapeFolder()
		result = append(result, rss...)
	}

	return result
}

func (s *Scraper) scrapeDocs() ([]scrape.DocDef, error) {
	var res []scrape.DocDef

	fld, err := getFile(s.path)
	if err != nil {
		return res, err
	}

	fls, err := fld.Readdir(0)
	if err != nil {
		return nil, err
	}

	err = fld.Sync()
	if err != nil {
		log.Fatalf("cannot read the directory", err)

	}

	err = fld.Close()
	if err != nil {
		log.Fatalf("cannot read the directory", err)

	}

	for _, f := range fls {
		dd, err := s.retrieveDocumentation(f)
		if err != nil {
			///fmt.Println(err)
		} else {
			res = append(res, *dd)
		}
	}
	return res, nil
}

func (s *Scraper) retrieveDocumentation(doc os.FileInfo) (*scrape.DocDef, error) {
	result := scrape.DocDef{}

	fld, err := getFile(s.path)
	if err != nil {
		return &result, err
	}
	err = fld.Sync()
	if err != nil {
		log.Fatalf("cannot sync the directory", err)

	}

	err = fld.Close()
	if err != nil {
		log.Fatalf("cannot close the directory", err)

	}

	switch doc.Name() {
	case syncFile:
		result.Type = scrape.Swagger
	case asyncFile:
		result.Type = scrape.Async
	default:
		return nil, fmt.Errorf("unsupported type: %s", doc.Name())
	}

	name := fld.Name() + "/" + doc.Name()

	result.RepoName = name

	rpl := strings.NewReplacer(".", "d", string(os.PathSeparator), "-")

	result.ID = b64.StdEncoding.EncodeToString([]byte(rpl.Replace(fmt.Sprintf("%s-%s", fld.Name(), result.Type))))
	definition, err := ioutil.ReadFile(fld.Name() + "/" + doc.Name())
	if err != nil {
		return nil, err
	}
	result.Definition = string(definition)

	err = s.watcher.Add(name)
	if err != nil {
		fmt.Println("Error watching", err)
	}
	return &result, nil
}

func getFile(path string) (*os.File, error) {
	fl, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	st, err := fl.Stat()
	if err != nil {
		return nil, err
	}
	if !st.Mode().IsDir() {
		return nil, errors.New("the path specified is not a directory")
	}

	return fl, nil
}

func (s *Scraper) newWatcher() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	s.watcher.Close()
	s.watcher = watcher

	go func() {
		for {
			select {
			// watch for events
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				triggers := map[fsnotify.Op]bool{fsnotify.Remove: true, fsnotify.Write: true, fsnotify.Create: true}
				if _, ok := triggers[event.Op]; ok {
					s.ntf <- true
				} else {
					log.Debugf("We received", event.String())
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Errorf("error event while watching files:", err)
			}
		}
	}()
	return err
}
