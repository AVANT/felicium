package controllers

import (
	"fmt"
	"strings"
	"time"

	"github.com/AVANT/felicium/moonrakr/app/lib/results"
	"github.com/AVANT/felicium/moonrakr/app/models"
	"github.com/gorilla/feeds"
	"github.com/robfig/revel"
)

type Feed struct {
	*revel.Controller
}

func (r *Feed) Rss() revel.Result {
	posts, _ := models.GetPostsByStatus("published")
	now := time.Now()
	rss := feeds.RssFeed{
		Title:         "AVANT.org",
		Description:   "/ˈavänt/ An online journal, forum, & project space for sharing ways of thinking with practical connections across time, scale, system, & discipline.",
		WebMaster:     "editorial @ avant.org (AVANT.org editorial)",
		LastBuildDate: now.Format("Mon, 02 Jan 2006 15:04:05 MST"),
		Copyright:     "Attribution-NonCommercial-ShareAlike 4.0 International",
		Link:          "http://testing.avant.org/rss",
		Image: &feeds.RssImage{
			Url:    "http://jeroboam.s3.amazonaws.com/production/2071b317-b31b-46d8-7bb1-2a1f0a76dc38.jpg",
			Title:  "AVANT.org",
			Link:   "http://avant.org",
			Width:  144,
			Height: 144,
		},
		Ttl: 1,
	}
	for _, v := range *posts {
		rss.Items = append(rss.Items, &feeds.RssItem{
			Title: v.GetTitle(),
			Link:  "http://avant.org/media/" + v.GetSlug(),
			Author: func() string {
				var toReturn []string
				for _, value := range v.GetAuthorsArray() {
					toReturn = append(toReturn, "editorial @ avant.org ("+value+")")
				}
				return strings.Join(toReturn, ",")
			}(),
			Description: v.GetExcerpt(),
			Enclosure: func(m *models.Media) *feeds.RssEnclosure {
				enclosure := new(feeds.RssEnclosure)
				enclosure.Url = m.GetUrl()
				enclosure.Length = "10152"
				enclosure.Type = m.GetContentType()
				return enclosure
			}(v.GetHeaderImage()),
		})
	}
	//
	//		feed := &feeds.Feed{
	//		Title:       "AVANT.org",
	//		Description: "/ˈavänt/ An online journal, forum, & project space for sharing ways of thinking with practical connections across time, scale, system, & discipline.",
	//		Author:      &feeds.Author{"Avant.org", "editorial@avant.org"},
	//		Created:     now,
	//		Copyright:   "Attribution-NonCommercial-ShareAlike 4.0 International",
	//		Link:        &feeds.Link{Href: "http://avant.org/rss"},
	//	}
	//	feed.Items = []*feeds.Item{}
	//	for _, v := range *posts {
	//		feed.Add(
	//			&feeds.Item{
	//				Title:       v.GetTitle(),
	//				Link:        &feeds.Link{Href: "http://avant.org/media/" + v.GetSlug()},
	//				Description: v.GetExcerpt(),
	//				Author: &feeds.Author{
	//					func() string {
	//						var toReturn []string
	//						for _, value := range v.GetAuthorsArray() {
	//							toReturn = append(toReturn, value)
	//						}
	//						return strings.Join(toReturn, ",")
	//					}(),
	//					"editorial@avant.org",
	//				},
	//				Created: now,
	//			},
	//		)
	//	}
	fmt.Printf("%+v\n", rss)
	toReturn := new(results.RenderXmlResultWithHeader)
	(*toReturn).Obj = rss.FeedXml()
	return toReturn
}

// func (r *Feed) Rss() revel.Result {
// 	posts, _ := models.GetPostsByStatus("published")
// 	rss := new(RssFeedXml)
// 	rss.Version = "2.0"
// 	rss.Channel = new(RssFeed)
// 	rss.Channel.Title = "AVANT.org"
// 	rss.Channel.Link = "http://avant.org/rss"
// 	rss.Channel.Description = "/ˈavänt/ An online journal, forum, & project space for sharing ways of thinking with practical connections across time, scale, system, & discipline."
// 	rss.Channel.Copyright = "Attribution-NonCommercial-ShareAlike 4.0 International"
// 	rss.Channel.PubDate = time.Now().Format("Jan 2, 2006 at 3:04pm (MST)")
// 	rss.Channel.LastBuildDate = time.Now().Format("Jan 2, 2006 at 3:04pm (MST)")
// 	rss.Channel.Ttl = 60
// 	for _, v := range *posts {
// 		doc := new(RssItem)
// 		doc.Link = "http://avant.org/media/" + v.GetSlug()
// 		doc.Title = v.GetTitle()
// 		doc.Description = v.GetExcerpt()
// 		doc.PubDate = v.GetCreatedAt().Format("Jan 2, 2006 at 3:04pm (MST)")
// 		doc.Guid = "avant.org/media/" + v.GetSlug()
// 		doc.Enclosure = func(m *models.Media) *RssEnclosure {
// 			enclosure := new(RssEnclosure)
// 			enclosure.Url = m.GetUrl()
// 			enclosure.Length = "10152"
// 			enclosure.Type = m.GetContentType()
// 			return enclosure
// 		}(v.GetHeaderImage())
// 		doc.Author = func() string {
// 			var toReturn []string
// 			// this is for later on
// 			// for _, value := range v.GetAuthors() {
// 			// 	toReturn = append(toReturn, value.GetFullName())
// 			// }
// 			for _, value := range v.GetAuthorsArray() {
// 				toReturn = append(toReturn, value)
// 			}
// 			return strings.Join(toReturn, ",")
// 		}()
// 		rss.Channel.Items = append(rss.Channel.Items, doc)
// 	}
// 	toReturn := new(results.RenderXmlResultWithHeader)
// 	(*toReturn).Obj = rss
// 	return toReturn
// }
