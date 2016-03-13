package html

import (
    // . "ml/dict"
    . "ml/array"
    . "ml/strings"
    . "ml/trace"

    "github.com/PuerkitoBio/goquery"
)

type Document struct {
    *goquery.Document
}

func Parse(html String) *Document {
    doc, err := goquery.NewDocumentFromReader(html.NewReader())
    RaiseIf(err)
    return &Document{
        Document: doc,
    }
}

func (self *Document) GetAllElements(ids Array) map[string]*goquery.Selection {
    data := map[string]*goquery.Selection{}

    self.Find("[id]").Each(func(i int, s *goquery.Selection) {
        i, exist := ids.Index(s.Attr2("id"))
        if exist {
            data[ids[i].(string)] = s
        }
    })

    return data
}
