package contract

type HttpResponseError struct {
	Code int
	Msg string
	Err error
	ResponseBody []byte
}

func (e *HttpResponseError) Error() string {
	return e.Msg
}

func (e *HttpResponseError) Unwrap() error {
	return e.Err
}