package main

import (
	"context"
	"fmt"
	"time"
)

/*
Идея заключается в том, чтобы создать канал с ограниченным количеством выходов.

	Допустим, мы читаем сообщения из потока данных, и нас интересуют только первые 5 сообщений.
	 Также давайте предположим, что мы не контролируем сам источник данных, мы можем только читать из него
*/
func main() {

	integI := make(chan interface{})

	go func() {
		defer close(integI)
		for i := 0; i < 1000; i++ {
			integI <- i
		}
	}()

	appCtx, cancel := context.WithCancel(context.Background())

	res := takeNData(appCtx, 100, integI)
	go func(res <-chan interface{}) {
		for r := range res {
			fmt.Println(r)
		}
	}(res)
	time.Sleep(100 * time.Microsecond)
	cancel()
}
func takeNData(ctx context.Context, n int, dataChannel <-chan interface{}) <-chan interface{} {
	result := make(chan interface{})

	go func() {
		defer close(result)

		// 3
		for i := 0; i < n; i++ {
			select {
			case val, ok := <-dataChannel:
				if !ok {
					return
				}
				result <- val
			case <-ctx.Done():
				return
			}
		}
	}()
	return result
}
