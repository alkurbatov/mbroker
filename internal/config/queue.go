package config

// Queue конфигурация очереди сообщений.
type Queue struct {
	// MaxSize максимальный размер очереди.
	MaxSize int64 `mapstructure:"max_size"`

	// MaxConsumers максимальное количество подписчиков.
	MaxConsumers int64 `mapstructure:"max_consumers"`
}
