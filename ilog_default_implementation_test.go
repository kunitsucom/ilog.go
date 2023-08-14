package ilog //nolint:testpackage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"path"
	"regexp"
	"testing"
	"time"
)

func TestScenario(t *testing.T) {
	t.Parallel()
	t.Run("success,JSON", func(t *testing.T) {
		t.Parallel()
		buf := bytes.NewBuffer(nil)
		expected := regexp.MustCompilePOSIX(`{"severity":"DEBUG","timestamp":"[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}\.?[0-9]*Z","caller":"ilog\.go/[a-z_]+_test\.go:[0-9]+","message":"Logf: format string","bool":true,"boolPointer":false,"boolPointer2":null,"byte":"\\u0001","bytes":"bytes","time\.Duration":"1h1m1.001001001s","error":"ilog: log entry not written","errorFormatter":"ilog: log entry not written","errorNil":"<nil>","float32":1\.234567,"float64":1\.23456789,"float64NaN":"NaN","float64\+Inf":"\+Inf","float64-Inf":"-Inf","int":-1,"int8":-1,"int16":-1,"int32":123456789,"int64":123456789,"string":"string","stringEscaped":"\\b\\f\\n\\r\\t","time\.Time":"2023-08-13T04:38:39\.123456789\+09:00","uint":1,"uint16":1,"uint32":123456789,"uint64":123456789,"fmt\.Formatter":"testFormatter","fmt\.Stringer":"testStringer","fmt\.Stringer2":"<nil>","func":"0x[0-9a-f]+"}`)

		le := NewBuilder(DebugLevel, buf).
			SetTimestampZone(time.UTC).
			Build().
			Any("bool", true).
			Any("boolPointer", new(bool)).
			Any("boolPointer2", (*bool)(nil)).
			Any("byte", byte(1)).
			Any("bytes", []byte("bytes")).
			Any("time.Duration", time.Hour+time.Minute+time.Second+time.Millisecond+time.Microsecond+time.Nanosecond).
			Any("error", ErrLogEntryIsNotWritten).
			Any("errorFormatter", &testFormatterError{ErrLogEntryIsNotWritten}).
			Any("errorNil", nil).
			Any("float32", float32(1.234567)).
			Any("float64", float64(1.23456789)).
			Any("float64NaN", math.NaN()).
			Any("float64+Inf", math.Inf(1)).
			Any("float64-Inf", math.Inf(-1)).
			Any("int", int(-1)).
			Any("int8", int8(-1)).
			Any("int16", int16(-1)).
			Any("int32", int32(123456789)).
			Any("int64", int64(123456789)).
			Any("string", "string").
			Any("stringEscaped", "\b\f\n\r\t").
			Any("time.Time", time.Date(2023, 8, 13, 4, 38, 39, 123456789, time.FixedZone("Asia/Tokyo", int(9*time.Hour/time.Second)))).
			Any("uint", uint(1)).
			Any("uint16", uint16(1)).
			Any("uint32", uint32(123456789)).
			Any("uint64", uint64(123456789)).
			Any("fmt.Formatter", &testFormatter{}).
			Any("fmt.Stringer", testStringer("testStringer")).
			Any("fmt.Stringer2", (*testStringer)(nil)).
			Any("func", func() {})

		le.Logf(DebugLevel, "Logf: %s", "format string")
		if !expected.Match(buf.Bytes()) {
			t.Errorf("❌: !expected.Match(buf.Bytes()):\n%s", buf)
		}

		decoded := make(map[string]interface{})
		if err := json.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(&decoded); err != nil {
			t.Errorf("❌: err: %+v", err)
		}

		if expected, actual := "DEBUG", decoded["severity"]; expected != actual {
			t.Errorf("❌: severity: expected(%s) != actual(%s)", expected, actual)
		}

		if expected, actual := "Logf: format string", decoded["message"]; expected != actual {
			t.Errorf("❌: message: expected(%s) != actual(%s)", expected, actual)
		}

		if expected, actual := "+Inf", decoded["float64+Inf"]; expected != actual {
			t.Errorf("❌: float64+Inf: expected(%s) != actual(%s)", expected, actual)
		}

		t.Logf("ℹ️: buf:\n%s", buf)
	})
}

