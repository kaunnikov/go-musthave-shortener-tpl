package errs

import "fmt"

type DoubleError struct {
	ShortURL string
	Err      error
}

type TokenNotFoundInCookieError struct {
	Err error
}

func (d *DoubleError) Error() string {
	return fmt.Sprintf("[%s] %v", d.ShortURL, d.Err)
}
func (d *TokenNotFoundInCookieError) Error() string {
	return fmt.Sprintf("%v", d.Err)
}
