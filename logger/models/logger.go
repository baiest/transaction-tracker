package models

type Properties interface {
	LogProperties() map[string]string
}

type LogProperties struct {
	Event            string
	Error            error
	AdditionalParams []Properties
}

type LogService interface {
	Log(string, LogProperties)
}

type Logger struct {
	Service LogService
}

func (l *Logger) Info(props LogProperties) {
	l.Service.Log("info", props)
}

func (l *Logger) Error(props LogProperties) {
	l.Service.Log("error", props)
}

func (l *Logger) Panic(props LogProperties) {
	l.Service.Log("panic", props)
}
