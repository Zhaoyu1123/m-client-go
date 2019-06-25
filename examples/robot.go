package main

import (
	"fmt"
	"gitlab.mfwdev.com/servicemesh/robot"
	"time"
)

func main() {
	r, err := robot.NewRobot(
		robot.Cluster{
			ConfigPath:"/Users/zy/.kube/config37",
			Resources:[]robot.RN{
				{robot.Services, "istio-system"},
				{robot.Pods, "istio-system"},
				{robot.Endpoints, "default"},
			},
		},
		robot.Cluster{
			ConfigPath:"/Users/zy/.kube/config39",
			Resources:[]robot.RN{
				{robot.Services, "istio-system"},
				{robot.Pods, "istio-system"},
				{robot.Pods, "default"},
			},
		},
	)
	if err != nil {
		panic(err)
	}
	go r.Run()

	for {
		obj, _ := r.Pop()
		if err := process(obj); err != nil {
			// r.GetByKey()
			// r.ListKeys()
			r.ReQueue(obj)
		} else {
			r.Finish(obj)
		}

	}

	// go r.Stop()
}

func process(obj robot.QueueObject) error {
	// your own logic
	fmt.Println(time.Now(), obj.Event, obj.RType, obj.Key)
	return nil
}
