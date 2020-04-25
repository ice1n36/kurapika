package handlers

import (
	"encoding/json"
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

func NewNewAppHandler(logger *zap.SugaredLogger, mobsf clients.MobSFHTTPClient) (*NewAppHandler, error) {
	logger.Infow("Executing NewNewAppHandler.")
	return &NewAppHandler{
		logger: logger,
		mobsf:  mobsf,
	}, nil
}

func (s *NewAppHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.logger.Infow("New App Event")

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

	s.logger.Infow("Downloading app",
		"req", sareq,
		"dir", dir,
	)

	// download app
	var app string

	if sareq.Os == "android" {
		app, err = s.downloadAPK(sareq.AppId, sareq.DeviceCodeName, dir)
		if err != nil {
			http.Error(w, "download app failed", http.StatusInternalServerError)
			s.logger.Errorw("error during downloading of app", "error", err.Error())
			return
		}
	} else if sareq.Os == "ios" {
		http.Error(w, "unsupported os", http.StatusBadRequest)
		return
	} else {
		http.Error(w, "unsupported os", http.StatusBadRequest)
		return
	}

	// upload app to mobsf
	hash, err := s.mobsf.Upload(app)
	if err != nil {
		http.Error(w, "Upload app failed", http.StatusInternalServerError)
		return
	}
	s.logger.Infow("Upload successful", "hash", hash)

	// TODO: spin up thread to wait for mobsf scan to finish, get the result and report (i.e. email)
	return
}

// downloads and returns the path to the apk downloaded
func (s *NewAppHandler) downloadAPK(appId string, deviceCodeName string, dir string) (string, error) {
	cmd := exec.Command("gplaycli", "-f", dir, "-dc", deviceCodeName, "-d", appId)

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	apkPath := appId + ".apk"
	return filepath.Join(dir, apkPath), nil
}