func TestLogger(t *testing.T) {
	t.Parallel()
	t.Run("success,Logger", func(t *testing.T) {
		t.Parallel()
		buf := bytes.NewBuffer(nil)
		expected := regexp.MustCompilePOSIX(`{"level":"DEBUG","time":"[A-Z][a-z]+ [A-Z][a-z]+ [0-9]{1,2} [0-9]{2}:[0-9]{2}:[0-9]{2} UTC [0-9]{4}","file":".+/ilog\.go/[a-z_]+_test\.go:[0-9]+","msg":"Logf"}` + "\r\n")

		const expectedLevel = DebugLevel
		l := NewBuilder(ErrorLevel, buf).
			SetLevelKey("level").
			SetLevels(copyLevels(defaultLevels)).
			SetTimestampKey("time").
			SetTimestampFormat(time.UnixDate).
			SetTimestampZone(time.UTC).
			SetCallerKey("file").
			UseShortCaller(false).
			SetMessageKey("msg").
			SetSeparator("\r\n").
			Build().
			SetLevel(expectedLevel).
			AddCallerSkip(10).
			AddCallerSkip(-10)

		if expected, actual := expectedLevel, l.Level(); expected != actual {
			t.Errorf("❌: expected(%d) != actual(%d)", expected, actual)
		}

		l.Logf(DebugLevel, "Logf")
		if !expected.Match(buf.Bytes()) {
			t.Errorf("❌: !expected.Match(buf.Bytes()):\n%s", buf)
		}

		t.Logf("ℹ️: buf:\n%s", buf)
	})

	t.Run("success,Logger,common", func(t *testing.T) {
		t.Parallel()
		buf := bytes.NewBuffer(nil)
		l := NewBuilder(DebugLevel, buf).SetTimestampKey("").SetCallerKey("").Build()
		l.Any("any", "any").Debugf("Debugf")
		l.Bool("bool", true).Debugf("Debugf")
		l.Bytes("bytes", []byte("bytes")).Debugf("Debugf")
		l.Duration("time.Duration", time.Hour+time.Minute+time.Second+time.Millisecond+time.Microsecond+time.Nanosecond).Debugf("Debugf")
		l.Err(io.ErrUnexpectedEOF).Debugf("Debugf")
		l.ErrWithKey("err", io.ErrUnexpectedEOF).Debugf("Debugf")
		l.Float32("float32", float32(1.234567)).Debugf("Debugf")
		l.Float64("float64", float64(1.23456789)).Debugf("Debugf")
		l.Int("int", int(-1)).Debugf("Debugf")
		l.Int32("int32", int32(-1)).Debugf("Debugf")
		l.Int64("int64", int64(-1)).Debugf("Debugf")
		l.String("string", "string").Debugf("Debugf")
		l.Time("time.Time", time.Date(2023, 8, 13, 4, 38, 39, 123456789, time.FixedZone("Asia/Tokyo", int(9*time.Hour/time.Second)))).Debugf("Debugf")
		l.Uint("uint", uint(1)).Debugf("Debugf")
		l.Uint32("uint32", uint32(123456789)).Debugf("Debugf")
		l.Uint64("uint64", uint64(123456789)).Debugf("Debugf")
		l.Debugf("Debugf")
		l.Infof("Infof")
		l.Warnf("Warnf")
		l.Errorf("Errorf")
		l.Any("any", "any").Debugf("Debugf")
		l.Any("any", "any").Infof("Infof")
		l.Any("any", "any").Warnf("Warnf")
		l.Any("any", "any").Errorf("Errorf")
		l.Any("any", "any").Logf(DebugLevel, "Logf")
		_, _ = l.Any("any", "any").Write([]byte("Write"))

		t.Logf("ℹ️: buf:\n%s", buf)
	})
}

type testFormatter struct{}

func (f *testFormatter) Format(s fmt.State, verb rune) {
	_, _ = fmt.Fprint(s, "testFormatter")
}

type testFormatterError struct {
	err error
}

func (f *testFormatterError) Format(s fmt.State, verb rune) {
	_, _ = fmt.Fprint(s, f.err)
}

func (f *testFormatterError) Error() string {
	return f.err.Error()
}

type testStringer string

func (s testStringer) String() string {
	return string(s)
}

