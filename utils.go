package plexrando

import (
	"net/http"
	"strconv"
	"time"
)

/*
func mustNewRequestWithBody(method, url string, body io.Reader) *http.Request {
	got, err := http.NewRequest(method, url, body)
	if err != nil {
		panic(err)
	}
	return got
}
*/

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
