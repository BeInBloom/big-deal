package models

type Option[T any] struct {
	value T
	set   bool
}

func Some[T any](value T) Option[T] {
	return Option[T]{
		value: value,
		set:   true,
	}
}

func None[T any]() Option[T] {
	return Option[T]{}
}

func (o Option[T]) Get() (T, bool) {
	return o.value, o.set
}

func (o Option[T]) IsSet() bool {
	return o.set
}
