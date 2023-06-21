package client

type Result[T interface{}] struct {
	Ok  *T
	Err error
}

func (r Result[T]) Unwrap() (T, error) {
	return *r.Ok, r.Err
}

func ResultOk[T interface{}](ok T) Result[T] {
	return Result[T]{
		Ok:  &ok,
		Err: nil,
	}
}

func ResultErr[T interface{}](err error) Result[T] {
	return Result[T]{
		Ok:  nil,
		Err: err,
	}
}

func ResultNull[T interface{}]() Result[T] {
	return Result[T]{
		Ok:  nil,
		Err: nil,
	}
}

func (r Result[T]) IsNull() bool {
	return r.Ok == nil && r.Err == nil
}
