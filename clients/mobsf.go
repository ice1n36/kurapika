package clients

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/config"
	"go.uber.org/zap"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

// MobSFHTTPClient is the interface for all things that interface with mobsf
type MobSFHTTPClient interface {
	Upload(string) (string, error)
	Scan(string, string, string) error
}

type mobSFHTTPClientImpl struct {
	logger *zap.SugaredLogger
	host   string
	port   int
	apiKey string
	client *http.Client
	config config.Provider
}

type uploadRespBody struct {
	ScanType string `json:"scan_type"`
	Hash     string `json:"hash"`
	Filename string `json:"file_name"`
}

// NewMobSFHTTPClient creates a new MobSFHTTPClient based on configuration
func NewMobSFHTTPClient(logger *zap.SugaredLogger, config config.Provider) MobSFHTTPClient {
	return &mobSFHTTPClientImpl{
		logger: logger,
		host:   config.Get("mobsf.host").String(),
		port:   config.Get("mobsf.port").Value().(int),
		apiKey: config.Get("mobsf.apikey").String(),
		client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

// Upload sends app (apk/ipa) to mobsf via /api/v1/upload
func (m *mobSFHTTPClientImpl) Upload(appPath string) (string, error) {
	m.logger.Infow("Uploading to mobsf", "file", appPath)

	extraParams := map[string]string{}
	url := fmt.Sprintf("http://%s:%d/api/v1/upload", m.host, m.port)
	request, err := m.newfileUploadRequest(url, extraParams, "file", appPath)
	if err != nil {
		m.logger.Errorw("error creating new upload request", "error", err.Error())
		return "", err
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		m.logger.Errorw("error uploading", "error", err.Error())
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		m.logger.Errorw("bad response code", "code", resp.StatusCode)
		return "", errors.New("Bad response code")
	}

	respBody := uploadRespBody{}
	respBodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		m.logger.Errorw("error reading resp body", "error", err.Error())
		return "", err
	}

	err = json.Unmarshal(respBodyData, &respBody)
	if err != nil {
		m.logger.Errorw("error unmarshalling", "error", err.Error())
		return "", err
	}
	return respBody.Hash, nil
}

func (m *mobSFHTTPClientImpl) Scan(appId string, appType string, hash string) error {
	return nil
}

// Creates a new file upload http request with optional extra params
func (m *mobSFHTTPClientImpl) newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, fi.Name())
	if err != nil {
		return nil, err
	}
	part.Write(fileContents)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", uri, body)
	request.Header.Add("Content-Type", writer.FormDataContentType())
	request.Header.Add("Authorization", m.apiKey)
	return request, nil
}
