package utils

import "strings"

func GetFileType(fileName string) string {
	if strings.HasSuffix(fileName, ".pdf") {
		return "pdf"
	} else if strings.HasSuffix(fileName, ".docx") {
		return "docx"
	}
	return "unknown"
}
