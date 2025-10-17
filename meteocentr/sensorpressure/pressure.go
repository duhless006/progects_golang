package sensorprssure

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func pressure(ctx context.Context,
	wg *sync.WaitGroup,
	dataPoint chan<- int,
	n int,
) {
	defer wg.Done()
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("датчик давления воздуха номер:", n, "--->>>закончид работу")
			return
		default:
			fmt.Println("датчик давления воздуха номер:", n, "включился")
			time.Sleep(1 * time.Second)
			info := r.Intn(800)
			fmt.Println("датчик давления воздуха номер:", n, "вычислил показания", info)

			dataPoint <- info
			fmt.Println("датчик давления воздуха номер:", n, "----передал показания---")
			fmt.Println("")
		}
	}
}

func PressurePool(ctx context.Context, pressureCount int) chan int {
	pressureTransferPoint := make(chan int)
	wg := &sync.WaitGroup{}

	for i := 1; i <= pressureCount; i++ {
		wg.Add(1)
		go pressure(ctx, wg, pressureTransferPoint, i)

	}
	go func() {
		wg.Wait()
		close(pressureTransferPoint)
	}()
	return pressureTransferPoint

}
