package robot

type Resource int

const (
	Services  Resource = iota

	Endpoints

	Pods

	ConfigMaps
)

func (t Resource) String() string {
	out := "unknown"
	switch t {
	case Services:
		out = "services"
	case Endpoints:
		out = "endpoints"
	case Pods:
		out = "pods"
	case ConfigMaps:
		out = "configmaps"
	}
	return out
}

// Event represents a registry update event
type Event int

const (
	// EventAdd is sent when an object is added
	EventAdd Event = iota

	// EventUpdate is sent when an object is modified
	// Captures the modified object
	EventUpdate

	// EventDelete is sent when an object is deleted
	// Captures the object at the last known state
	EventDelete
)

func (event Event) String() string {
	out := "unknown"
	switch event {
	case EventAdd:
		out = "add"
	case EventUpdate:
		out = "update"
	case EventDelete:
		out = "delete"
	}
	return out
}

type QueueObject struct {
	typ   Event
	key   string
}
