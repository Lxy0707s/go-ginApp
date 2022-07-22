package logtool

import "github.com/getsentry/raven-go"

// Logger provide same apis to log message
type Logger interface {
	Debug(msg string, values ...interface{})
	Info(msg string, values ...interface{})
	Warn(msg string, values ...interface{})
	Error(msg string, values ...interface{})
	Fatal(msg string, values ...interface{})
	Sync()
	NewPrefix(string) Logger
	SetDebug(bool)
	SetRavenClient(*raven.Client)
}
