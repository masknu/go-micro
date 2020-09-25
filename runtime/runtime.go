// Package runtime is a service runtime manager
package runtime

import (
	"errors"
	"time"
)

var (
	ErrAlreadyExists   = errors.New("already exists")
	ErrInvalidResource = errors.New("invalid resource")
)

// Runtime is a service runtime manager
type Runtime interface {
	// Init initializes runtime
	Init(...Option) error
	// Create a resource
	Create(Resource, ...CreateOption) error
	// Read a resource
	Read(...ReadOption) ([]*Service, error)
	// Update a resource
	Update(Resource, ...UpdateOption) error
	// Delete a resource
	Delete(Resource, ...DeleteOption) error
	// Logs returns the logs for a resource
	Logs(Resource, ...LogsOption) (Logs, error)
	// Start starts the runtime
	Start() error
	// Stop shuts down the runtime
	Stop() error
	// String defines the runtime implementation
	String() string
}

// Logs returns a log stream
type Logs interface {
	Error() error
	Chan() chan Log
	Stop() error
}

// Log is a log message
type Log struct {
	Message  string
	Metadata map[string]string
}

// Scheduler is a runtime service scheduler
type Scheduler interface {
	// Notify publishes schedule events
	Notify() (<-chan Event, error)
	// Close stops the scheduler
	Close() error
}

// EventType defines schedule event
type EventType int

const (
	// Create is emitted when a new build has been craeted
	Create EventType = iota
	// Update is emitted when a new update become available
	Update
	// Delete is emitted when a build has been deleted
	Delete
)

// String returns human readable event type
func (t EventType) String() string {
	switch t {
	case Create:
		return "create"
	case Delete:
		return "delete"
	case Update:
		return "update"
	default:
		return "unknown"
	}
}

// Event is notification event
type Event struct {
	// ID of the event
	ID string
	// Type is event type
	Type EventType
	// Timestamp is event timestamp
	Timestamp time.Time
	// Service the event relates to
	Service *Service
	// Options to use when processing the event
	Options *CreateOptions
}

// ServiceStatus defines service statuses
type ServiceStatus int

const (
	// Unknown indicates the status of the service is not known
	Unknown ServiceStatus = iota
	// Pending is the initial status of a service
	Pending
	// Building is the status when the service is being built
	Building
	// Starting is the status when the service has been started but is not yet ready to accept traffic
	Starting
	// Running is the status when the service is active and accepting traffic
	Running
	// Stopping is the status when a service is stopping
	Stopping
	// Stopped is the status when a service has been stopped or has completed
	Stopped
	// Error is the status when an error occured, this could be a build error or a run error. The error
	// details can be found within the service's metadata
	Error
)

type Namespace struct {
	// Name of the namespace
	Name string
}

type NetworkPolicy struct {
	// The labels allowed ingress by this policy
	AllowedLabels map[string]string
	// Name of the namespace
	Name string
	// Namespace the network policy belongs to
	Namespace string
}

// Service is runtime resource
type Service struct {
	// Name of the service
	Name string
	// Version of the service
	Version string
	// url location of source
	Source string
	// Metadata stores metadata
	Metadata map[string]string
	// Status of the service
	Status ServiceStatus
}

// Resource represents any resource handled by runtime
type Resource struct {
	Namespace     *Namespace
	NetworkPolicy *NetworkPolicy
	Service       *Service
}

// Resources which are allocated to a serivce
type Resources struct {
	// CPU is the maximum amount of CPU the service will be allocated (unit millicpu)
	// e.g. 0.25CPU would be passed as 250
	CPU int
	// Mem is the maximum amount of memory the service will be allocated (unit mebibyte)
	// e.g. 128 MiB of memory would be passed as 128
	Mem int
	// Disk is the maximum amount of disk space the service will be allocated (unit mebibyte)
	// e.g. 128 MiB of memory would be passed as 128
	Disk int
}
