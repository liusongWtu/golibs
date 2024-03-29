package logs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

const (
	apacheFormatPattern = "%s - - [%s] \"%s %d %d\" %f %s %s\n"
	apacheFormat        = "APACHE_FORMAT"
	jsonFormat          = "JSON_FORMAT"
)

// AccessLogRecord struct for holding access log data.
type AccessLogRecord struct {
	RemoteAddr     string        `json:"remote_addr"`
	RequestTime    time.Time     `json:"request_time"`
	RequestMethod  string        `json:"request_method"`
	Request        string        `json:"request"`
	ServerProtocol string        `json:"server_protocol"`
	Host           string        `json:"host"`
	Status         int           `json:"status"`
	BodyBytesSent  int64         `json:"body_bytes_sent"`
	ElapsedTime    time.Duration `json:"elapsed_time"`
	HTTPReferrer   string        `json:"http_referrer"`
	HTTPUserAgent  string        `json:"http_user_agent"`
	RemoteUser     string        `json:"remote_user"`
}

func (r *AccessLogRecord) json() ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	disableEscapeHTML(encoder)

	err := encoder.Encode(r)
	return buffer.Bytes(), err
}

func disableEscapeHTML(i interface{}) {
	e, ok := i.(interface {
		SetEscapeHTML(bool)
	})
	if ok {
		e.SetEscapeHTML(false)
	}
}

// AccessLog - Format and print access log.
func AccessLog(r *AccessLogRecord, format string) {
	var msg string

	switch format {

	case apacheFormat:
		timeFormatted := r.RequestTime.Format("02/Jan/2006 03:04:05")
		msg = fmt.Sprintf(apacheFormatPattern, r.RemoteAddr, timeFormatted, r.Request, r.Status, r.BodyBytesSent,
			r.ElapsedTime.Seconds(), r.HTTPReferrer, r.HTTPUserAgent)
	case jsonFormat:
		fallthrough
	default:
		jsonData, err := r.json()
		if err != nil {
			msg = fmt.Sprintf(`{"Error": "%s"}`, err)
		} else {
			msg = string(jsonData)
		}
	}

	beeLogger.Debug(msg)
}
