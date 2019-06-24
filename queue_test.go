package robot

import (
	"testing"
)

func TestQueue(t *testing.T) {
	q := newWorkQueue()

	objOne := QueueObject{EventAdd, Endpoints, "one"}
	objTwo := QueueObject{EventAdd, Endpoints, "two"}

	q.push(objOne)

	q.push(objOne)

	q.push(objOne)

	q.push(objOne)

	q.push(objOne)

	if e, a := 5, q.NumRequeues(objOne); e != a {
		t.Errorf("expected %v, got %v", e, a)
	}

	obj, _ := q.Pop()
	if obj != objOne {
		t.Errorf("expected %v, got %v", objOne, obj)
	}
	q.Finish(objOne)

	q.push(objOne)
	obj, _ = q.Pop()
	if obj != objOne {
		t.Errorf("expected %v, got %v", objOne, obj)
	}

	q.push(objTwo)
	obj, _ = q.Pop()
	if obj != objTwo {
		t.Errorf("expected %v, got %v", objTwo, obj)
	}

	if e, a := 1, q.NumRequeues(objOne); e != a {
		t.Errorf("expected %v, got %v", e, a)
	}

	if e, a := 1, q.NumRequeues(objTwo); e != a {
		t.Errorf("expected %v, got %v", e, a)
	}

	_ = q.ReQueue(objOne)
	_ = q.ReQueue(objOne)
	if e, a := 3, q.NumRequeues(objOne); e != a {
		t.Errorf("expected %v, got %v", e, a)
	}

	_ = q.ReQueue(objOne)
	if e, a := 0, q.NumRequeues(objOne); e != a {
		t.Errorf("expected %v, got %v", e, a)
	}

	q.Forget(objOne)
	q.Forget(objOne)
	q.Forget(objTwo)
	q.Forget(objTwo)
	if e, a := 0, q.NumRequeues(objOne); e != a {
		t.Errorf("expected %v, got %v", e, a)
	}

	if e, a := 0, q.NumRequeues(objTwo); e != a {
		t.Errorf("expected %v, got %v", e, a)
	}
}
