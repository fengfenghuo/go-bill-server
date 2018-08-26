package rpc

type shutdownError struct{}

func (e *shutdownError) Error() string { return "server is shutting down" }

type InvalidParamsError struct {
	Msg string
}

func (e *InvalidParamsError) Error() string {
	if len(e.Msg) > 0 {
		return e.Msg
	} else {
		return "invalid parameter"
	}
}

type ErrSubscriptionNotFound struct{}

func (e *ErrSubscriptionNotFound) Error() string { return "subscription not found" }
