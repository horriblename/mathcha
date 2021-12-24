package latex

const ERROR_TOLERANCE = 1

type ErrCode int

const (
	ERR_TOKEN           = iota
	ERR_MISSING_CLOSE   // A Closing expression is missing e.g. "}" or "\right"
	ERR_UNMATCHED_CLOSE // Unmatched Closing expression e.g. "}" or "\right"
)

var errType = [...]string{
	ERR_TOKEN:           "ERR_TOKEN",
	ERR_MISSING_CLOSE:   "ERR_MISSING_CLOSE",
	ERR_UNMATCHED_CLOSE: "ERR_UNMATCHED_CLOSING",
}

func (e ErrCode) String() string { return errType[e] }

type ErrorHandler struct {
	errorList []ParseErr
}

type ParseErr struct {
	errType ErrCode
	desc    string
}

func (eh *ErrorHandler) AddErr(e ErrCode, desc string) {
	eh.errorList = append(eh.errorList, ParseErr{errType: e, desc: desc})
	if eh.Errors() >= ERROR_TOLERANCE {
		// FIXME
		println("Last encountered error: ", e.String())
		println("details: ", desc)
		panic("Too many errors encountered!")
	}
}

func (eh *ErrorHandler) Errors() int { return len(eh.errorList) }

func (eh *ErrorHandler) Trace() {
}
