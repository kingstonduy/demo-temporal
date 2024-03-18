package bootstrap

import "github.com/lengocson131002/go-clean-core/config"

type ServerConfig struct {
	Name       string
	AppVersion string
	HttpPort   int
	GrpcPort   int
	BaseURI    string
}

func GetServerConfig(cfg config.Configure) *ServerConfig {
	name := cfg.GetString("SERVER_NAME")
	version := cfg.GetString("SERVER_VERSION")
	httpPort := cfg.GetInt("SERVER_HTTP_PORT")
	grpcPort := cfg.GetInt("SERVER_GRPC_PORT")
	baseUrl := cfg.GetString("SERVER_BASE_URL")

	return &ServerConfig{
		Name:       name,
		AppVersion: version,
		HttpPort:   httpPort,
		GrpcPort:   grpcPort,
		BaseURI:    baseUrl,
	}
}
