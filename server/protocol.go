package main

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const (
	RequestPutCommand = "put"
	RequestGetCommand = "get"

	ResponseStatusOK  = "ok"
	ResponseStatusErr = "error"
)

type request struct {
	command   string
	key       string
	value     float64
	timestamp time.Time
}

type response struct {
	status, errMsg string
	records        []model
}

func newResponse() response {
	return response{}
}

func (r request) validate() (err error) {
	switch r.command {
	case RequestGetCommand:
		if r.value != 0 || !r.timestamp.IsZero() {
			err = errors.New(BadRequestFormatErr)
		} else if r.key == "" {
			err = errors.New(BlankKeyErr)
		}
	case RequestPutCommand:
		timeDiff := time.Now().Sub(r.timestamp).Hours()
		if r.key == "" {
			err = errors.New(BlankKeyErr)
		} else if timeDiff < 0 || r.timestamp.IsZero() {
			err = errors.New(BadTimestampErr)
		}
	default:
		err = errors.New(WrongCommandErr)
	}

	return
}

func (r *response) setOk() {
	r.status = ResponseStatusOK
}

func (r *response) setErr(msg string) {
	r.status = ResponseStatusErr
	r.errMsg = msg
}

func (r *response) build() string {
	var builder strings.Builder
	builder.WriteString(r.status)

	if r.status == ResponseStatusErr {
		builder.WriteString("\n")
		builder.WriteString(r.errMsg)
	} else if r.status == ResponseStatusOK {
		for _, m := range r.records {
			builder.WriteString("\n")
			builder.WriteString(m.key)
			builder.WriteString(" ")
			builder.WriteString(strconv.FormatFloat(m.value, 'f', -1, 64))
			builder.WriteString(" ")
			builder.WriteString(strconv.FormatInt(m.timestamp.Unix(), 10))
		}
	}

	builder.WriteString("\n\n")
	return builder.String()
}
