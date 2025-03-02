package domain

import "errors"

var ErrBufferOverflow = errors.New("max buffer size reached")

// RingBuffer кольцевой буфер сообщений.
type RingBuffer struct {
	// beg начало буфера.
	beg int

	// end конец буфера.
	end int

	messages []Message
}

func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		messages: make([]Message, size),
	}
}

// Size возвращает размер буфера.
func (b *RingBuffer) Size() int {
	if b.beg <= b.end {
		return b.end - b.beg
	}

	return len(b.messages) - b.beg + b.end
}

// PushBack помещает сообщение в конец буфера.
func (b *RingBuffer) PushBack(msg Message) error {
	if b.Size()+1 == len(b.messages) {
		return ErrBufferOverflow
	}

	b.end++

	if b.end > len(b.messages) {
		b.end = 0
	}

	b.messages[b.end] = msg

	return nil
}

// SpaceLeft возвращает количество сообщений, которые можно положить в буфер до переполнения.
func (b *RingBuffer) SpaceLeft() int {
	return len(b.messages) - b.Size()
}
