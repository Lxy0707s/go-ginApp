package logtool

import (
	"github.com/getsentry/raven-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

// Sugar implements Logger with zap sugar logger
type Sugar struct {
	prefix      string
	sugar       *zap.SugaredLogger
	ravenClient *raven.Client
}

// NewSugar return logger with prefix string, and debug level setting
// if prefix is set, print log message like 'prefix : message'
// if set debug to true, print debug level logs. otherwise, print info level and above
func NewSugar(prefix string, debug bool) Logger {
	s := &Sugar{
		prefix: prefix,
	}
	s.initSugar(debug)
	return s
}

func (s *Sugar) SetRavenClient(ravenClient *raven.Client) {
	s.ravenClient = ravenClient
}

func (s *Sugar) initSugar(isDebug bool) {
	config := zap.NewProductionConfig()
	lv := zap.NewAtomicLevel()
	lv.SetLevel(zap.InfoLevel)
	if isDebug {
		lv.SetLevel(zap.DebugLevel)
		// config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	config.Sampling = &zap.SamplingConfig{
		Initial:    1e5,
		Thereafter: 1e5,
	}
	config.Level = lv
	config.Encoding = "console"
	config.EncoderConfig.EncodeTime = timeEncoder
	config.DisableCaller = true
	// config.DisableStacktrace = true

	l, err := config.Build()
	if err != nil {
		panic("log init error : " + err.Error())
	}
	if s.sugar != nil {
		su := l.Sugar()
		*s.sugar = *su
		return
	}
	s.sugar = l.Sugar()
}

// Sync sync logs
func (s *Sugar) Sync() {
	s.sugar.Sync()
}

func (s *Sugar) prefixMsg(msg string) string {
	if s.prefix == "" {
		return msg
	}
	return s.prefix + " : " + msg
}

// Debug logs message as debug level
func (s *Sugar) Debug(msg string, values ...interface{}) {
	msg = s.prefixMsg(msg)
	if len(values) == 0 {
		s.sugar.Debug(msg)
		return
	}
	s.sugar.Debugw(msg, values...)
}

// Info logs message as info level
func (s *Sugar) Info(msg string, values ...interface{}) {
	msg = s.prefixMsg(msg)
	if len(values) == 0 {
		s.sugar.Info(msg)
		return
	}
	s.sugar.Infow(msg, values...)
}

// Warn logs message as warn level
func (s *Sugar) Warn(msg string, values ...interface{}) {
	msg = s.prefixMsg(msg)
	if len(values) == 0 {
		s.sugar.Warn(msg)
		return
	}
	s.sugar.Warnw(msg, values...)
}

// Error logs message as error level
// it logs caller stack together
func (s *Sugar) Error(msg string, values ...interface{}) {
	msg = s.prefixMsg(msg)
	if len(values) == 0 {
		if s.ravenClient != nil {
			s.ravenClient.CaptureMessage(msg, nil)
		}
		s.sugar.Error(msg)
		return
	}

	if s.ravenClient != nil {
		var err error
		var tags map[string]string

		for _, value := range values {
			switch v := value.(type) {
			case error:
				err = v
			case map[string]string:
				tags = v
			default:
				break
			}
		}

		if tags == nil {
			if err == nil {
				s.ravenClient.CaptureMessage(msg, nil)
			} else {
				s.ravenClient.CaptureError(err, nil)
			}
		} else {
			if err == nil {
				s.ravenClient.CaptureMessage(msg, tags)
			} else {
				s.ravenClient.CaptureError(err, tags)
			}
		}
	}

	s.sugar.Errorw(msg, values...)
}

// Fatal logs message as fatal level
// it calls os.Exit(1) after logged
func (s *Sugar) Fatal(msg string, values ...interface{}) {
	msg = s.prefixMsg(msg)
	if len(values) == 0 {
		if s.ravenClient != nil {
			s.ravenClient.CaptureMessageAndWait(msg, nil)
		}
		s.sugar.Fatal(msg)
		return
	}

	if s.ravenClient != nil {
		var err error
		var tags map[string]string

		for _, value := range values {
			switch v := value.(type) {
			case error:
				err = v
			case map[string]string:
				tags = v
			default:
				break
			}
		}

		if tags == nil {
			if err == nil {
				s.ravenClient.CaptureMessageAndWait(msg, nil)
			} else {
				s.ravenClient.CaptureErrorAndWait(err, nil)
			}
		} else {
			if err == nil {
				s.ravenClient.CaptureMessageAndWait(msg, tags)
			} else {
				s.ravenClient.CaptureErrorAndWait(err, tags)
			}
		}
	}

	s.sugar.Fatalw(msg, values...)
}

// NewPrefix return new suger logger with prefix message
func (s *Sugar) NewPrefix(prefix string) Logger {
	return &Sugar{
		sugar:  s.sugar,
		prefix: prefix,
	}
}

// SetDebug sets logger to print debug
func (s *Sugar) SetDebug(debug bool) {
	s.initSugar(debug)
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("01/02 15:04:05.999"))
}
