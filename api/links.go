package api

import (
	"github.com/domeos/alarm/g"
	"github.com/toolkits/net/httplib"
	"time"
)

func LinkToSMS(content string) (string, error) {
	uri := g.Config().Api.DomeOS + "/api/alarm/link/store"
	req := httplib.Post(uri).SetTimeout(3*time.Second, 10*time.Second)
	req.Body([]byte(content))
	return req.String()
}
