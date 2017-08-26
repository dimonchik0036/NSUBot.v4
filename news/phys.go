package news

const (
	PhysHref       = "http://www.phys.nsu.ru"
	PhysTimeLayout = "02-01-2006"
)

/*func Phys(href string, count int) (news []News, err error) {
	//44896
	body, err := getNewsPage(PhysHref + href + "?limit=" + strconv.Itoa(count))
	if err != nil {
		return []News{}, err
	}
	println(len(body))
	rg, err := regexp.Compile("<table class=\"tablelist\">.*?</table>")
	if err != nil {
		return []News{}, err
	}

	body = rg.Find(body)
	println(string(body))
	os.Exit(0)
	hrefs := hrefProcessing(body, count)

	for _, v := range hrefs {
		news = append(news, News{
			ID:    idScan(string(v[0])),
			Title: html.UnescapeString(string(v[1])),
			OptionalURL:   PhysHref + string(v[0]),
		})
	}

	return
}*/
