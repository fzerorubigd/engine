package mockery

import (
	"context"

	"google.golang.org/grpc/metadata"
)

// AuthorizeToken add to authorize token to context
func AuthorizeToken(ctx context.Context, token string) context.Context {
	md := metadata.MD{
		"authorization": []string{"Bearer " + token},
	}

	return metadata.NewOutgoingContext(ctx, md)
}
