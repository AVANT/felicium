package results

import (
	"encoding/xml"
	"net/http"

	"github.com/robfig/revel"
)

type RenderXmlResultWithHeader struct {
	Obj interface{}
}

func (r RenderXmlResultWithHeader) Apply(req *revel.Request, resp *revel.Response) {
	resp.WriteHeader(http.StatusOK, "application/xml; charset=utf-8")
	resp.Out.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>`))
	xmlEncoder := xml.NewEncoder(resp.Out)
	xmlEncoder.Encode(r.Obj)
}
