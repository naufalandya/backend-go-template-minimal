package imagekit

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type ImageKit struct {
	PublicKey  string
	PrivateKey string
	Endpoint   string
}

func New() *ImageKit {
	return &ImageKit{
		PublicKey:  os.Getenv("IMAGEKIT_PUBLIC_KEY"),
		PrivateKey: os.Getenv("IMAGEKIT_PRIVATE_KEY"),
		Endpoint:   os.Getenv("IMAGEKIT_URL_ENDPOINT"),
	}
}

type UploadResponse struct {
	Url string `json:"url"`
}

func (ik *ImageKit) UploadFile(fileBytes []byte, fileName string) (*UploadResponse, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Create "file" field
	fileField, err := writer.CreateFormField("file")
	if err != nil {
		return nil, err
	}

	encodedFile := base64.StdEncoding.EncodeToString(fileBytes)
	_, err = io.WriteString(fileField, "data:image/jpeg;base64,"+encodedFile)
	if err != nil {
		return nil, err
	}

	// Create "fileName" field
	_ = writer.WriteField("fileName", fileName)

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	// Prepare the request
	req, err := http.NewRequest("POST", "https://upload.imagekit.io/api/v1/files/upload", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.SetBasicAuth(ik.PublicKey, ik.PrivateKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed upload: %s", string(respBody))
	}

	var uploadResp UploadResponse
	err = json.Unmarshal(respBody, &uploadResp)
	if err != nil {
		return nil, err
	}

	return &uploadResp, nil
}
