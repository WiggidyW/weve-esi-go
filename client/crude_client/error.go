package crude_client

import "fmt"

// type InputError struct{ error }
// type EsiError struct{ error }

type CacheCreateError struct{ error }
type CacheGetError struct{ error }
type CacheSetError struct{ error }
type RequestParamsError struct{ error }
type HttpError struct{ error }
type MalformedResponse struct{ error }

type StatusError struct {
	Code int
	Text string
}

func (e StatusError) Error() string {
	return fmt.Sprintf("ESI Server returned Response Code %s", e.Text)
}

func newStatusError(
	code int,
	text string,
) StatusError {
	return StatusError{
		Code: code,
		Text: text,
	}
}

// type ServerError struct{ error }
