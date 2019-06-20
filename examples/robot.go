package main

import (
	"fmt"
	"robot"
	"time"
)

func main() {
	r, err := robot.NewRobot([]string{}, []string{"/Users/zy/.kube/config37"})
	if err != nil {
		panic(err)
	}
	r.Discover([]robot.Resource{robot.Services, robot.Endpoints}, []string{"default/details"})
	go r.Run()

	for {
		obj, _ := r.Pop()
		if err := process(obj); err != nil {
			r.ReQueue(obj)
		} else {
			r.Finish(obj)
		}
	}
}

func process(obj robot.QueueObject) error {
	// your own logic
	fmt.Println(time.Now(), obj.Event, obj.RType, obj.Key)
	return nil
}
