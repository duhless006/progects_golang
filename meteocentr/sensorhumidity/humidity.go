package sensorhumidity

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func humidity(ctx context.Context,
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
			fmt.Println("датчик влажности воздуха номер:", n, "---->>>закончил вычисление")
			return
		default:
			fmt.Println("датчик влажности воздуха номер:", n, "начал вычисление")
			time.Sleep(1 * time.Second)
			info := r.Intn(800)
			fmt.Println("датчик влажности воздуха номер:", n, "вычислил показания", info)

			dataPoint <- info
			fmt.Println("датчик влажности воздуха номер:", n, "----показания переданы----")
			fmt.Println("")
		}
	}
}

func HumidityPool(ctx context.Context, humidityCount int) <-chan int {
	humidityTransferPoint := make(chan int)
	wg := &sync.WaitGroup{}

	for i := 1; i <= humidityCount; i++ {
		wg.Add(1)
		go humidity(ctx, wg, humidityTransferPoint, i)
	}
	go func() {
		wg.Wait()
		close(humidityTransferPoint)
	}()
	return humidityTransferPoint
}
