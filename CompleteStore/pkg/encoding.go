package pkg

import "io"

type Decoder interface {
	Decoder(io.Reader, any) error
}
