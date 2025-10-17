package main

import (
	"context"
	"fmt"
	"progect_golang/meteocentr/sensorhumidity"
	sensorprssure "progect_golang/meteocentr/sensorpressure"
	"progect_golang/meteocentr/sensorseismic"
	"sync"

	"time"
)

var mtx sync.Mutex

func main() {

	humidityContext, humidityCancel := context.WithCancel(context.Background())
	pressureContext, pressureCancel := context.WithCancel(context.Background())
	seismicContext, seismicCancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(3 * time.Second)
		pressureCancel()
	}()

	go func() {
		time.Sleep(6 * time.Second)
		humidityCancel()
	}()

	go func() {
		time.Sleep(9 * time.Second)
		seismicCancel()
	}()

	pressureTransferPoint := sensorprssure.PressurePool(pressureContext, 2)
	humidityTransferPoint := sensorhumidity.HumidityPool(humidityContext, 2)
	seismicTransferPoint := sensorseismic.SeismicPool(seismicContext, 2)

	var humidity []int
	var pressure []int
	var seismics []int

	for v := range pressureTransferPoint {
		mtx.Lock()
		pressure = append(pressure, v)
		mtx.Unlock()
	}

	for v := range humidityTransferPoint {
		mtx.Lock()
		humidity = append(humidity, v)
		mtx.Unlock()
	}

	for v := range seismicTransferPoint {
		mtx.Lock()
		seismics = append(seismics, v)
		mtx.Unlock()
	}

	mtx.Lock()
	fmt.Println("датчик давления воздуха:", pressure)
	fmt.Println("датчик влажности воздуха:", humidity)
	fmt.Println("датчик сесмической активности:", seismics)
	mtx.Unlock()

}
