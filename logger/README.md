# Logger Package

Un sistema de logging robusto y extensible para Go que implementa el patrón adaptador para mantener la separación de responsabilidades.

## Características

- ✅ **Múltiples niveles de logging**: Debug, Info, Warning, Error, Fatal
- ✅ **Eventos estructurados**: Cada log incluye evento, mensaje, TTL y parámetros
- ✅ **Patrón adaptador**: Sentry y otras integraciones no contaminan el logger principal
- ✅ **Modo debug**: Logging en consola con formato JSON estructurado
- ✅ **Integraciones extensibles**: Fácil agregar nuevos servicios de logging
- ✅ **Contexto opcional**: Soporte para context.Context en logs
- ✅ **TTL (Time To Live)**: Configuración de tiempo de vida para logs
- ✅ **Graceful degradation**: El logger funciona incluso si las integraciones fallan

## Instalación

```bash
go get github.com/getsentry/sentry-go
go get github.com/sirupsen/logrus
```

## Uso Básico

### Logger Simple

```go
package main

import (
    "time"
    "your-project/logger"
)

func main() {
    // Configuración básica
    config := &logger.LoggerConfig{
        Level:      logger.InfoLevel,
        DebugMode:  true,
        Console:    true,
    }

    // Crear logger
    log, err := logger.NewLogger(config)
    if err != nil {
        panic(err)
    }
    defer log.Close()

    // Logging básico
    log.Info("app_start", "Application started", 0, map[string]interface{}{
        "version": "1.0.0",
        "port":    8080,
    })

    // Logging con TTL
    log.Info("temp_data", "Temporary data created", 30*time.Minute, map[string]interface{}{
        "data_type": "cache",
        "size":      "1MB",
    })

    // Logging de errores
    log.Error("db_error", "Database connection failed", 6*time.Hour, map[string]interface{}{
        "database": "postgres",
        "host":     "localhost:5432",
    })
}
```

### Logger Global

```go
package main

import (
    "your-project/logger"
)

func main() {
    // Inicializar logger global
    if err := logger.InitGlobalLogger(); err != nil {
        panic(err)
    }
    defer logger.CloseGlobalLogger()

    // Usar funciones de conveniencia
    logger.Info("app_start", "Application started", 0, nil)
    logger.Debug("config_loaded", "Configuration loaded", 0, nil)
    logger.Error("api_error", "API request failed", 0, nil)
}
```

## Integración con Sentry

### Usando el Adaptador

```go
package main

import (
    "your-project/logger"
)

func main() {
    // Configuración de Sentry
    sentryConfig := &logger.SentryConfig{
        DSN:              "https://your-dsn@sentry.io/project",
        Environment:      "production",
        Release:          "v1.0.0",
        Debug:            false,
        TracesSampleRate: 0.1,
        AttachStacktrace: true,
    }

    // Crear adaptador de Sentry
    sentryAdapter, err := logger.NewSentryAdapter(sentryConfig)
    if err != nil {
        panic(err)
    }

    // Configurar logger con Sentry
    config := &logger.LoggerConfig{
        Level:        logger.InfoLevel,
        Console:      true,
        Integrations: []logger.LogIntegration{sentryAdapter},
    }

    log, err := logger.NewLogger(config)
    if err != nil {
        panic(err)
    }
    defer log.Close()

    // Los errores se enviarán automáticamente a Sentry
    log.Error("critical_error", "System failure", 0, map[string]interface{}{
        "component": "database",
        "severity":  "critical",
    })
}
```

### Configuración Automática con Variables de Entorno

```bash
export SENTRY_DSN="https://your-dsn@sentry.io/project"
export SENTRY_ENVIRONMENT="production"
export SENTRY_RELEASE="v1.0.0"
export DEBUG="true"
```

```go
package main

import (
    "your-project/logger"
)

func main() {
    // El logger se configurará automáticamente con Sentry
    if err := logger.InitGlobalLogger(); err != nil {
        panic(err)
    }
    defer logger.CloseGlobalLogger()

    // Sentry está configurado automáticamente
    logger.Error("test_error", "This will go to Sentry", 0, nil)
}
```

## Integraciones Personalizadas

### Crear una Nueva Integración

```go
package main

import (
    "your-project/logger"
)

type CustomIntegration struct {
    webhookURL string
}

func NewCustomIntegration(webhookURL string) *CustomIntegration {
    return &CustomIntegration{webhookURL: webhookURL}
}

func (c *CustomIntegration) Log(event logger.LogEvent) error {
    // Implementar lógica personalizada
    // Por ejemplo, enviar a un webhook
    return nil
}

func (c *CustomIntegration) Close() error {
    // Limpiar recursos
    return nil
}

func main() {
    config := logger.DefaultConfig()
    
    // Agregar integración personalizada
    customIntegration := NewCustomIntegration("https://webhook.example.com")
    config.Integrations = append(config.Integrations, customIntegration)
    
    log, err := logger.NewLogger(config)
    if err != nil {
        panic(err)
    }
    defer log.Close()
    
    // Los logs se enviarán a la integración personalizada
    log.Info("custom_event", "Event sent to custom integration", 0, nil)
}
```

### Integraciones Compuestas

