package config

// Queue конфигурация очереди сообщений.
type Queue struct {
	// MaxSize максимальное количество сообщений, ожидающих отправки.
	MaxSize int `mapstructure:"max_size"`

	// MaxConsumers максимальное количество подписчиков.
	MaxConsumers int `mapstructure:"max_consumers"`
}
