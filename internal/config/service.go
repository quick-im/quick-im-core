package config

import (
	"github.com/quick-im/quick-im-core/internal/logger"
	"github.com/quick-im/quick-im-core/internal/logger/innerzap"
	"go.uber.org/zap/zapcore"
)

type Option func(*serviceConfigImpl)

type ServiceConfig interface {
	GetIp() string
	GetPort() uint16
	GetOpenTracing() bool
	GetJeagerAgentHostPort() string
	GetUseConsulRegistry() bool
	GetConsulServers() []string
	GetNatsServers() []string
	GetNatsEnableJetstream() bool
	GetLogger() logger.Logger
}

type serviceConfigImpl struct {
	ip                  string
	port                uint16
	openTracing         bool
	trackAgentHostPort  string
	useConsulRegistry   bool
	consulServers       []string
	natsServers         []string
	natsEnableJetstream bool
	logger              logger.Logger
	loggerLevel         int8
	logPath             string
}

func NewServer(SERVER_NAME string, opts ...Option) ServiceConfig {
	ser := &serviceConfigImpl{
		consulServers:       make([]string, 0),
		natsServers:         make([]string, 0),
		natsEnableJetstream: true,
		loggerLevel:         int8(zapcore.DebugLevel),
		logPath:             "logs",
	}
	for i := range opts {
		opts[i](ser)
	}
	if ser.logger == nil {
		ser.logger = innerzap.NewZapLoggerAdapter(
			innerzap.NewLoggerWithOpt(
				innerzap.WithLogLevel(zapcore.Level(ser.loggerLevel)),
				innerzap.WithServiceName(SERVER_NAME),
				innerzap.WithLogPath(ser.logPath),
			).NewLogger(),
		)
	}
	return ser
}

// serviceConfig实现serviceConfigInterface接口
func WithIp(ip string) Option {
	return func(rs *serviceConfigImpl) {
		rs.ip = ip
	}
}

func WithPort(port uint16) Option {
	return func(rs *serviceConfigImpl) {
		rs.port = port
	}
}

func WithOpenTracing(disable bool) Option {
	return func(rs *serviceConfigImpl) {
		rs.openTracing = disable
	}
}

func WithJeagerAgentHostPort(trackAgentHostPort string) Option {
	return func(rs *serviceConfigImpl) {
		rs.trackAgentHostPort = trackAgentHostPort
	}
}

func WithUseConsulRegistry(useConsulRegistry bool) Option {
	return func(rs *serviceConfigImpl) {
		rs.useConsulRegistry = useConsulRegistry
	}
}

func WithConsulServer(consulServer string) Option {
	return func(rs *serviceConfigImpl) {
		if rs.consulServers == nil {
			rs.consulServers = make([]string, 0)
		}
		rs.consulServers = append(rs.consulServers, consulServer)
	}
}

func WithConsulServers(consulServers ...string) Option {
	return func(rs *serviceConfigImpl) {
		rs.consulServers = consulServers
	}
}

func WithNatsServer(natsServer string) Option {
	return func(rs *serviceConfigImpl) {
		if rs.natsServers == nil {
			rs.natsServers = make([]string, 0)
		}
		rs.natsServers = append(rs.natsServers, natsServer)
	}
}

func WithNatsServers(natsServers ...string) Option {
	return func(rs *serviceConfigImpl) {
		rs.natsServers = natsServers
	}
}

func WithNatsDisableJetstream() Option {
	return func(rs *serviceConfigImpl) {
		rs.natsEnableJetstream = false
	}
}

func WithLoggerLevel(level int8) Option {
	return func(rs *serviceConfigImpl) {
		rs.loggerLevel = level
	}
}

func WithLogrPath(path string) Option {
	return func(rs *serviceConfigImpl) {
		rs.logPath = path
	}
}

func WithLogger(logger logger.Logger) Option {
	return func(rs *serviceConfigImpl) {
		rs.logger = logger
	}
}

func (s *serviceConfigImpl) GetIp() string {
	return s.ip
}

func (s *serviceConfigImpl) GetPort() uint16 {
	return s.port
}

func (s *serviceConfigImpl) GetOpenTracing() bool {
	return s.openTracing
}

func (s *serviceConfigImpl) GetJeagerAgentHostPort() string {
	return s.trackAgentHostPort
}

func (s *serviceConfigImpl) GetUseConsulRegistry() bool {
	return s.useConsulRegistry
}

func (s *serviceConfigImpl) GetConsulServers() []string {
	return s.consulServers
}

func (s *serviceConfigImpl) GetNatsServers() []string {
	return s.natsServers
}

func (s *serviceConfigImpl) GetNatsEnableJetstream() bool {
	return s.natsEnableJetstream
}

func (s *serviceConfigImpl) GetLogger() logger.Logger {
	return s.logger
}
