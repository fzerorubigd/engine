package impl

import (
	"context"
	"crypto"
	"crypto/x509"
	"encoding/pem"

	"elbix.dev/engine/pkg/assert"

	miscpb "elbix.dev/engine/modules/misc/proto"
	"elbix.dev/engine/pkg/health"
	"elbix.dev/engine/pkg/version"
	"github.com/gogo/protobuf/types"
)

type miscController struct {
	pub crypto.PublicKey
}

func (mc miscController) PublicKey(context.Context, *miscpb.PubKeyRequest) (*miscpb.PubKeyResponse, error) {
	resp := &miscpb.PubKeyResponse{}
	pubBytes, err := x509.MarshalPKIXPublicKey(mc.pub)
	assert.Nil(err)
	resp.Pub = string(pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubBytes,
		},
	))
	return resp, nil
}

func (mc miscController) Health(ctx context.Context, _ *miscpb.HealthRequest) (*miscpb.HealthResponse, error) {
	err := health.Healthy(ctx)
	if err != nil {
		return nil, err
	}

	return &miscpb.HealthResponse{}, nil
}

func (mc miscController) Version(context.Context, *miscpb.VersionRequest) (*miscpb.VersionResponse, error) {
	ver := version.GetVersion()
	bd, _ := types.TimestampProto(ver.BuildDate)
	cd, _ := types.TimestampProto(ver.Date)
	return &miscpb.VersionResponse{
		BuildDate:  bd,
		CommitDate: cd,
		CommitHash: ver.Hash,
		ShortHash:  ver.Short,
		Count:      ver.Count,
	}, nil
}

// NewMiscController return a new misc controller
func NewMiscController(pub crypto.PublicKey) miscpb.MiscSystemServer {
	return &miscController{pub: pub}
}
