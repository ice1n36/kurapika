package handlers

import (
	"go.uber.org/zap"
	"net/http"
)

type ScanApkHandler struct {
	logger *zap.SugaredLogger
}

func NewScanApkHandler(logger *zap.SugaredLogger) (*ScanApkHandler, error) {
	logger.Infow("Executing NewScanApkHandler.")
	return &ScanApkHandler{
		logger: logger,
	}, nil
}

func (s *ScanApkHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	s.logger.Infow("ScanApkHandler")
}
