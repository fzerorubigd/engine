package impl

import (
	"context"

	"github.com/gogo/protobuf/types"

	"elbix.dev/engine/modules/misc/proto"
	"elbix.dev/engine/pkg/grpcgw"
	"elbix.dev/engine/pkg/health"
	"elbix.dev/engine/pkg/version"
)

type miscController struct {
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

func init() {
	grpcgw.Register(miscpb.NewWrappedMiscSystemServer(&miscController{}))
}
