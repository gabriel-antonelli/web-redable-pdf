package createpdf

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type FileSystem interface {
	ReadFile(filename string) ([]byte, error)
}

type RealFileSystem struct{}

func (r *RealFileSystem) ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

func CreateCombinedPDF(tempDir, output string, PDFFromBuffer PDFFromBuffer) error {
	defer os.RemoveAll(tempDir)

	buffer, combinedOutput, err := PDFFromBuffer.combineHTMLFilesToBuffer(tempDir)
	if err != nil {
		return err
	}

	if output == "ARTICLENAME" {
		output = combinedOutput
	}

	err = PDFFromBuffer.generatePDFFromBuffer(buffer, output, "")
	if err != nil {
		return err
	}

	return nil
}

func CreateSeparatedPDFFiles(tempDir, output string, fs FileSystem, PDFFromBuffer PDFFromBuffer) error {
	defer os.RemoveAll(tempDir)

	tempDirRead, err := os.ReadDir(tempDir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %v", err)
	}

	for num, file := range tempDirRead {
		log.Printf("Reading file: %s/%s", tempDir, file.Name())
		readFile, err := fs.ReadFile(filepath.Join(tempDir, file.Name()))
		if err != nil {
			return fmt.Errorf("failed to read file %s: %v", file.Name(), err)
		}
		buffer := new(bytes.Buffer)
		buffer.Write(readFile)

		userOutput := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))

		if output != "ARTICLENAME" {
			userOutput = fmt.Sprintf("%d-%s", num+1, output)
		}

		PDFFromBuffer.generatePDFFromBuffer(buffer, userOutput, file.Name())
	}
	return nil
}
