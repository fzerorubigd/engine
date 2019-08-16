package misc

import (
	"crypto"

	// protobuf/grpc code
	"elbix.dev/engine/modules/misc/impl"
	_ "elbix.dev/engine/modules/misc/proto"
)

// SetPublicKey set the public key for the route
func SetPublicKey(pub crypto.PublicKey) {
	impl.SetPublicKey(pub)
}
