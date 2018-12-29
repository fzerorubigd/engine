package version

import "os"

var (
	hash  = "--"
	short = "--"
	date  = "--"
	build = "--"
	count = "-1"
)

func init() {
	if hash == "--" {
		if o := os.Getenv("LONG_HASH"); o != "" {
			hash = o
		}
	}

	if short == "--" {
		if o := os.Getenv("SHORT_HASH"); o != "" {
			short = o
		}
	}

	if date == "--" {
		if o := os.Getenv("COMMIT_DATE"); o != "" {
			date = o
		}
	}

	if build == "--" {
		if o := os.Getenv("BUILD_DATE"); o != "" {
			build = o
		}
	}

	if count == "-1" {
		if o := os.Getenv("COMMITCOUNT"); o != "" {
			count = o
		}
	}

}
