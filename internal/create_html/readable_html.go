package createhtml

import (
	"fmt"
	"log"
	"os"
	"time"

	readability "github.com/go-shiori/go-readability"
)

type ReadabilityService interface {
	FromURL(url string, timeout time.Duration) (readability.Article, error)
}

func GenerateReadableHtml(args []string, readabilityService ReadabilityService) (string, error) {
	tempDir, err := os.MkdirTemp("", "redable-html")
	if err != nil {
		return "", err
	}

	for num, webPage := range args {
		log.Printf("Finding text for article number %d in link: %s", num+1, webPage)

		article, err := readabilityService.FromURL(webPage, 30*time.Second)
		if err != nil {
			return "", err
		}

		// Add <meta charset="utf-8"> to the content
		article.Content = fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
</head>
<body>
%s
</body>
</html>`, article.Content)

		dstHTMLFile, err := os.Create(fmt.Sprintf("%s/%d---%s.html", tempDir, num+1, article.Title))
		if err != nil {
			return "", err
		}
		defer dstHTMLFile.Close()

		dstHTMLFile.WriteString(article.Content)
	}

	return tempDir, nil
}
