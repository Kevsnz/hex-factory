package utils

import "fmt"

type RingBuffer[T any] struct {
	data []T
	head int
	tail int
	size int
}

func NewRingBuffer[T any](size int) RingBuffer[T] {
	return RingBuffer[T]{
		data: make([]T, size),
		head: 0,
		tail: 0,
		size: 0,
	}
}

func (rb *RingBuffer[T]) Push(t T) error {
	if rb.size == len(rb.data) {
		return fmt.Errorf("ring buffer overflow")
	}
	rb.data[rb.tail] = t
	rb.tail = (rb.tail + 1) % len(rb.data)
	rb.size++
	return nil
}

func (rb *RingBuffer[T]) Pop() (T, bool) {
	var tNull T
	if rb.size == 0 {
		return tNull, false
	}
	t := rb.data[rb.head]
	rb.data[rb.head] = tNull
	rb.head = (rb.head + 1) % len(rb.data)
	rb.size--
	return t, true
}

func (rb *RingBuffer[T]) PopLast() (T, bool) {
	var tNull T
	if rb.size == 0 {
		return tNull, false
	}
	rb.tail = (rb.tail + len(rb.data) - 1) % len(rb.data)
	t := rb.data[rb.tail]
	rb.data[rb.tail] = tNull
	rb.size--
	return t, true
}

func (rb *RingBuffer[T]) PopAt(i int) (T, bool) {
	l := len(rb.data)
	var tNull T
	if i < 0 || i >= rb.size {
		return tNull, false
	}

	item := rb.data[rb.head]
	rb.data[rb.head] = tNull
	for j := 1; j <= i; j++ {
		rb.data[(rb.head+j)%l], item = item, rb.data[(rb.head+j)%l]
	}
	rb.head = (rb.head + 1) % l
	rb.size--
	return item, true
}

func (rb *RingBuffer[T]) Clear() {
	var t T
	for rb.size > 0 {
		rb.data[rb.head] = t
		rb.head = (rb.head + 1) % len(rb.data)
		rb.size--
	}
}

func (rb *RingBuffer[T]) Resize(size int) error {
	if size < rb.size {
		return fmt.Errorf("cannot make ring buffer smaller than number of currently occupied elements")
	}

	newData := make([]T, size)
	i := 0
	for rb.size > 0 {
		newData[i], _ = rb.Pop()
		i++
	}
	rb.head = 0
	rb.tail = i
	return nil
}

func (rb *RingBuffer[T]) Len() int {
	return rb.size
}

func (rb *RingBuffer[T]) Peek(i int) T {
	if i >= rb.size {
		panic("index out of bounds")
	}

	return rb.data[(rb.head+i)%len(rb.data)]
}

func (rb *RingBuffer[T]) PeekFirst() (T, bool) {
	var tNull T
	if rb.size == 0 {
		return tNull, false
	}
	t := rb.data[rb.head]
	return t, true
}

func (rb *RingBuffer[T]) PeekLast() (T, bool) {
	var tNull T
	if rb.size == 0 {
		return tNull, false
	}
	t := rb.data[(rb.tail+len(rb.data)-1)%len(rb.data)]
	return t, true
}

func (rb *RingBuffer[T]) Insert(i int, t T) error {
	if rb.size == len(rb.data) {
		return fmt.Errorf("ring buffer overflow")
	}
	if i < 0 || i > rb.size {
		return fmt.Errorf("index out of bounds")
	}

	for j := rb.head + rb.size - 1; j >= rb.head+i; j-- {
		rb.data[(j+1)%len(rb.data)] = rb.data[j%len(rb.data)]
	}
	rb.data[(rb.head+i)%len(rb.data)] = t

	rb.tail = (rb.tail + 1) % len(rb.data)
	rb.size++
	return nil
}
