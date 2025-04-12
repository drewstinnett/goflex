package goflex

import (
	"net/http"
	"strconv"
	"time"
)

func mustNewRequest(method, url string) *http.Request {
	got, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err)
	}
	return got
}

func dateFromUnixString(s string) (*time.Time, error) {
	var vInt64 int64
	if s != "" {
		var err error
		vInt64, err = strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}
	}
	return toPTR(time.Unix(vInt64, 0)), nil
}

func toPTR[V any](v V) *V {
	return &v
}

func fromPTR[T any](ptr *T) T {
	if ptr != nil {
		return *ptr
	}
	var zero T
	return zero
}

const DAY_HOURS = 24

func daysToDuration(days int) time.Duration {
	return time.Duration(days) * time.Hour * DAY_HOURS
}
