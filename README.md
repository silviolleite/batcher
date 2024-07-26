# Batcher

![coverage](https://raw.githubusercontent.com/silviolleite/batcher/badges/.badges/main/coverage.svg)

---

This package provides a worker pool to batch items. Create a Batcher, start it with a handler function and add items to it. Items are accumulated by BatchSize or Timeout, when on of these threshold was satisfied the handler function is called with items accumulated. You can repeat this processes for many times.

### Install

---

Manual install:

```bash
go get -u github.com/silviolleite/batcher
```

Golang import:

```go
import "github.com/silviolleite/batcher"
```


### Usage

---

```go 
package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/silviolleite/batcher"
)

type myData struct {
	name string
	age  int
}

var processed int

func batchHandler(ctx context.Context, data []any) {
	var totalCount int
	var mu sync.Mutex
	for _, v := range data {
		d, ok := v.(myData)
		if !ok {
			fmt.Printf("data is not an myData")
			return
		}
		fmt.Printf("%+v\n", d)
		mu.Lock()
		totalCount += 1
		mu.Unlock()
		time.Sleep(300 * time.Millisecond)
	}
	processed += totalCount
}

func main() {
	maxCount := 30
	b := batcher.New(&batcher.Options{Workers: 3, BatchSize: 4})
	b.Start(context.Background(), batchHandler)
	for i := 0; i < maxCount; i++ {
		err := b.Add(myData{
			name: "xpto",
			age:  i,
		})
		if err != nil {
			panic(err)
		}
	}
	b.Close()
	fmt.Printf("processed: %d\n", processed)
}
```
