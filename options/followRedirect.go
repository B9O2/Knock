package options

import "fmt"

type FollowRedirectsOpt struct {
	enabled bool
	max     int
}

func (fr FollowRedirectsOpt) Detail() (string, []string) {
	return "FollowRedirects", []string{fmt.Sprint(fr.enabled), fmt.Sprint(fr.max)}
}

func (fr FollowRedirectsOpt) Handle(opts *ClientOptions) error {
	opts.FollowRedirects = fr.enabled
	opts.MaxRedirects = fr.max
	return nil
}

func SetFollowRedirects(enable bool, max int) FollowRedirectsOpt {
	return FollowRedirectsOpt{
		enabled: enable,
		max:     max,
	}
}
