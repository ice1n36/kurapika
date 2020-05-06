package handlers

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/ice1n36/kurapika/clients"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

type NewAppHandler struct {
	logger *zap.SugaredLogger
	mobsf  clients.MobSFHTTPClient
}

type NewAppRequest struct {
	AppId          string `json:"app_id"`
	AppVersion     string `json:"app_version"`
	DeviceCodeName string `json:"device_codename"`
	Os             string `json:"os"`
}

type NewAppResponse struct {
	Hash       string `json:"hash"`
	AppId      string `json:"app_id"`
	AppVersion string `json:"app_version"`
	KId        string `json:"k_id"`
	Status     string `json:"status"`
}

const (
	StatusDownloading = "Downloading"
	StatusUploading   = "Uploading"
	StatusScanning    = "Scanning"
	StatusComplete    = "Complete"
	StatusFailure     = "Failure"
)

func NewNewAppHandler(logger *zap.SugaredLogger, mobsf clients.MobSFHTTPClient) (*NewAppHandler, error) {
	logger.Infow("Executing NewNewAppHandler.")
	return &NewAppHandler{
		logger: logger,
		mobsf:  mobsf,
	}, nil
}

func (s *NewAppHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// generate a kurapika id
	kId := uuid.New().String()
	s.logger.Infow("New App Event", "kid", kId)

	sareq := NewAppRequest{}

	err := json.NewDecoder(req.Body).Decode(&sareq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dir, err := ioutil.TempDir("", "*-apps")
	if err != nil {
		http.Error(w, "fs failure", http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(dir)

	response := NewAppResponse{
		AppId:      sareq.AppId,
		AppVersion: sareq.AppVersion,
		KId:        kId,
		Status:     StatusDownloading,
	}

	s.logger.Infow("Downloading app",
		"req", sareq,
		"dir", dir,
		"kid", kId,
	)

	// download app
	var app string

	if sareq.Os == "android" {
		app, err = s.downloadAPK(sareq.AppId, sareq.DeviceCodeName, dir)
		if err != nil {
			http.Error(w, "download app failed", http.StatusInternalServerError)
			s.logger.Errorw("error during downloading of app",
				"kid", kId,
				"error", err.Error(),
			)
			return
		}
	} else if sareq.Os == "ios" {
		http.Error(w, "unsupported os", http.StatusBadRequest)
		return
	} else {
		http.Error(w, "unsupported os", http.StatusBadRequest)
		return
	}

	response.Status = StatusUploading

	// upload app to mobsf
	appPath := filepath.Join(dir, app)
	hash, err := s.mobsf.Upload(appPath)
	if err != nil {
		http.Error(w, "Upload app failed", http.StatusInternalServerError)
		return
	}
	s.logger.Infow("Upload successful", "hash", hash)
	response.Hash = hash

	response.Status = StatusScanning
	go func() {
		err := s.mobsf.Scan(app, "apk", hash)
		if err != nil {
			s.logger.Errorw("Scan failed",
				"kid", kId,
				"error", err.Error(),
			)
			// TODO: alert something that scan failed
		}

		// TODO: get the result and report (i.e. email)
	}()

	responseJson, err := json.Marshal(response)
	if err != nil {
		s.logger.Errorw("Error marshalling response",
			"kid", kId,
			"error", err.Error(),
		)
		http.Error(w, "Response creation failure", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJson)
	return
}

// downloads and returns the path to the apk downloaded
func (s *NewAppHandler) downloadAPK(appId string, deviceCodeName string, dir string) (string, error) {
	cmd := exec.Command("gplaycli", "-f", dir, "-av", "-dc", deviceCodeName, "-d", appId)

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}
	if len(files) != 1 {
		return "", errors.New("more than one file in app download directory")
	}

	return files[0].Name(), nil
}
