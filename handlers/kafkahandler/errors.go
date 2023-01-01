package kafkahandler

import (
	"github.com/twothicc/common-go/errortype"
)

const pkg = "handlers/kafkahandler"

//nolint:gomnd // error code
var (
	ErrConstructor       = errortype.ErrorType{Code: 1, Pkg: pkg}
	ErrUnmarshal         = errortype.ErrorType{Code: 2, Pkg: pkg}
	ErrIndex             = errortype.ErrorType{Code: 3, Pkg: pkg}
	ErrInvalidCtimestamp = errortype.ErrorType{Code: 4, Pkg: pkg}
	ErrUniqueId          = errortype.ErrorType{Code: 5, Pkg: pkg}
)
