package logger

type Logger interface {
	Debug(msg string, args ...string)
	Info(msg string, args ...string)
	Warn(msg string, args ...string)
	Error(msg string, args ...string)
	Panic(msg string, args ...string)
	Fatal(msg string, args ...string)
}
