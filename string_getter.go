package glog

import (
	"encoding/json"
	"fmt"
	"time"

	"git.veep.tech/veep/ent-kiosk-backend/common/sin"
	stringUtil "git.veep.tech/veep/ent-kiosk-backend/common/util/strings"
)

type (
	StringGetter func() string
)

// AsJSON Convert object to JSON when needed
func AsJSON(object interface{}) StringGetter {
	return func() string {
		if object == nil {
			return ""
		}
		data, _ := json.MarshalIndent(object, "", "  ")
		return string(data)
	}
}

// AsISOTime Convert object to JSON when needed
func AsISOTime(t time.Time) StringGetter {
	return func() string {
		return t.Format(time.RFC3339)
	}
}

func AsErrStrackTrace(err error) StringGetter {
	return func() string {
		var errAsString string
		if realErr, isSin := err.(sin.Sin); isSin {
			errAsBytes, _ := json.Marshal(realErr)
			errAsString = fmt.Sprintf("%s: %s", realErr.Error(), string(errAsBytes))
		} else {
			errAsString = err.Error()
		}
		return errAsString
	}
}

func Last(value string, length int) StringGetter {
	return func() string {
		return stringUtil.Last(value, length)
	}
}
