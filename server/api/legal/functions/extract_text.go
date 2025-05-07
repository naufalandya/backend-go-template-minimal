package functions

import (
	"io"
	"os"
	"sort"
	"strings"
	"unicode"

	"baliance.com/gooxml/document"

	"rsc.io/pdf"
)

type pdfChar struct {
	X float64
	Y float64
	S string
}

func ExtractTextFromPDF(file io.Reader) (string, error) {
	tempFile, err := os.CreateTemp("", "upload-*.pdf")
	if err != nil {
		return "", err
	}
	defer os.Remove(tempFile.Name())

	_, err = io.Copy(tempFile, file)
	if err != nil {
		return "", err
	}
	tempFile.Close()

	doc, err := pdf.Open(tempFile.Name())
	if err != nil {
		return "", err
	}

	var textBuilder strings.Builder

	for i := 1; i <= doc.NumPage() && i <= 15; i++ {
		page := doc.Page(i)
		if page.V.IsNull() {
			continue
		}
		content := page.Content()

		var chars []pdfChar
		for _, txt := range content.Text {
			clean := sanitizeText(txt.S)
			if clean != "" {
				chars = append(chars, pdfChar{X: txt.X, Y: txt.Y, S: clean})
			}
		}

		lines := groupLines(chars, 2.0)

		var yKeys []float64
		for y := range lines {
			yKeys = append(yKeys, y)
		}
		sort.Sort(sort.Reverse(sort.Float64Slice(yKeys)))

		for _, y := range yKeys {
			line := lines[y]
			sort.Slice(line, func(i, j int) bool {
				return line[i].X < line[j].X
			})

			var lineBuilder strings.Builder
			prevX := -1000.0
			for _, ch := range line {
				if ch.X-prevX > 2.5 {
					lineBuilder.WriteString(" ")
				}
				lineBuilder.WriteString(ch.S)
				prevX = ch.X
			}
			textBuilder.WriteString(lineBuilder.String())
			textBuilder.WriteString(" ")
		}
	}

	rawText := textBuilder.String()
	cleaned := strings.Join(strings.Fields(rawText), " ")
	return cleaned, nil
}

func sanitizeText(input string) string {
	var b strings.Builder
	for _, r := range input {
		if unicode.IsPrint(r) && (unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsPunct(r) || unicode.IsSpace(r)) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func groupLines(chars []pdfChar, tolerance float64) map[float64][]pdfChar {
	lines := make(map[float64][]pdfChar)

	for _, c := range chars {
		found := false
		for y := range lines {
			if abs(c.Y-y) <= tolerance {
				lines[y] = append(lines[y], c)
				found = true
				break
			}
		}
		if !found {
			lines[c.Y] = []pdfChar{c}
		}
	}
	return lines
}

func abs(f float64) float64 {
	if f < 0 {
		return -f
	}
	return f
}

func ExtractTextFromDocx(file io.Reader) (string, error) {
	tempFile, err := os.CreateTemp("", "upload-*.docx")
	if err != nil {
		return "", err
	}
	defer os.Remove(tempFile.Name())

	_, err = io.Copy(tempFile, file)
	if err != nil {
		return "", err
	}
	tempFile.Close()

	doc, err := document.Open(tempFile.Name())
	if err != nil {
		return "", err
	}

	var textBuilder strings.Builder
	for _, para := range doc.Paragraphs() {
		for _, run := range para.Runs() {
			textBuilder.WriteString(run.Text())
		}
		textBuilder.WriteString("\n")
	}

	return textBuilder.String(), nil
}
