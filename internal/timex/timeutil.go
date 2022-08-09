package timex

import (
	"time"

	"github.com/jackc/pgtype"
	"github.com/spotlightpa/almanack/internal/syncx"
	"github.com/spotlightpa/almanack/internal/try"
)

var getNewYork = syncx.Once(func() *time.Location {
	return try.To(time.LoadLocation("America/New_York"))
})

func ToEST(t time.Time) time.Time {
	return t.In(getNewYork())
}

func Unwrap(v any) (t time.Time, ok bool) {
	if t, _ = v.(time.Time); !t.IsZero() {
		ok = true
		return
	}
	s, _ := v.(string)
	if s == "" {
		return
	}
	ok = t.UnmarshalText([]byte(s)) == nil
	return
}

const timeWindow = 5 * time.Minute

func Equalish(old, new pgtype.Timestamptz) bool {
	if old.Status != new.Status {
		return false
	}
	if old.Status != pgtype.Present {
		return true
	}
	diff := old.Time.Sub(new.Time).Abs()
	return diff < timeWindow
}
