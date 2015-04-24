package results

import (
	"bufio"
	"bytes"
	"fmt"

	"github.com/emilsjolander/goson"
	"github.com/robfig/revel"
)

//A type that allows us to define our own result in revel
type JsonResult struct {
	*bytes.Buffer
}

type RenderableCollection interface {
	goson.Collection
	RenderElement(i int) (*JsonResult, error)
}

//Apply allows this to be used as a result in revel
func (json *JsonResult) Apply(req *revel.Request, resp *revel.Response) {
	if resp.Status == 0 {
		resp.Status = 200
	}
	resp.WriteHeader(resp.Status, "application/json")
	json.WriteTo(resp.Out)
}

func RenderRenderableCollection(collection RenderableCollection) (*JsonResult, error) {

	buffer := new(bytes.Buffer)
	w := bufio.NewWriter(buffer)
	defer w.Flush()

	fmt.Fprint(w, "[")

	for i := 0; i < collection.Len(); i++ {
		partial, err := collection.RenderElement(i)
		//bail out on error
		if err != nil {
			w.Flush()
			return new(JsonResult), err
		}
		fmt.Fprint(w, partial.String())
		//skip the comma on the last itteration
		if i < collection.Len()-1 {
			fmt.Fprint(w, ",")
		}
	}
	fmt.Fprint(w, "]")
	return &JsonResult{buffer}, nil
}
