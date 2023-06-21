package env

import (
	"fmt"
)

type EnvMissingError struct {
	Key string
}

func (e EnvMissingError) Error() string {
	return fmt.Sprintf("Environment Variable %s is missing", e.Key)
}

func newEnvMissingError(
	key string,
) EnvMissingError {
	return EnvMissingError{
		Key: key,
	}
}

type EnvInvalidError struct {
	Key  string
	Val  string
	Type string
}

func (e EnvInvalidError) Error() string {
	return fmt.Sprintf(
		"Environment Variable %s is invalid: %s is not a valid %s",
		e.Key,
		e.Val,
		e.Type,
	)
}

func newEnvInvalidError[T interface{}](
	key string,
	val string,
) EnvInvalidError {
	return EnvInvalidError{
		Key:  key,
		Val:  val,
		Type: fmt.Sprintf("%T", *new(T)),
	}
}
