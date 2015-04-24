package interfaces

import (
	"github.com/AVANT/felicium/moonrakr/app/lib/results"
)

type Renderable interface {
	Render() (*results.JsonResult, error)
}
