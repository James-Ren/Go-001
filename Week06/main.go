package main

import (
	"fmt"
	"time"

	"github.com/James-Ren/Go-001/tree/main/Week06/rolling"
)

func main() {
	opt := rolling.RollingCounterOpts{Size: 5, BucketDuration: 100 * time.Millisecond}
	rollingCounter := rolling.NewRollingCounter(opt)
	rollingCounter.Add(2)
	time.Sleep(300 * time.Millisecond)
	rollingCounter.Add(5)
	fmt.Printf("rolling counter avg:%f,Value:%d\n", rollingCounter.Avg(), rollingCounter.Value())
	time.Sleep(200 * time.Millisecond)
	rollingCounter.Add(1)
	rollingCounter.Add(3)
	fmt.Printf("rolling counter avg:%f,Value:%d\n", rollingCounter.Avg(), rollingCounter.Value())
}
