package stackdriver

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/kr/pretty"

	"github.com/sirupsen/logrus"
)

func TestFormatter(t *testing.T) {
	skipTimestamp = true

	for _, tt := range formatterTests {
		var out bytes.Buffer

		logger := logrus.New()
		logger.Out = &out
		logger.Formatter = NewFormatter(
			WithService("test"),
			WithVersion("0.1"),
		)

		tt.run(logger)

		var got map[string]interface{}
		json.Unmarshal(out.Bytes(), &got)

		if !reflect.DeepEqual(got, tt.out) {
			t.Errorf("unexpected output = %# v; want = %# v", pretty.Formatter(got), pretty.Formatter(tt.out))
		}
	}
}

var formatterTests = []struct {
	run func(*logrus.Logger)
	out map[string]interface{}
}{
	{
		run: func(logger *logrus.Logger) {
			logger.WithField("foo", "bar").Info("my log entry")
		},
		out: map[string]interface{}{
			"severity": "INFO",
			"message":  "my log entry",
			"context": map[string]interface{}{
				"data": map[string]interface{}{
					"foo": "bar",
				},
			},
		},
	},
	{
		run: func(logger *logrus.Logger) {
			logger.WithField("foo", "bar").Error("my log entry")
		},
		out: map[string]interface{}{
			"severity": "ERROR",
			"message":  "my log entry",
			"serviceContext": map[string]interface{}{
				"service": "test",
				"version": "0.1",
			},
			"context": map[string]interface{}{
				"data": map[string]interface{}{
					"foo": "bar",
				},
				"reportLocation": map[string]interface{}{
					"filePath":     "github.com/TV4/logrus-stackdriver-formatter/formatter_test.go",
					"lineNumber":   59.0,
					"functionName": "glob..func2",
				},
			},
		},
	},
	{
		run: func(logger *logrus.Logger) {
			logger.
				WithField("foo", "bar").
				WithError(errors.New("test error")).
				Error("my log entry")
		},
		out: map[string]interface{}{
			"severity": "ERROR",
			"message":  "my log entry: test error",
			"serviceContext": map[string]interface{}{
				"service": "test",
				"version": "0.1",
			},
			"context": map[string]interface{}{
				"data": map[string]interface{}{
					"foo": "bar",
				},
				"reportLocation": map[string]interface{}{
					"filePath":     "github.com/TV4/logrus-stackdriver-formatter/formatter_test.go",
					"lineNumber":   85.0,
					"functionName": "glob..func3",
				},
			},
		},
	},
	{
		run: func(logger *logrus.Logger) {
			logger.
				WithFields(logrus.Fields{
					"foo": "bar",
					"httpRequest": map[string]interface{}{
						"method": "GET",
					},
				}).
				Error("my log entry")
		},
		out: map[string]interface{}{
			"severity": "ERROR",
			"message":  "my log entry",
			"serviceContext": map[string]interface{}{
				"service": "test",
				"version": "0.1",
			},
			"context": map[string]interface{}{
				"data": map[string]interface{}{
					"foo": "bar",
				},
				"httpRequest": map[string]interface{}{
					"method": "GET",
				},
				"reportLocation": map[string]interface{}{
					"filePath":     "github.com/TV4/logrus-stackdriver-formatter/formatter_test.go",
					"lineNumber":   115.0,
					"functionName": "glob..func4",
				},
			},
		},
	},
	{
		run: func(logger *logrus.Logger) {
			logger.
				WithField("user", "testuser").
				Error("my log entry")
		},
		out: map[string]interface{}{
			"severity": "ERROR",
			"message":  "my log entry",
			"serviceContext": map[string]interface{}{
				"service": "test",
				"version": "0.1",
			},
			"context": map[string]interface{}{
				"user": "testuser",
				"reportLocation": map[string]interface{}{
					"filePath":     "github.com/TV4/logrus-stackdriver-formatter/formatter_test.go",
					"lineNumber":   143.0,
					"functionName": "glob..func5",
				},
			},
		},
	},
	{
		run: func(logger *logrus.Logger) {
			logger.
				WithField("stack_trace", "goroutine 1 [running]:\nmain.main()\n\t/tmp/sandbox3332884018/prog.go:11 +0x45\n").
				WithError(errors.New("test error")).
				Error("my log entry")
		},
		out: map[string]interface{}{
			"severity":    "ERROR",
			"message":     "my log entry: test error",
			"stack_trace": "my log entry: test error\ngoroutine 1 [running]:\nmain.main()\n\t/tmp/sandbox3332884018/prog.go:11 +0x45\n",
			"serviceContext": map[string]interface{}{
				"service": "test",
				"version": "0.1",
			},
			"context": map[string]interface{}{
				"reportLocation": map[string]interface{}{
					"filePath":     "github.com/TV4/logrus-stackdriver-formatter/formatter_test.go",
					"lineNumber":   167.0,
					"functionName": "glob..func6",
				},
			},
		},
	},
}
