package controllers

import (
	"encoding/xml"
	"strings"

	"github.com/AVANT/felicium/moonrakr/app/lib/results"
	"github.com/AVANT/felicium/moonrakr/app/models"
	"github.com/robfig/revel"
)

type SiteMap struct {
	*revel.Controller
}

type SiteMapPages struct {
	XMLName    xml.Name       `xml:"urlset"`
	XmlNS      string         `xml:"xmlns,attr"`
	XmlImageNS string         `xml:"xmlns:image,attr"`
	XmlNewsNS  string         `xml:"xmlns:news,attr"`
	Pages      []*SiteMapPage `xml:"url"`
}

type SiteMapPage struct {
	XMLName  xml.Name `xml:"url"`
	Loc      string   `xml:"loc"`
	Name     string   `xml:"news:news>news:publication>news:name"`
	Language string   `xml:"news:news>news:publication>news:language"`
	Title    string   `xml:"news:news>news:title"`
	Keywords string   `xml:"news:news>news:keywords"`
	Image    string   `xml:"image:image>image:loc"`
}

func (s *SiteMap) Index() revel.Result {
	posts, _ := models.GetPostsByStatus("published")
	docList := new(SiteMapPages)
	docList.XmlImageNS = "http://www.google.com/schemas/sitemap-image/1.1"
	docList.XmlNewsNS = "http://www.google.com/schemas/sitemap-news/0.9"
	for _, v := range *posts {
		doc := new(SiteMapPage)
		doc.Loc = "http://www.example.org/something/" + v.GetSlug()
		doc.Name = "VVVNT"
		doc.Language = "en"
		doc.Title = v.GetTitle()
		doc.Keywords = strings.Join(v.GetTags(), ",")
		doc.Image = v.GetHeaderImage().GetUrl()
		docList.Pages = append(docList.Pages, doc)
	}
	toReturn := new(results.RenderXmlResultWithHeader)
	(*toReturn).Obj = docList
	return toReturn
}
