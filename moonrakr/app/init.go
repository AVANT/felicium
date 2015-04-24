package app

import (
	"github.com/AVANT/felicium/moonrakr/app/lib/boot"
	"github.com/robfig/revel"
)

func init() {

	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter, // Recover from panics and display an error page instead.
		// Nginx should make this less necessary we can take a closer look later
		//LogFilter,
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,
		revel.InterceptorFilter, // Run interceptors around the action.
		revel.ActionInvoker,     // Invoke the action.
	}

	///
	// On Start
	///

	//connect to database
	revel.OnAppStart(func() {
		boot.NormalBoot()
	})
}

var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	// Add CORS
	c.Response.Out.Header().Add("Access-Control-Allow-Origin", "*")
	c.Response.Out.Header().Add("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
	c.Response.Out.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	c.Response.Out.Header().Add("Content-Type", "application/json")

	fc[0](c, fc[1:]) // Execute the next filter stage.
}
