package examplecode

import (
	"sync"

	"github.com/Sirupsen/logrus"
)

func bob() {
	logrus.Infof("Hello World from Bob")
}

func main() {
	logrus.Infof("Hello World")
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		bob()
		wg.Done()
	}()
	wg.Wait()
}
