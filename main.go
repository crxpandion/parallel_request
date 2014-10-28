package main

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/crxpandion/parallel_request/para"
)

type Worker struct {
	id   int
	wait time.Duration
	e    error
}

func (w *Worker) Work() (string, error) {
	time.Sleep(w.wait)
	if w.e != nil {
		return "FAILED " + strconv.Itoa(w.id), w.e
	}
	return "SUCCESS " + strconv.Itoa(w.id), nil
}

func main() {
	runtime.GOMAXPROCS(1)
	f := para.NewFetcher()
	start := time.Now()
	for _, w := range []Worker{
		Worker{
			id:   1,
			wait: time.Second * 10,
		},
		Worker{
			id:   2,
			wait: time.Second * 5,
			e:    errors.New("error 1"),
		},
		Worker{
			id:   3,
			wait: time.Second * 15,
		},
		Worker{
			id:   4,
			wait: time.Second * 10,
			e:    errors.New("error 2"),
		},
		Worker{
			id:   5,
			wait: time.Second * 20,
		},
	} {
		ww := w
		f.Fetch(ww.Work)
	}
	finalResult := <-f.Get
	fmt.Printf("Took %s\n", time.Now().Sub(start))
	fmt.Printf("%+v\n", finalResult)

	<-f.Done
	fmt.Printf("Cleanup took %s\n", time.Now().Sub(start))
}
