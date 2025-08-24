package logger

import (
	"context"
	"sync"
	"transaction-tracker/logger/models"
	"transaction-tracker/logger/services"
)

var (
	lock = &sync.Mutex{}

	logger *models.Logger
)

func GetLogger(ctx context.Context, serviceName string) (*models.Logger, error) {
	lock.Lock()
	defer lock.Unlock()

	if logger == nil {
		service, err := services.NewFmtLogger(ctx, serviceName)
		if err != nil {
			return nil, err
		}

		logger = &models.Logger{
			Service: service,
		}
	}

	return logger, nil
}

type PropertiesImpl struct {
	props map[string]string
}

func (p *PropertiesImpl) LogProperties() map[string]string {
	return p.props
}

func MapToProperties(props map[string]string) models.Properties {
	return &PropertiesImpl{props: props}
}
