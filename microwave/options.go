package microwave

import (
	"github.com/pkg/errors"
	"github.com/sqrt-7/microwave/log"
	"github.com/sqrt-7/microwave/tools"
)

type Option interface {
	apply(*Microwave) error
}

type optionFn func(*Microwave) error

func (o optionFn) apply(l *Microwave) error {
	return o(l)
}

func CustomLogger(customLogger log.Logger) Option {
	return optionFn(func(input *Microwave) error {
		input.logger = customLogger
		return nil
	})
}

func RequireEnvs(envs ...string) Option {
	return optionFn(func(input *Microwave) error {
		if len(envs) > 0 {
			res, err := tools.EnvLookup(envs...)
			if err != nil {
				return errors.Wrap(err, MsgBootError)
			}
			input.envs = res
		}

		return nil
	})
}
