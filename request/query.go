package request

import (
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

func ExtractSlotQueryParams(req *http.Request) (time.Time, time.Duration, error) {
	start, err := extractNumParam(req, "start_timestamp")
	if err != nil {
		return time.Now(), 0, errors.Wrap(err, "failed to extract start_timestamp query param")
	}

	dur, err := extractNumParam(req, "duration")
	if err != nil {
		return time.Now(), 0, errors.Wrap(err, "failed to extract duration query param")
	}

	startTime := time.Unix(start, 0)
	durTime := time.Duration(dur * 1e9) //seconds to nanoseconds

	return startTime, durTime, nil
}

func extractNumParam(req *http.Request, param string) (int64, error) {
	vals, ok := req.URL.Query()[param]
	if !ok || len(vals[0]) < 1 {
		return 0, errors.Errorf("request is missing the '%s' query parameter", param)
	}
	valString := vals[0] //should just be one, but take the first one regardless

	val, err := strconv.ParseInt(valString, 10, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "invalid query parameter '%d'", val)
	}
	return val, nil
}
