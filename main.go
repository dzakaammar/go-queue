package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
)

var workerCount int

func main() {
	flag.IntVar(&workerCount, "worker_count", 0, "number of worker")
	flag.Parse()

	if workerCount <= 0 {
		log.Fatal("invalid number of worker")
	}

	stopCh := make(chan struct{})

	// listen for interrupt
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		<-sigCh
		log.Println("Ctrl + C. Interrupted...")
		close(stopCh)
	}()

	q := []int{1, 2, 4, 2, 3, 5, 2, 3, 1, 3}
	dataCh := queue(q, stopCh)

	wg := new(sync.WaitGroup)
	for i := 1; i <= workerCount; i++ {
		wg.Add(1)
		go worker(wg, dataCh)
	}

	wg.Wait()
}

func queue(queues []int, stopCh <-chan struct{}) <-chan int {
	dataCh := make(chan int, len(queues))
	go func() {
		defer close(dataCh)
		for _, q := range queues {
			select {
			case <-stopCh:
				log.Println("Closing the queue...")
				return
			default:
				dataCh <- q
			}
		}

	}()

	return dataCh
}

func worker(wg *sync.WaitGroup, dataCh <-chan int) {
	defer wg.Done()

	for data := range dataCh {
		time.Sleep(time.Duration(data) * time.Second)
		log.Printf("queue processed with %d second\n", data)
	}
}
