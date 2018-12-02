package initializer

// Simple is the simple initializer for system
type Simple interface {
	Initialize()
}

// DoInitialize is a helper function to initialize the object if it have an initialize method
func DoInitialize(in interface{}) {
	if i, ok := in.(Simple); ok {
		i.Initialize()
	}
}
