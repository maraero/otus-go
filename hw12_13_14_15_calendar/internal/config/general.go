package config

import "errors"

func validateConfigLogger(c Logger) error {
	if len(c.OutputPaths) == 0 {
		return errors.New(ErrMissingOutputPaths)
	}

	if len(c.ErrorOutputPaths) == 0 {
		return errors.New(ErrMissingErrOutputPaths)
	}

	return nil
}
