package elastic

import (
	"github.com/twothicc/common-go/errortype"
)

const pkg = "infra/elasticsearch"

//nolint:gomnd // error code
var (
	ErrConstructor = errortype.ErrorType{Code: 1, Pkg: pkg}
)