func TestLogEntry(t *testing.T) {
	t.Parallel()
	t.Run("success,LogEntry", func(t *testing.T) {
		t.Parallel()
		buf := bytes.NewBuffer(nil)
		expected := regexp.MustCompilePOSIX(`{"severity":"DEBUG","timestamp":"[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}\.?[0-9]*Z","caller":"ilog\.go/[a-z_]+_test\.go:[0-9]+","message":"Logf: format string","bool":true,"boolPointer":false,"boolPointer2":null,"byte":"\\u0001","bytes":"bytes","time\.Duration":"1h1m1.001001001s","error":"ilog: log entry not written","errorFormatter":"ilog: log entry not written","errorNil":"<nil>","float32":1\.234567,"float64":1\.23456789,"float64NaN":"NaN","float64\+Inf":"\+Inf","float64-Inf":"-Inf","int":-1,"int8":-1,"int16":-1,"int32":123456789,"int64":123456789,"string":"string","stringEscaped":"\\b\\f\\n\\r\\t","time\.Time":"2023-08-13T04:38:39\.123456789\+09:00","uint":1,"uint16":1,"uint32":123456789,"uint64":123456789,"fmt\.Formatter":"testFormatter","fmt\.Stringer":"testStringer","fmt\.Stringer2":"<nil>","func":"0x[0-9a-f]+"}`)

		le := NewBuilder(DebugLevel, buf).
			SetTimestampZone(time.UTC).
			Build().
			Any("bool", true).
			Any("boolPointer", new(bool)).
			Any("boolPointer2", (*bool)(nil)).
			Any("byte", byte(1)).
			Any("bytes", []byte("bytes")).
			Any("time.Duration", time.Hour+time.Minute+time.Second+time.Millisecond+time.Microsecond+time.Nanosecond).
			Any("error", ErrLogEntryIsNotWritten).
			Any("errorFormatter", &testFormatterError{ErrLogEntryIsNotWritten}).
			Any("errorNil", nil).
			Any("float32", float32(1.234567)).
			Any("float64", float64(1.23456789)).
			Any("float64NaN", math.NaN()).
			Any("float64+Inf", math.Inf(1)).
			Any("float64-Inf", math.Inf(-1)).
			Any("int", int(-1)).
			Any("int8", int8(-1)).
			Any("int16", int16(-1)).
			Any("int32", int32(123456789)).
			Any("int64", int64(123456789)).
			Any("string", "string").
			Any("stringEscaped", "\b\f\n\r\t").
			Any("time.Time", time.Date(2023, 8, 13, 4, 38, 39, 123456789, time.FixedZone("Asia/Tokyo", int(9*time.Hour/time.Second)))).
			Any("uint", uint(1)).
			Any("uint16", uint16(1)).
			Any("uint32", uint32(123456789)).
			Any("uint64", uint64(123456789)).
			Any("fmt.Formatter", &testFormatter{}).
			Any("fmt.Stringer", testStringer("testStringer")).
			Any("fmt.Stringer2", (*testStringer)(nil)).
			Any("func", func() {})

		if expected, actual := ErrLogEntryIsNotWritten.Error(), le.Error(); expected != actual {
			t.Errorf("❌: expected(%s) != actual(%s)", expected, actual)
		}

		le.Logf(DebugLevel, "Logf: %s", "format string")
		if !expected.Match(buf.Bytes()) {
			t.Errorf("❌: !expected.Match(buf.Bytes()):\n%s", buf)
		}

		t.Logf("ℹ️: buf:\n%s", buf)
	})

	t.Run("success,default,DEBUG", func(t *testing.T) {
		t.Parallel()
		buf := bytes.NewBuffer(nil)
		expected := regexp.MustCompilePOSIX(`{"severity":"DEBUG","timestamp":"[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}\.?[0-9]*Z","caller":"ilog\.go/[a-z_]+_test\.go:[0-9]+","message":"default"}`)

		NewBuilder(-128, buf).SetTimestampZone(time.UTC).Build().Logf(-128, "default")

		if !expected.Match(buf.Bytes()) {
			t.Errorf("❌: !expected.Match(buf.Bytes()):\n%s", buf)
		}
	})

	t.Run("success,Debugf", func(t *testing.T) {
		t.Parallel()
		buf := bytes.NewBuffer(nil)
		expected := regexp.MustCompilePOSIX(`{"severity":"DEBUG","timestamp":"[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}\.?[0-9]*Z","caller":"ilog\.go/[a-z_]+_test\.go:[0-9]+","message":"Debugf"}`)

		NewBuilder(DebugLevel, buf).SetTimestampZone(time.UTC).Build().Debugf("Debugf")

		if !expected.Match(buf.Bytes()) {
			t.Errorf("❌: !expected.Match(buf.Bytes()):\n%s", buf)
		}
	})

	t.Run("success,Infof", func(t *testing.T) {
		t.Parallel()
		buf := bytes.NewBuffer(nil)
		expected := regexp.MustCompilePOSIX(`{"severity":"INFO","timestamp":"[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}\.?[0-9]*Z","caller":"ilog\.go/[a-z_]+_test\.go:[0-9]+","message":"Infof"}`)

		NewBuilder(DebugLevel, buf).SetTimestampZone(time.UTC).Build().Infof("Infof")

		if !expected.Match(buf.Bytes()) {
			t.Errorf("❌: !expected.Match(buf.Bytes()):\n%s", buf)
		}
	})

	t.Run("success,Warnf", func(t *testing.T) {
		t.Parallel()
		buf := bytes.NewBuffer(nil)
		expected := regexp.MustCompilePOSIX(`{"severity":"WARN","timestamp":"[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}\.?[0-9]*Z","caller":"ilog\.go/[a-z_]+_test\.go:[0-9]+","message":"Warnf"}`)

		NewBuilder(DebugLevel, buf).SetTimestampZone(time.UTC).Build().Warnf("Warnf")

		if !expected.Match(buf.Bytes()) {
			t.Errorf("❌: !expected.Match(buf.Bytes()):\n%s", buf)
		}
	})

	t.Run("success,Errorf", func(t *testing.T) {
		t.Parallel()
		buf := bytes.NewBuffer(nil)
		expected := regexp.MustCompilePOSIX(`{"severity":"ERROR","timestamp":"[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}\.?[0-9]*Z","caller":"ilog\.go/[a-z_]+_test\.go:[0-9]+","message":"Errorf"}`)

		NewBuilder(DebugLevel, buf).SetTimestampZone(time.UTC).Build().Errorf("Errorf")

		if !expected.Match(buf.Bytes()) {
			t.Errorf("❌: !expected.Match(buf.Bytes()):\n%s", buf)
		}
	})
}

