package crude_client

import (
	"fmt"
	"io"
	"net/http"
)

// type InputError struct{ error }
// type EsiError struct{ error }

type CacheCreateError struct{ error }
type CacheGetError struct{ error }
type CacheSetError struct{ error }
type RequestParamsError struct{ error }
type HttpError struct{ error }
type MalformedResponse struct{ error }

type StatusError struct {
	Url      string
	CodeText string
	EsiText  string
}

func (e StatusError) Error() string {
	errstr := fmt.Sprintf(
		"ESI Server Request '%s' returned Response Code '%s'",
		e.Url,
		e.CodeText,
	)
	if e.EsiText == "" {
		errstr += " with no error message"
	} else {
		errstr += fmt.Sprintf(
			" with error message '%s'",
			e.EsiText,
		)
	}
	return errstr
}

type EsiError struct {
	Error string `json:"error"`
}

func newStatusError(rep *http.Response) StatusError {
	var body_str string
	body_bytes, err := io.ReadAll(rep.Body)
	if err != nil {
		body_str = ""
	} else {
		body_str = string(body_bytes)
	}
	return StatusError{
		Url:      rep.Request.URL.String(),
		CodeText: rep.Status,
		EsiText:  body_str,
	}
}

// type ServerError struct{ error }
