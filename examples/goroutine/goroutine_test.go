package goroutine

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"sync"
	"testing"
)

func TestGoroutine(t *testing.T) {
	words := []string{"Go,PHP"}
	strings, err := coSearchWithContext(context.Background(), words)
	if err != nil {
		return
	}
	fmt.Println(strings)

}
func search(ctx context.Context, word string) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		if word == "Go" || word == "Java" {
			return "", errors.New("Go or Java")
		}
		return fmt.Sprintf("result: %s", word), nil // 模拟结果
	}
}

func coSearchWithContext(ctx context.Context, words []string) ([]string, error) {
	// 创建goroutine的管理器
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var (
		// 并发等待
		wg = sync.WaitGroup{}
		// 只执行一次
		once = sync.Once{}

		results = make([]string, len(words))
		// 控制goroutine数量
		tokens = make(chan struct{}, 2)

		err error
	)

	for i, word := range words {
		tokens <- struct{}{}
		wg.Add(1)

		go func(word string, i int) {
			defer func() {
				wg.Done()
				<-tokens
			}()

			result, e := search(ctx, word)
			if e != nil {
				once.Do(func() {
					err = e
					cancel()
				})

				return
			}

			results[i] = result
		}(word, i)
	}

	wg.Wait()

	return results, err
}

func TestErrGroup(t *testing.T) {
	group, err := coSearchWithErrGroup(context.Background(), []string{"java", "js", "html"})
	if err != nil {
		return
	}
	fmt.Println(group)
}
func coSearchWithErrGroup(ctx context.Context, words []string) ([]string, error) {
	// 上下文
	g, ctx := errgroup.WithContext(ctx)
	// 并发goroutine控制
	g.SetLimit(10)

	results := make([]string, len(words))

	for i, word := range words {
		i, word := i, word
		// 启动goroutine
		g.Go(func() error {
			result, err := search(ctx, word)
			if err != nil {
				return err
			}

			results[i] = result
			return nil
		})
	}
	// wait
	err := g.Wait()

	return results, err
}
