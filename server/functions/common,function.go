package functions

import (
	"io"

	"github.com/h2non/filetype"
)

func Contains(validExtensions []string, ext string) bool {
	for _, v := range validExtensions {
		if v == ext {
			return true
		}
	}
	return false
}

func IsValidFileType(file io.Reader) bool {
	head := make([]byte, 261)
	_, err := file.Read(head)
	if err != nil {
		return false
	}
	kind, err := filetype.Match(head)
	if err != nil {
		return false
	}
	return kind.MIME.Type == "image" || kind.MIME.Type == "application/pdf"
}

// Simulate virus scanning (this function is a placeholder for real virus scanning logic)
func ScanForVirus(file io.Reader) error {
	// Placeholder for real virus scanning logic.
	// You could integrate with ClamAV or a cloud service for real virus checking.
	return nil
}
