package config

type Config struct {
	Profile  string          `json:"profile"`
	Services []ServiceConfig `yaml:"services"`
	Argocd   ArgocdConfig    `yaml:"argocd"`
	Kafka    KafkaConfig     `yaml:"kafka"`
	Logging  LoggingConfig   `yaml:"logging"`
}

type ServiceConfig struct {
	Name    string                `yaml:"name"`
	Version string                `yaml:"version"`
	Depends []ServiceDependConfig `yaml:"depends"`
}

type ServiceDependConfig struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

type ArgocdConfig struct {
	URL        string               `yaml:"url"`
	Token      string               `yaml:"token"`
	Repository string               `yaml:"repository"`
	Metadata   ArgocdMetadataConfig `yaml:"metadata"`
}

type ArgocdMetadataConfig struct {
	Namespace string `yaml:"namespace"`
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
