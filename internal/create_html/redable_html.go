package createhtml

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-shiori/go-readability"
)

func GenerateRedableHtml(args []string) string {
	tempDir, err := os.MkdirTemp("", "redable-html")
	if err != nil {
		log.Fatalf("Error creating temp directory: %s", err)
	}

	for num, webPage := range args {
		log.Printf("Finding text for article number %d in link: %s", num+1, webPage)

		article, err := readability.FromURL(webPage, 30*time.Second)
		if err != nil {
			log.Fatalf("failed to parse %s, %v\n", webPage, err)
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
			log.Fatalf("Error creating HTML file: %s", err)
		}
		defer dstHTMLFile.Close()

		dstHTMLFile.WriteString(article.Content)
	}

	return tempDir
}
