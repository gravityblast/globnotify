# globnotify
golang glob watcher. it watches a filesystem tree and notifies on file changes based on a glob pattern.

```go
package main

import (
	"fmt"
	"log"

	"github.com/pilu/globnotify"
)

func main() {
	w, err := globnotify.New("./**/*.css")
	if err != nil {
		log.Fatal(err)
	}

	events, err := w.Watch()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("watching...\n")

	for {
		select {
		case event := <-events:
			fmt.Printf("%+v\n", event)
		}
	}
}
```
