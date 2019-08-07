package grpcgw

import (
	"github.com/fzerorubigd/engine/pkg/config"
)

var (
	maxPerPage = config.RegisterInt("grpcw.max_per_page", 100, "http maximum item per page")
	minPerPage = config.RegisterInt("grpcw.min_per_page", 1, "http minimum item per page")
	perPage    = config.RegisterInt("grpcw.per_page", 10, "http default item per page")
)

// Pager is interface for payload with pagination support
type Pager interface {
	GetPage() int64
	GetPerPage() int64
}

// GetPageAndCount return the p and c variable from the request, if not available
// return the default value
func GetPageAndCount(pager Pager, offset bool) (int64, int64) {
	p := pager.GetPage()
	if p < 1 {
		p = 1
	}

	c := pager.GetPerPage()
	if c > maxPerPage.Int64() || c < minPerPage.Int64() {
		c = perPage.Int64()
	}

	if offset {
		// If i need to make it to offset model then do it here
		p = (p - 1) * c
	}

	return p, c
}
