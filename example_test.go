package batcher_test

import (
	"context"
	"fmt"

	"github.com/silviolleite/batcher"
)

type myData struct {
	name string
	age  int
}

func ExampleNew() {
	ctx := context.Background()
	var loggedHandlers []string

	b := batcher.New(&batcher.Options{
		Workers: 3,
		Logger: batcher.LoggerFunc(func(args ...interface{}) {
			loggedHandlers = append(loggedHandlers, args[0].(string))
		}),
	})

	fn := func(ctx context.Context, data []any) {
		for _, v := range data {
			_, ok := v.(myData)
			if !ok {
				fmt.Printf("data is not an myData")
				return
			}
		}
	}

	b.Start(ctx, fn)

	for i := 0; i < 1; i++ {
		_ = b.Add(myData{
			name: fmt.Sprintf("xtpo-%d", i),
			age:  i,
		})
	}
	b.Close()

	fmt.Println(len(loggedHandlers))
	// Output:
	// 7
}
