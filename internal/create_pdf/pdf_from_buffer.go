package createpdf

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

type PDFFromBuffer interface {
	combineHTMLFilesToBuffer(tempDir string) (*bytes.Buffer, string, error)
	generatePDFFromBuffer(buffer io.Reader, output, filename string) error
}

type PDFFromBufferImpl struct{}

func (PDFFromBufferImpl) combineHTMLFilesToBuffer(tempDir string) (*bytes.Buffer, string, error) {

	tempDirRead, err := os.ReadDir(tempDir)
	if err != nil {
		return nil, "output", fmt.Errorf("failed to read directory: %v", err)
	}

	var output strings.Builder
	buffer := new(bytes.Buffer)

	for i, file := range tempDirRead {
		log.Printf("Reading file: %s/%s", tempDir, file.Name())
		readFile, err := os.ReadFile(filepath.Join(tempDir, file.Name()))
		if err != nil {
			return nil, output.String(), fmt.Errorf("failed to read file %s: %v", file.Name(), err)
		}
		buffer.Write(readFile)

		if i > 0 {
			output.WriteRune('-')
		}
		output.WriteString(strings.TrimSuffix(file.Name(), filepath.Ext(file.Name())))
	}

	return buffer, output.String(), nil
}

func (PDFFromBufferImpl) generatePDFFromBuffer(buffer io.Reader, output, filename string) error {
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return fmt.Errorf("failed to create PDF generator: %v", err)
	}

	pdfg.AddPage(wkhtmltopdf.NewPageReader(buffer))

	err = pdfg.Create()
	if err != nil {
		return fmt.Errorf("failed to create PDF: %v", err)
	}

	if output == "ARTICLENAME" {
		output = strings.TrimSuffix(filename, filepath.Ext(filename))
	}

	log.Printf("Creating %s.pdf", output)
	err = pdfg.WriteFile(fmt.Sprintf("%s.pdf", output))
	if err != nil {
		return fmt.Errorf("failed to write PDF to file: %v", err)
	}

	return nil
}
