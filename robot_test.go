package robot

import (
	"fmt"
	"testing"
	"time"

	"github.com/k0kubun/pp"
)

func TestHello(t *testing.T) {
	r, err := NewRobot([]string{}, []string{"./config/pioneer"})
	if err != nil {
		panic(err)
	}
	r.Discover([]Resource{Services, Endpoints}, []string{"default/productpage"})
	go r.Run()

	for {
		obj, _ := r.Pop()
		fmt.Println(time.Now(), obj.Event, obj.Resource, obj.Key)
		if obj.Key == "default/productpage" {
			pp.Println(r.GetByKey(Endpoints, obj.Key))
		}

		//fmt.Println(r.ListKeys(), len(r.ListKeys(Services)))
	}
}
