package model

type Config struct {
	Common CommonConf
	Mongo  MongoConf
	Redis  RedisConf
	OpenAi OpenAiConf
}

type CommonConf struct {
	HttpAddress string
	ServiceName string
}

type RedisConf struct {
	Address     string
	MaxPoolSize int
}

type MongoConf struct {
	Address     string
	MaxPoolSize int
}

type OpenAiConf struct {
	Secret string
	Url    string
}
