package version

import "os"

var (
	hash  = "NOT_EXTRACTED"
	short = "NOT_EXTRACTED"
	date  = "04-08-17-00-10-22"
	build = "04-08-17-00-10-22"
	count = "115"
)

func init() {
	if o := os.Getenv("LONGHASH"); o != "" {
		hash = o
	}

	if o := os.Getenv("SHORTHASH"); o != "" {
		short = o
	}

	if o := os.Getenv("COMMITDATE"); o != "" {
		date = o
	}

	if o := os.Getenv("COMMITCOUNT"); o != "" {
		count = o
	}

}
