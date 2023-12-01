package options

import "time"

type TimeoutOpt struct {
	timeout time.Duration
}

func (p TimeoutOpt) Detail() (string, []string) {
	return "Timeout", []string{
		p.timeout.String(),
	}
}

func (p TimeoutOpt) Handle(opts *ClientOptions) error {
	opts.Timeout = p.timeout
	return nil
}

func SetTimeoutOpt(timeout time.Duration) TimeoutOpt {
	return TimeoutOpt{
		timeout: timeout,
	}
}
