package errs

import "fmt"

type DoubleError struct {
	ShortURL string
	Err      error
}

// Error добавляет поддержку интерфейса error для типа LabelError.
func (d *DoubleError) Error() string {
	return fmt.Sprintf("[%s] %v", d.ShortURL, d.Err)
}
