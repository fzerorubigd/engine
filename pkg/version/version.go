package version

import (
	"context"
	"strconv"
	"time"

	"github.com/fzerorubigd/balloon/pkg/log"

	"github.com/fzerorubigd/balloon/pkg/initializer"
	"go.uber.org/zap/zapcore"
)

// Data is the application version in detail
type Data struct {
	Hash      string    `json:"hash"`
	Short     string    `json:"short_hash"`
	Date      time.Time `json:"commit_date"`
	Count     int64     `json:"build_number"`
	BuildDate time.Time `json:"build_date"`
}

// GetVersion return the application version in detail
func GetVersion() Data {
	c := Data{}
	c.Count, _ = strconv.ParseInt(count, 10, 64)
	c.Date, _ = time.Parse("01-02-06-15-04-05", date)
	c.Hash = hash
	c.Short = short
	c.BuildDate, _ = time.Parse("01-02-06-15-04-05", build)

	return c
}

// LogVersion return an logrus entry with version information attached
func LogVersion() []zapcore.Field {
	ver := GetVersion()
	return []zapcore.Field{
		log.String("commit_hash", ver.Hash),
		log.String("short_hash", ver.Short),
		log.Time("commit_date", ver.Date),
		log.Time("build_time", ver.BuildDate),
	}
}

type show struct {
}

func (show) Initialize(ctx context.Context) {
	log.Debug("Start", LogVersion()...)
	go func() {
		<-ctx.Done()
		log.Debug("Done", LogVersion()...)
	}()
}

func init() {
	initializer.Register(&show{}, -10)
}
