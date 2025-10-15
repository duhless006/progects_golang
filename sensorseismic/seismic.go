package sensorseismic

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func seismics(ctx context.Context,
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
			fmt.Println("датчик сейсмической активности номер:", n, "---->>>закончил вычисление")
			return
		default:
			fmt.Println("датчик сейсмической активности номер:", n, "начинает вычисление")
			time.Sleep(1 * time.Second)
			info := r.Intn(800)
			fmt.Println("датчик сейсмической активности номер:", n, "вычислил показания", info)

			dataPoint <- info
			fmt.Println("датчик сейсмической активности номер:", n, "----показания переданы----")
			fmt.Println("")
		}
	}
}

func SeismicPool(ctx context.Context, seismicCount int) <-chan int {
	seismicTransferPoint := make(chan int)
	wg := &sync.WaitGroup{}

	for i := 1; i <= seismicCount; i++ {
		wg.Add(1)
		go seismics(ctx, wg, seismicTransferPoint, i)
	}
	go func() {
		wg.Wait()
		close(seismicTransferPoint)
	}()
	return seismicTransferPoint
}
