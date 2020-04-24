package clients

import (
	"go.uber.org/zap"
)

type MobSFHTTPClient interface {
	Upload(string) error
}

type mobSFHTTPClientImpl struct {
	logger *zap.SugaredLogger
	apiKey string
}

func NewMobSFHTTPClient(logger *zap.SugaredLogger) MobSFHTTPClient {
	return &mobSFHTTPClientImpl{
		logger: logger,
		apiKey: "none", // TODO: pull apiKey from configuration file
	}
}

func (m *mobSFHTTPClientImpl) Upload(appPath string) error {
	// TODO: impl
	m.logger.Infow("Uploading to mobsf", "file", appPath)
	return nil
}