func TestLogger_logf(t *testing.T) {
	t.Parallel()
	t.Run("success,empty", func(t *testing.T) {
		t.Parallel()
		buf := bytes.NewBuffer(nil)
		const expected = ""

		NewBuilder(InfoLevel, buf).SetTimestampZone(time.UTC).Build().Debugf("Debugf")
		if expected != buf.String() {
			t.Errorf("❌: expected(%s) != actual(%s)", expected, buf.String())
		}
		buf.Reset()

		NewBuilder(WarnLevel, buf).SetTimestampZone(time.UTC).Build().Infof("Infof")
		if expected != buf.String() {
			t.Errorf("❌: expected(%s) != actual(%s)", expected, buf.String())
		}
		buf.Reset()

		NewBuilder(ErrorLevel, buf).SetTimestampZone(time.UTC).Build().Warnf("Warnf")
		if expected != buf.String() {
			t.Errorf("❌: expected(%s) != actual(%s)", expected, buf.String())
		}
		buf.Reset()
	})

	t.Run("success,{}", func(t *testing.T) {
		t.Parallel()
		buf := bytes.NewBuffer(nil)
		const expected = "{}\n"

		NewBuilder(DebugLevel, buf).
			SetLevelKey("").
			SetTimestampKey("").
			SetCallerKey("").
			SetMessageKey("").Build().Debugf("{}")

		if expected != buf.String() {
			t.Errorf("❌: expected(%s) != actual(%s)", expected, buf.String())
		}
	})
}

type testWriter struct {
	err error
}

func (w *testWriter) Write(p []byte) (n int, err error) {
	return 0, w.err
}

