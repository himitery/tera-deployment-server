package config

type Config struct {
	Profile  string          `yaml:"profile"`
	Services []ServiceConfig `yaml:"services"`
	Kafka    KafkaConfig     `yaml:"kafka"`
	Logging  LoggingConfig   `yaml:"logging"`
}

type ServiceConfig struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Depends []struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	} `yaml:"depends"`
}

type KafkaConfig struct {
	BootstrapServers []KafkaBootstrapServerConfig `yaml:"bootstrap_servers"`
	Protocol         string                       `yaml:"protocol"`
	Sasl             KafkaSaslConfig              `yaml:"sasl"`
	Topic            string                       `yaml:"topic"`
}

type KafkaSaslConfig struct {
	Mechanism string `yaml:"mechanism"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
}

type KafkaBootstrapServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type LoggingConfig struct {
	Level string `yaml:"level"`
}
