package batcher_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/silviolleite/batcher"
)

func TestBatcher_Start(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		b := batcher.New(nil)
		ctx := context.Background()
		var count int
		want := 42

		fn := func(ctx context.Context, data []any) {
			for _, v := range data {
				_, ok := v.(myData)
				if !ok {
					fmt.Printf("data is not an myData")
					return
				}
			}
			count += len(data)
		}

		b.Start(ctx, fn)

		for i := 0; i < want; i++ {
			_ = b.Add(myData{
				name: fmt.Sprintf("xtpo-%d", i),
				age:  i,
			})
		}
		b.Close()
		assert.Equal(t, want, count)
	})

	t.Run("Success with timeout", func(t *testing.T) {
		b := batcher.New(&batcher.Options{
			BatchSize:    10,
			BatchTimeout: 10 * time.Millisecond,
		})
		ctx := context.Background()
		var count int
		want := 1

		fn := func(ctx context.Context, data []any) {
			for _, v := range data {
				_, ok := v.(myData)
				if !ok {
					fmt.Printf("data is not an myData")
					return
				}
			}
			count += len(data)
		}

		b.Start(ctx, fn)

		err := b.Add(myData{
			name: fmt.Sprintf("xtpo"),
			age:  2,
		})
		assert.NoError(t, err)
		time.Sleep(50 * time.Millisecond)
		b.Close()
		assert.Equal(t, want, count)
	})
}

func TestBatcher_Add(t *testing.T) {
	t.Run("should return error when the channel is closed", func(t *testing.T) {
		b := batcher.New(nil)
		ctx := context.Background()

		fn := func(ctx context.Context, data []any) {
			for _, v := range data {
				_, ok := v.(myData)
				if !ok {
					t.Error("Expecting an myData got")
					return
				}
			}
		}

		b.Start(ctx, fn)

		err := b.Add(myData{
			name: fmt.Sprintf("xtpo"),
			age:  1,
		})
		assert.NoError(t, err)

		err = b.Add(myData{
			name: fmt.Sprintf("xtpo2"),
			age:  2,
		})
		assert.NoError(t, err)

		b.Close()
		err = b.Add(myData{
			name: fmt.Sprintf("xtpo3"),
			age:  3,
		})
		assert.Error(t, err)
		assert.ErrorIs(t, batcher.ErrInputChannelClosed, err)
	})

	t.Run("should return error when the input is nil", func(t *testing.T) {
		b := batcher.New(nil)
		ctx := context.Background()

		fn := func(ctx context.Context, data []any) {
			for _, v := range data {
				_, ok := v.(myData)
				if !ok {
					t.Error("Expecting an myData got")
					return
				}
			}
		}

		b.Start(ctx, fn)

		err := b.Add(myData{
			name: fmt.Sprintf("xtpo"),
			age:  1,
		})
		assert.NoError(t, err)

		err = b.Add(nil)
		assert.Error(t, err)
		assert.ErrorIs(t, batcher.ErrInputIsNil, err)

		b.Close()
	})
}
