package logger

import "go.uber.org/zap"

func New() (*zap.Logger, error) {
	log, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}

	return log, nil
}
