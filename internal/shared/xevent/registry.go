package xevent

import "fmt"

// Registry defines the interface for registering and retrieving payload factories.
type Registry interface {
	Register(eventType string, factory func() interface{})
	GetFactory(eventType string) (func() interface{}, error)
}

// DefaultPayloadFactoryRegistry is the default implementation of PayloadFactoryRegistry.
type DefaultPayloadFactoryRegistry struct {
	factories map[string]func() interface{}
}

// NewPayloadRegistry creates a new instance of DefaultPayloadFactoryRegistry.
func NewPayloadRegistry() *DefaultPayloadFactoryRegistry {
	return &DefaultPayloadFactoryRegistry{
		factories: make(map[string]func() interface{}),
	}
}

// Register adds a payload factory function for the given event type.
func (r *DefaultPayloadFactoryRegistry) Register(eventType string, factory func() interface{}) {
	r.factories[eventType] = factory
}

// GetFactory retrieves the payload factory function for the given event type.
func (r *DefaultPayloadFactoryRegistry) GetFactory(eventType string) (func() interface{}, error) {
	factory, exists := r.factories[eventType]
	if !exists {
		return nil, fmt.Errorf("unknown event type: %s", eventType)
	}
	return factory, nil
}

// Register registers a payload factory function for the given event type.
func Register[Event any](registry Registry, name string, factory func() Event) {
	registry.Register(name, func() interface{} {
		return factory()
	})
}

// typeOf returns the name of the type of the given value.
// If its any, it returns the fmt.Sprintf("%T", value).
// If its a pointer, it returns the name of the type of the pointer.
// If its a Stringer, it returns the result of calling the String method.
func typeOf(value any) string {
	switch v := value.(type) {
	case fmt.Stringer:
		return v.String()
	case fmt.GoStringer:
		return v.GoString()
	default:
		return fmt.Sprintf("%T", value)
	}
}
