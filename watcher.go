package globnotify

import (
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	zglob "github.com/mattn/go-zglob"
	"github.com/pilu/treenotify"
)

type Watcher interface {
	Watch() (chan fsnotify.Event, error)
	Close()
}

type watcher struct {
	tw     treenotify.Watcher
	glob   string
	events chan fsnotify.Event
	stop   chan struct{}
}

func New(glob string) (Watcher, error) {
	tw, err := treenotify.New()
	if err != nil {
		return nil, err
	}

	absGlob, err := filepath.Abs(glob)
	if err != nil {
		return nil, err
	}

	return &watcher{
		tw:     tw,
		glob:   absGlob,
		events: make(chan fsnotify.Event),
		stop:   make(chan struct{}),
	}, nil
}

func root(glob string) string {
	parts := strings.Split(glob, "*")
	return parts[0]
}

func (w *watcher) Watch() (chan fsnotify.Event, error) {
	root := root(w.glob)
	events, err := w.tw.Watch(root)
	if err != nil {
		return nil, err
	}

	go w.watch(events)

	return w.events, nil
}

func (w *watcher) watch(events chan fsnotify.Event) {
	working := true
	for working {
		select {
		case <-w.stop:
			working = false
		case event := <-events:
			match, _ := zglob.Match(w.glob, event.Name)
			if match {
				w.events <- event
			}
		}
	}
}

func (w *watcher) Close() {
	w.stop <- struct{}{}
	w.tw.Close()
}
