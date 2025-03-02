package config

// Queue конфигурация очереди сообщений.
type Queue struct {
	// MaxSize максимальный размер очереди.
	MaxSize int `mapstructure:"max_size"`

	// MaxConsumers максимальное количество подписчиков.
	MaxConsumers int `mapstructure:"max_consumers"`
}
