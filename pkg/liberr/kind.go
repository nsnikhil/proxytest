package liberr

type Kind string

const (
	ValidationError  Kind = "validationError"
	InternalError    Kind = "internalError"
	RateLimitedError Kind = "rateLimitedError"
)
