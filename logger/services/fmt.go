package services

import (
	"context"
	"fmt"
	"log"
	"transaction-tracker/logger/models"
)

type FmtLogger struct {
	ServiceName string
}

func NewFmtLogger(_ context.Context, serviceName string) (*FmtLogger, error) {
	return &FmtLogger{
		ServiceName: serviceName,
	}, nil
}

func (l *FmtLogger) Log(level string, props models.LogProperties) {
	stringProps := ""

	for _, additionalParam := range props.AdditionalParams {
		if additionalParam == nil {
			continue
		}

		for key, value := range additionalParam.LogProperties() {
			stringProps += fmt.Sprintf("| %s: %s", key, value)
		}
	}

	if (level == "error" || level == "panic") && props.Error != nil {
		stringProps += fmt.Sprintf("| error: %v", props.Error)
	}

	log.Printf("service: %s | level: %s | event: %s %s\n", l.ServiceName, level, props.Event, stringProps)
}
