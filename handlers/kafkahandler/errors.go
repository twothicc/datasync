package kafkahandler

import (
	"github.com/twothicc/common-go/errortype"
)

const pkg = "handlers/kafkahandler"

//nolint:gomnd // error code
var (
	ErrConstructor = errortype.ErrorType{Code: 1, Pkg: pkg}
)
