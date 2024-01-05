package createhtml

import (
	"time"

	readability "github.com/go-shiori/go-readability"
)

type RealReadabilityService struct{}

func (r *RealReadabilityService) FromURL(url string, timeout time.Duration) (readability.Article, error) {
	return readability.FromURL(url, timeout)
}