```go
package main

import (
    "your-project/logger"
)

func main() {
    // Crear múltiples integraciones
    fileIntegration := logger.NewFileIntegration("/var/log/app.log")
    slackIntegration := logger.NewSlackIntegration("webhook-url", "#alerts")
    
    // Combinar en una integración compuesta
    compositeIntegration := logger.NewCompositeIntegration(
        fileIntegration,
        slackIntegration,
    )
    
    config := logger.DefaultConfig()
    config.Integrations = append(config.Integrations, compositeIntegration)
    
    log, err := logger.NewLogger(config)
    if err != nil {
        panic(err)
    }
    defer log.Close()
    
    // Los logs se enviarán a todas las integraciones
    log.Error("alert", "System alert", 0, nil)
}
```

### Integración Filtrada

```go
package main

import (
    "your-project/logger"
)

func main() {
    slackIntegration := logger.NewSlackIntegration("webhook-url", "#alerts")
    
    // Solo enviar errores a Slack
    errorOnlySlack := logger.NewFilteredIntegration(
        slackIntegration,
        func(event logger.LogEvent) bool {
            return event.Level >= logger.ErrorLevel
        },
    )
    
    config := logger.DefaultConfig()
    config.Integrations = append(config.Integrations, errorOnlySlack)
    
    log, err := logger.NewLogger(config)
    if err != nil {
        panic(err)
    }
    defer log.Close()
    
    // Solo los errores irán a Slack
    log.Info("info", "This won't go to Slack", 0, nil)
    log.Error("error", "This will go to Slack", 0, nil)
}
```

## Funciones Helper

### Logging de Transacciones

```go
logger.LogTransaction("payment_processed", "txn-123", 99.99, "completed", 7*24*time.Hour)
```

### Logging de Errores

```go
err := errors.New("database connection failed")
logger.LogError("db_error", err, 6*time.Hour, map[string]interface{}{
    "database": "postgres",
})
```

### Logging de Rendimiento

```go
start := time.Now()
// ... operación ...
duration := time.Since(start)

logger.LogPerformance("db_query", "SELECT users", duration, 1*time.Hour, map[string]interface{}{
    "table": "users",
    "rows":  100,
})
```

### Logging de Seguridad

```go
logger.LogSecurity("login_attempt", "user123", "login", "/auth/login", 24*time.Hour, map[string]interface{}{
    "ip":        "192.168.1.100",
    "user_agent": "Mozilla/5.0...",
})
```

## Configuración Avanzada

### Niveles de Logging

```go
config := &logger.LoggerConfig{
    Level:      logger.DebugLevel,  // Debug, Info, Warning, Error, Fatal
    DebugMode:  true,               // Habilita modo debug
    Console:    true,               // Habilita logging en consola
}
```

### Modo Debug

```go
// Habilitar modo debug
logger.SetDebugMode(true)

// O configurar nivel manualmente
logger.SetLevel(logger.DebugLevel)
```

## Patrón Adaptador

El logger implementa el patrón adaptador para mantener la separación de responsabilidades:

- **Logger Core**: Maneja la lógica de logging básica
- **Integrations**: Implementan la interfaz `LogIntegration`
- **Adapters**: Conectan servicios externos (como Sentry) con el logger

### Ventajas del Patrón Adaptador

1. **Separación de responsabilidades**: Sentry no contamina el logger principal
2. **Fácil testing**: Se pueden mockear las integraciones
3. **Extensibilidad**: Nuevas integraciones sin modificar el core
4. **Graceful degradation**: El logger funciona aunque las integraciones fallen
5. **Configuración opcional**: Las integraciones son completamente opcionales

## Testing

```go
package logger_test

import (
    "testing"
    "your-project/logger"
)

func TestLogger(t *testing.T) {
    // Crear mock integration
    mockIntegration := logger.NewMockIntegration()
    
    config := &logger.LoggerConfig{
        Level:        logger.InfoLevel,
        Console:      false, // Deshabilitar consola para tests
        Integrations: []logger.LogIntegration{mockIntegration},
    }
    
    log, err := logger.NewLogger(config)
    if err != nil {
        t.Fatalf("Failed to create logger: %v", err)
    }
    defer log.Close()
    
    // Test logging
    log.Info("test_event", "test message", 0, map[string]interface{}{
        "test": true,
    })
    
    // Verificar que la integración recibió el evento
    logs := mockIntegration.GetLogs()
    if len(logs) != 1 {
        t.Errorf("Expected 1 log, got %d", len(logs))
    }
}
```

## Variables de Entorno

| Variable | Descripción | Default |
|----------|-------------|---------|
| `DEBUG` | Habilita modo debug | `false` |
| `SENTRY_DSN` | DSN de Sentry | - |
| `SENTRY_ENVIRONMENT` | Ambiente de Sentry | `development` |
| `SENTRY_RELEASE` | Release de Sentry | `v1.0.0` |

## Estructura del Proyecto

```
logger/
├── logger.go          # Core del logger e interfaz principal
├── sentry_adapter.go  # Adaptador para Sentry
├── integrations.go    # Integraciones predefinidas
├── utils.go          # Funciones de utilidad y logger global
├── example.go        # Ejemplos de uso
├── logger_test.go    # Tests
└── README.md         # Esta documentación
```

## Contribuir

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## Licencia

Este proyecto está bajo la Licencia MIT. Ver el archivo `LICENSE` para más detalles.
