package client

import (
	"fmt"
	"reflect"
	"time"
)

func u64FirstDigit(v uint64) uint64 {
	for v >= 10 {
		v = v / 10
	}
	return v
}

func rfc3339ToUnix(s string) (uint64, error) {
	datetime, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return 0, err
	}
	timestamp := datetime.Unix()
	if timestamp < 0 {
		return 0, fmt.Errorf(
			"unix timestamp was negative: %d",
			timestamp,
		)
	}
	return uint64(timestamp), nil
}

// func getJsonError[T interface{}](
// 	key string,
// 	value interface{},
// ) error {
// 	return MalformedJson{fmt.Errorf(
// 		"json value for key: %s was not of type: %T, was: %s",
// 		key,
// 		*new(T),
// 		reflect.TypeOf(value).String(),
// 	)}
// }

func getJson[T interface{}](
	json map[string]interface{},
	key string,
) (*T, error) {
	value, ok := json[key]
	if ok && value != nil {
		switch t := value.(type) {
		case T:
			return &t, nil
		default:
			return nil, MalformedJson{fmt.Errorf(
				"%s: %s was not of type: %T, was: %s",
				"json value for key",
				key,
				*new(T),
				reflect.TypeOf(value).String(),
			)}
		}
	} else {
		return nil, MalformedJson{fmt.Errorf(
			"json did not contain key: %s",
			key,
		)}
	}
}

func getValueOrPanic[T interface{}](
	json map[string]interface{},
	key string,
) T {
	value, err := getJson[T](json, key)
	if err != nil {
		panic(err)
	}
	return *value
}

func getValueOrDefault[T interface{}](
	json map[string]interface{},
	key string,
	default_value T,
) T {
	value, err := getJson[T](json, key)
	if err != nil {
		return default_value
	}
	return *value
}

func getTimestamp(
	json map[string]interface{},
	key string,
) (uint64, error) {
	value, err := getJson[string](json, key)
	if err != nil {
		return 0, err
	}
	timestamp, err := rfc3339ToUnix(*value)
	if err != nil {
		return 0, MalformedJson{fmt.Errorf(
			"%s: %s was not valid RFC3339, was: %s",
			"json value for key",
			key,
			*value,
		)}
	}
	return timestamp, nil
}

func getTimestampOrPanic(
	json map[string]interface{},
	key string,
) uint64 {
	value, err := getTimestamp(json, key)
	if err != nil {
		panic(err)
	}
	return value
}

// func getTimestampOrDefault(
// 	json map[string]interface{},
// 	key string,
// 	default_value uint64,
// ) uint64 {
// 	value, err := getTimestamp(json, key)
// 	if err != nil {
// 		return default_value
// 	}
// 	return value
// }
