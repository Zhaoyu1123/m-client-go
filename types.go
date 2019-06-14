package robot

type Resource int

const (
	All 		Resource = iota

	Services

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
type event int

const (
	// EventAdd is sent when an object is added
	EventAdd event = iota

	// EventUpdate is sent when an object is modified
	// Captures the modified object
	EventUpdate

	// EventDelete is sent when an object is deleted
	// Captures the object at the last known state
	EventDelete
)

func (e event) String() string {
	out := "unknown"
	switch e {
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
	Event   event
	Key   	string
}