//nolint:paralleltest,tparallel
func TestLogger_Write(t *testing.T) {
	//nolint:paralleltest,tparallel
	t.Run("failure,Logger,Write", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		defer SetGlobal(NewBuilder(DebugLevel, buf).SetTimestampZone(time.UTC).Build())()

		i, err := NewBuilder(DebugLevel, &testWriter{err: io.ErrUnexpectedEOF}).SetTimestampZone(time.UTC).Build().Write([]byte("ERROR"))
		if expected := regexp.MustCompilePOSIX(`w.logf: w.logger.writer.Write: p={"severity":"DEBUG","timestamp":"[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}\.?[0-9]*Z","caller":"ilog\.go/[a-z_]+_test\.go:[0-9]+","message":"ERROR"}: unexpected EOF`); err == nil || !expected.MatchString(err.Error()) {
			t.Errorf("❌: err != nil: %v", err)
		}
		if i != 0 {
			t.Errorf("❌: i != 0: %d", i)
		}
		if expected := regexp.MustCompilePOSIX(`{"severity":"ERROR","timestamp":"[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}\.?[0-9]*Z","caller":"ilog\.go/[a-z_]+\.go:[0-9]+","message":"w.logger.writer.Write: p={\\"severity\\":\\"DEBUG\\",\\"timestamp\\":\\"[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}\.?[0-9]*Z\\",\\"caller\\":\\"ilog\.go/[a-z_]+_test\.go:[0-9]+\\",\\"message\\":\\"ERROR(\\n)?\\"}: unexpected EOF"}`); !expected.Match(buf.Bytes()) {
			t.Errorf("❌: !expected.Match(buf.Bytes()):\n%s", buf)
		}
		t.Logf("ℹ️: buf:\n%s", buf)
		buf.Reset()
	})

	t.Run("failure,LogEntry,Write", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		defer SetGlobal(NewBuilder(DebugLevel, buf).SetTimestampZone(time.UTC).Build())()

		i, err := NewBuilder(DebugLevel, &testWriter{err: io.ErrUnexpectedEOF}).SetTimestampZone(time.UTC).Build().Any("any", "any").Write([]byte("ERROR"))
		if expected := regexp.MustCompilePOSIX(`w.logf: w.logger.writer.Write: p={"severity":"DEBUG","timestamp":"[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}\.?[0-9]*Z","caller":"ilog\.go/[a-z_]+_test\.go:[0-9]+","message":"ERROR","any":"any"}: unexpected EOF`); err == nil || !expected.MatchString(err.Error()) {
			t.Errorf("❌: err != nil: %v", err)
		}
		if i != 0 {
			t.Errorf("❌: i != 0: %d", i)
		}
		if expected := regexp.MustCompilePOSIX(`{"severity":"ERROR","timestamp":"[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}\.?[0-9]*Z","caller":"ilog\.go/[a-z_]+\.go:[0-9]+","message":"w.logger.writer.Write: p={\\"severity\\":\\"DEBUG\\",\\"timestamp\\":\\"[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}\.?[0-9]*Z\\",\\"caller\\":\\"ilog\.go/[a-z_]+_test\.go:[0-9]+\\",\\"message\\":\\"ERROR(\\n)?\\",\\"any\\":\\"any\\"}: unexpected EOF"}`); !expected.Match(buf.Bytes()) {
			t.Errorf("❌: !expected.Match(buf.Bytes()):\n%s", buf)
		}
		t.Logf("ℹ️: buf:\n%s", buf)
		buf.Reset()
	})
}

func Test_extractShortPath(t *testing.T) {
	t.Parallel()
	t.Run("success,noIndex", func(t *testing.T) {
		t.Parallel()
		const expected = "expected"
		actual := extractShortPath(expected)
		if expected != actual {
			t.Errorf("❌: expected(%s) != actual(%s)", expected, actual)
		}
	})
	t.Run("success,1Index", func(t *testing.T) {
		t.Parallel()
		const expected = "expected/expected"
		actual := extractShortPath(expected)
		if expected != actual {
			t.Errorf("❌: expected(%s) != actual(%s)", expected, actual)
		}
	})

	t.Run("success,1Index", func(t *testing.T) {
		t.Parallel()
		const expected = "expected/expected"
		actual := extractShortPath(path.Join(expected, expected))
		if expected != actual {
			t.Errorf("❌: expected(%s) != actual(%s)", expected, actual)
		}
	})
}
