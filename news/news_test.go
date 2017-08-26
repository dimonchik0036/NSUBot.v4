package news

import "testing"

func TestGetNewsPage(t *testing.T) {
	_, err := getNewsPage("http://fit.nsu.ru/")
	if err != nil {
		t.Fatal(err)
	}
}
