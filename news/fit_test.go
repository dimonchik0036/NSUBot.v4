package news

import "testing"

func TestFitAnnounce(t *testing.T) {
	Fit("http://fit.nsu.ru/news/announc", 5)
}
