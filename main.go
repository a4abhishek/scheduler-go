package main

import (
	"fmt"
	"time"
)

type Meta struct {
	driverID string
	message  string
	sig      <-chan time.Time
}

func main() {
	chanList := []Meta{}

	chanList = append(chanList, Meta{
		driverID: "123",
		message:  "msg1",
		sig:      time.NewTimer(5 * time.Second).C})

	chanList = append(chanList, Meta{
		driverID: "345",
		message:  "msg2",
		sig:      time.NewTimer(1 * time.Second).C})

	fanIn := make(chan Meta)

	go func(sig <-chan Meta) {
		for x := range sig {
			fmt.Println("driver-id: ", x.driverID, " message: ", x.message)
		}
	}(fanIn)

	i := 0
	for {
		if len(chanList) != 0 {
			select {
			case <-chanList[i].sig:
				fanIn <- chanList[i]
				chanList = append(chanList[:i], chanList[i+1:]...)
			default:
				time.Sleep(time.Second)
			}
		} else {
			//time.Sleep(time.Minute) // TODO: Uncomment it if you want to run it always
			break // TODO: Comment it if you want to run it always
		}

		i++
		if i >= len(chanList) {
			i = 0
		}
	}
}
