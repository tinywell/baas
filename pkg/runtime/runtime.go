package runtime

import "context"

// ServiceRunner ...
type ServiceRunner interface {
	Run(ctx context.Context, data ServiceMetadata) error
}

// ServiceMetadata ...
type ServiceMetadata interface {
	ServiceType() string
	RuntimeType() string
	DataID() string
	Action() string
}
