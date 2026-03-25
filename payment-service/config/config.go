package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
)

type Config struct {
	Port           int    `env:"PORT" envDefault:"8081"`
	GinMode        string `env:"GIN_MODE" envDefault:"debug"`
	LogLevel       string `env:"LOG_LEVEL" envDefault:"debug"`
	ServerCertFile string `env:"SERVER_CERT_FILE" envDefault:""`
	ServerKeyFile  string `env:"SERVER_KEY_FILE" envDefault:""`
	CACertFile     string `env:"CA_CERT_FILE" envDefault:""`
}

func InitOpentel() {
	tp := trace.NewTracerProvider()
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})
}

func Load() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)
	return cfg, err
}

func LoadTLSConfig(cfg *Config) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(
		cfg.ServerCertFile,
		cfg.ServerKeyFile,
	)
	if err != nil {
		return nil, fmt.Errorf("load server cert: %w", err)
	}

	caCert, err := os.ReadFile(cfg.CACertFile)
	if err != nil {
		return nil, fmt.Errorf("read CA cert: %w", err)
	}

	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		return nil, fmt.Errorf("failed to append CA cert")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},

		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs:  caCertPool,

		MinVersion: tls.VersionTLS12,
	}

	return tlsConfig, nil
}
