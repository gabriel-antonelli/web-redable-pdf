package createpdf

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	wkhtmltopdf "github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

func CreateCombinedPDF(tempDir, output string) error {
	defer os.RemoveAll(tempDir)

	buffer, combinedOutput, err := combineHTMLFilesToBuffer(tempDir)
	if err != nil {
		log.Fatal(err)
	}

	if output == "ARTICLENAME" {
		output = combinedOutput
	}

	err = generatePDFFromBuffer(buffer, output, "")
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
		readFile, err := os.ReadFile(filepath.Join(tempDir, file.Name()))
		if err != nil {
			log.Fatalf("failed to read file %s: %v", file.Name(), err)
		}
		buffer := new(bytes.Buffer)
		buffer.Write(readFile)

		userOutput := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))

		if output != "ARTICLENAME" {
			userOutput = fmt.Sprintf("%d-%s", num+1, output)
		}

		generatePDFFromBuffer(buffer, userOutput, file.Name())
	}
}

func combineHTMLFilesToBuffer(tempDir string) (*bytes.Buffer, string, error) {
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

func generatePDFFromBuffer(buffer io.Reader, output, filename string) error {
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
