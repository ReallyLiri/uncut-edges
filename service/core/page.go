package core

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func createHeaderFile(catalogID string) (Header, error) {
	var res Header

	header, err := fetchHeader(catalogID)
	if err != nil {
		return res, err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(header))
	if err != nil {
		return res, fmt.Errorf("error parsing html: %w", err)
	}

	res.CatalogID = catalogID
	res.Title = doc.Find("[itemprop='name']").Text()
	res.Properties = []Property{}
	res.Links = []string{}

	doc.Find("#sidebar a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		if strings.HasPrefix(href, "/") {
			return
		}
		res.Links = append(res.Links, href)
	})

	doc.Find("dl").Children().Each(func(i int, child *goquery.Selection) {
		if child.Is("dt") {
			dtText := child.Text()
			ddText := child.Next().Text()
			res.Properties = append(res.Properties, Property{dtText, ddText})
		}
	})

	return res, nil
}

func fetchHeader(catalogID string) ([]byte, error) {
	resp, err := http.Get(fmt.Sprintf(pageUrlFmt, catalogID))
	if err != nil {
		return nil, fmt.Errorf("error getting header: %w", err)
	}
	defer resp.Body.Close()
	header, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading header: %w", err)
	}
	return header, nil
}
