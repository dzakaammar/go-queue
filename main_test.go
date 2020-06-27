package main

import (
	"sync"
	"testing"
	"time"
)

func TestQueue(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		data := []int{1, 2, 3}
		stopCh := make(chan struct{})
		defer close(stopCh)

		dataCh := queue(data, stopCh)

		first := <-dataCh
		if first != data[0] {
			t.Fatalf("invalid result. expecting %d, but got %d", data[0], first)
		}

		second := <-dataCh
		if second != data[1] {
			t.Fatalf("invalid result. expecting %d, but got %d", data[1], second)
		}

		third := <-dataCh
		if third != data[2] {
			t.Fatalf("invalid result. expecting %d, but got %d", data[2], third)
		}
	})

	t.Run("closing the stop channel", func(t *testing.T) {
		data := []int{1, 2, 3}
		stopCh := make(chan struct{})
		close(stopCh)

		dataCh := queue(data, stopCh)

		d, valid := <-dataCh
		if valid || d != 0 {
			t.Fatal("should be invalid")
		}
	})
}

func TestWorker(t *testing.T) {
	data := []int{1, 2}
	stopCh := make(chan struct{})
	defer close(stopCh)

	dataCh := queue(data, stopCh)
	start := time.Now()

	wg := new(sync.WaitGroup)

	wg.Add(1)
	worker(wg, dataCh)
	wg.Wait()

	finish := -start.Sub(time.Now()).Round(time.Second)
	if finish != time.Duration(3*time.Second) {
		t.Fatalf("invalid time. expecting 3 second, got %+v", finish)
	}

}
