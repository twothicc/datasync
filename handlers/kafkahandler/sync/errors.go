package sync

import (
	"github.com/twothicc/common-go/errortype"
)

const pkg = "handlers/kafkahandler/sync"

//nolint:gomnd // error code
var (
	ErrConstructor = errortype.ErrorType{Code: 1, Pkg: pkg}
	ErrConsume     = errortype.ErrorType{Code: 2, Pkg: pkg}
	ErrUnmarshal   = errortype.ErrorType{Code: 3, Pkg: pkg}
)
