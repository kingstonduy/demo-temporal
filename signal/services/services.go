package services

import (
	"fmt"
	"time"
)

func PrintService[K any](input K) {
	time.Sleep(time.Second * 5)
	fmt.Println("💡 Blocking activity countdown:", input)

	return
}

func InputService[K any](input K) {
	time.Sleep(time.Second * 10)
	fmt.Println("💡 Input activity:", input)

	return
}
