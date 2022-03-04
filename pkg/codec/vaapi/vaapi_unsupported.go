// +build !dragonfly,!freebsd,!linux,!netbsd,!openbsd,!solaris

package vaapi

// // Dummy CGO import to avoid `C source files not allowed when not using cgo or SWIG`
import "C"

import (
	"errors"

	"github.com/cpdevs/mediadevices/pkg/codec"
	"github.com/cpdevs/mediadevices/pkg/io/video"
	"github.com/cpdevs/mediadevices/pkg/prop"
)

var errNotSupported = errors.New("vaapi is not supported on this environment")

func newVP8Encoder(r video.Reader, p prop.Media, params ParamsVP8) (codec.ReadCloser, error) {
	return nil, errNotSupported
}

func newVP9Encoder(r video.Reader, p prop.Media, params ParamsVP9) (codec.ReadCloser, error) {
	return nil, errNotSupported
}
