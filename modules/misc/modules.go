package misc

import (
	// protobuf/grpc code
	"elbix.dev/engine/modules/misc/impl"
	_ "elbix.dev/engine/modules/misc/proto"
)

// SetPublicKey set the public key for the route
func SetPublicKey(pub string) error {
	return impl.SetPublicKey(pub)
}
