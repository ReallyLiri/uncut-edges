package core

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func createPennHeader(catalogID string) (Header, error) {
	var res Header

	header, err := fetchPennHeader(catalogID)
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

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if s.Text() == "IIIF presentation manifest" {
			if href, exists := s.Attr("href"); exists {
				res.ManifestURL = href
			}
		}
	})

	return res, nil
}

func fetchPennHeader(catalogID string) ([]byte, error) {
	return fetch(fmt.Sprintf(pennPageUrlFmt, catalogID))
}

func createShakespearHeader(catalogID string) (Header, error) {
	var res Header

	header, err := fetchShakespearHeader(catalogID)
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

	doc.Find("div.field--name-field-display-title").Children().Each(func(i int, s *goquery.Selection) {
		res.Title = s.Text()
	})

	doc.Find("div.field").Children().Each(func(i int, child *goquery.Selection) {
		if child.HasClass("field__label") {
			labelText := child.Text()
			items := strings.TrimSpace(child.Next().Text())
			res.Properties = append(res.Properties, Property{labelText, items})
		}
	})

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if s.Text() == "IIIF Manifest" {
			if href, exists := s.Attr("href"); exists {
				res.ManifestURL = href
			}
		}
	})

	return res, nil
}

func fetchShakespearHeader(catalogID string) ([]byte, error) {
	return fetch(fmt.Sprintf(shakespeareUrlFmt, catalogID))
}

func fetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
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
