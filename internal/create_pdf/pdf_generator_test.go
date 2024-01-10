package createpdf

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockedPDFFromBuffer struct {
	mock.Mock
}

type MockFileSystem struct {
	mock.Mock
}

func (m *MockFileSystem) ReadFile(filename string) ([]byte, error) {
	args := m.Called(filename)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockedPDFFromBuffer) combineHTMLFilesToBuffer(tempDir string) (*bytes.Buffer, string, error) {
	args := m.Called(tempDir)
	return args.Get(0).(*bytes.Buffer), args.String(), args.Error(2)
}

func (m *MockedPDFFromBuffer) generatePDFFromBuffer(buffer io.Reader, output, filename string) error {
	args := m.Called(buffer, output, filename)
	return args.Error(0)
}

func TestCreateSeparatedPDFFilesNoErrors(t *testing.T) {
	tempDir := t.TempDir()
	testObj := new(MockedPDFFromBuffer)

	buffer := new(bytes.Buffer)
	testObj.On("combineHTMLFilesToBuffer", tempDir).Return(buffer, "combinedOutput", nil)
	testObj.On("generatePDFFromBuffer", mock.AnythingOfType("*bytes.Buffer"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)

	err := CreateCombinedPDF(tempDir, "ARTICLENAME", testObj)

	if err != nil {
		t.Fatalf("Error testing CreateCombinedPDF: %v", err)
	}

	assert.Nil(t, err)
	testObj.AssertExpectations(t)
}

func TestCreateSeparatedPDFFilesErrorCreatingCombinedBuffer(t *testing.T) {
	tempDir := t.TempDir()
	testObj := new(MockedPDFFromBuffer)
	buffer := new(bytes.Buffer)

	testObj.On("combineHTMLFilesToBuffer", tempDir).Return(buffer, "combinedOutput", errors.New("Error creating buffer"))

	err := CreateCombinedPDF(tempDir, "combinedOutput", testObj)

	if assert.Error(t, err) {
		assert.EqualError(t, err, "Error creating buffer")
	}
	testObj.AssertExpectations(t)
}

func TestCreateSeparatedPDFFilesErrorCreatingPDF(t *testing.T) {
	tempDir := t.TempDir()
	testObj := new(MockedPDFFromBuffer)
	buffer := new(bytes.Buffer)

	testObj.On("combineHTMLFilesToBuffer", tempDir).Return(buffer, "combinedOutput", nil)
	testObj.On("generatePDFFromBuffer", buffer, "combinedOutput", "").Return(errors.New("Error creating PDF"))

	err := CreateCombinedPDF(tempDir, "combinedOutput", testObj)

	if assert.Error(t, err) {
		assert.EqualError(t, err, "Error creating PDF")
	}
	testObj.AssertExpectations(t)
}

func TestCreateSeparatedPDFFiles(t *testing.T) {
	tempDir := t.TempDir()

	mockObj := new(MockedPDFFromBuffer)
	mockFileSystem := new(MockFileSystem)
	mockFileSystem.On("ReadFile", filepath.Join(tempDir, "file1.html")).Return([]byte("<html><body>File 1 Content</body></html>"), nil).Once()
	mockFileSystem.On("ReadFile", filepath.Join(tempDir, "file2.html")).Return([]byte("<html><body>File 2 Content</body></html>"), nil).Once()
	mockFileSystem.On("ReadFile", filepath.Join(tempDir, "file3.html")).Return([]byte("<html><body>File 3 Content</body></html>"), nil).Once()
	mockObj.On("generatePDFFromBuffer", mock.AnythingOfType("*bytes.Buffer"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)

	createTestFile(t, tempDir, "file1.html", "<html><body>File 1 Content</body></html>")
	createTestFile(t, tempDir, "file2.html", "<html><body>File 2 Content</body></html>")
	createTestFile(t, tempDir, "file3.html", "<html><body>File 3 Content</body></html>")

	err := CreateSeparatedPDFFiles(tempDir, "output", mockFileSystem, mockObj)
	if err != nil {
		t.Fatalf("Error testing CreateSeparatedPDFFiles: %v", err)
	}

	assert.Nil(t, err)
	mockObj.AssertExpectations(t)
}

func TestCreateSeparatedPDFFilesErrorReadingFile(t *testing.T) {
	tempDir := t.TempDir()

	mockObj := new(MockedPDFFromBuffer)
	mockFileSystem := new(MockFileSystem)

	mockObj.On("generatePDFFromBuffer", mock.AnythingOfType("*bytes.Buffer"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
	mockFileSystem.On("ReadFile", filepath.Join(tempDir, "file1.html")).Return([]byte("<html><body>File 1 Content</body></html>"), nil).Once()
	mockFileSystem.On("ReadFile", mock.AnythingOfType("string")).Return([]byte(""), errors.New("Error reading file"))

	createTestFile(t, tempDir, "file1.html", "<html><body>File 1 Content</body></html>")
	createTestFile(t, tempDir, "file2.html", "<html><body>File 2 Content</body></html>")

	err := CreateSeparatedPDFFiles(tempDir, "output", mockFileSystem, mockObj)
	if assert.Error(t, err) {
		assert.EqualError(t, err, "failed to read file file2.html: Error reading file")
	}

	mockObj.AssertExpectations(t)
}

func TestCreateSeparatedPDFFilesTempDirError(t *testing.T) {
	mockObj := new(MockedPDFFromBuffer)
	MockFileSystem := new(MockFileSystem)

	err := CreateSeparatedPDFFiles("not_a_dir", "output", MockFileSystem, mockObj)

	if assert.Error(t, err) {
		assert.EqualError(t, err, "failed to read directory: open not_a_dir: no such file or directory")
	}
}

func TestReadFile(t *testing.T) {
	tempDir := t.TempDir()

	tempFile, err := os.CreateTemp(tempDir, "testfile")
	if err != nil {
		t.Fatal(err)
	}
	defer tempFile.Close()

	content := []byte("test content")
	if _, err := tempFile.Write(content); err != nil {
		t.Fatal(err)
	}

	fs := &RealFileSystem{}

	filePath := tempFile.Name()
	result, err := fs.ReadFile(filePath)

	if err != nil {
		t.Errorf("ReadFile failed: %v", err)
	}

	assert.Equal(t, string(result), string(content))
}

func createTestFile(t *testing.T, tempDir, filename, content string) {
	filePath := filepath.Join(tempDir, filename)
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file %s: %v", filename, err)
	}
}
