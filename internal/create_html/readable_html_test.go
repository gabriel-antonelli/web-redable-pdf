package createhtml_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	createhtml "github.com/gabriel-antonelli/web-redable-pdf/internal/create_html"
	"github.com/go-shiori/go-readability"
	"github.com/stretchr/testify/assert"
)

var title = "testing it"

type ReadabilityService struct{}
type FailingReadabilityService struct{}

func (r *FailingReadabilityService) FromURL(url string, timeout time.Duration) (readability.Article, error) {
	return readability.Article{}, fmt.Errorf("error")
}

func (r *ReadabilityService) FromURL(url string, timeout time.Duration) (readability.Article, error) {
	return readability.Article{
		Title:   title,
		Content: "<p>teste</p>",
	}, nil
}

func TestGenerateRedableHtmlSuccess(t *testing.T) {
	urls := []string{
		"https://example.com/article1",
		"https://example.com/article2",
	}

	readabilityService := ReadabilityService{}
	tempDir, err := createhtml.GenerateReadableHtml(urls, &readabilityService)
	if err != nil {
		t.Fatalf("error while testing GenerateReadableHtml: %v", err)
	}
	defer os.RemoveAll(tempDir)

	_, errStat := os.Stat(tempDir)
	assert.NoError(t, errStat, "Temp directory should be created")

	assert.NoError(t, err, "Temp directory should be created")
	assert.FileExists(t, fmt.Sprintf("%s/1---%s.html", tempDir, title))
}

func TestGenerateRedableHtmlFailToParse(t *testing.T) {
	urls := []string{
		"https://example.com/articleFail",
	}

	readabilityService := FailingReadabilityService{}
	_, err := createhtml.GenerateReadableHtml(urls, &readabilityService)
	if err != nil {
		assert.Error(t, err)
	}
}
