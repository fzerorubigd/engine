package impl

import (
	"context"

	"github.com/fzerorubigd/balloon/modules/misc/proto"
	"github.com/fzerorubigd/balloon/pkg/grpcgw"
	"github.com/fzerorubigd/balloon/pkg/version"
	"github.com/gogo/protobuf/types"
)

type miscController struct {
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
