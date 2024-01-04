package createpdf

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	wkhtmltopdf "github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

func combineHTMLFilesToBuffer(tempDir string) (*bytes.Buffer, error) {
	tempDirRead, err := os.ReadDir(tempDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %v", err)
	}

	buffer := new(bytes.Buffer)
	for _, file := range tempDirRead {
		log.Printf("Reading file: %s/%s", tempDir, file.Name())
		readFile, err := os.ReadFile(fmt.Sprintf("%s/%s", tempDir, file.Name()))
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %v", file.Name(), err)
		}
		buffer.Write(readFile)
	}

	return buffer, nil
}

func generatePDFFromBuffer(buffer io.Reader, output string) error {
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return fmt.Errorf("failed to create PDF generator: %v", err)
	}

	pdfg.AddPage(wkhtmltopdf.NewPageReader(buffer))

	err = pdfg.Create()
	if err != nil {
		return fmt.Errorf("failed to create PDF: %v", err)
	}

	err = pdfg.WriteFile(fmt.Sprintf("%s.pdf", output))
	if err != nil {
		return fmt.Errorf("failed to write PDF to file: %v", err)
	}

	return nil
}

func CreateCombinedPDF(tempDir, output string) error {
	defer os.RemoveAll(tempDir)

	buffer, err := combineHTMLFilesToBuffer(tempDir)
	if err != nil {
		log.Fatal(err)
	}

	err = generatePDFFromBuffer(buffer, output)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func CreateSeparatedPDFFiles(tempDir, output string) {
	defer os.RemoveAll(tempDir)

	tempDirRead, err := os.ReadDir(tempDir)
	if err != nil {
		log.Fatalf("failed to read directory: %v", err)
	}

	for num, file := range tempDirRead {
		log.Printf("Reading file: %s/%s", tempDir, file.Name())
		readFile, err := os.ReadFile(fmt.Sprintf("%s/%s", tempDir, file.Name()))
		if err != nil {
			log.Fatalf("failed to read file %s: %v", file.Name(), err)
		}
		buffer := new(bytes.Buffer)
		buffer.Write(readFile)
		generatePDFFromBuffer(buffer, fmt.Sprintf("%s-%s", output, fmt.Sprint(num)))
	}
}
