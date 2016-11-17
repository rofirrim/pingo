package app

import "github.com/revel/revel"
import "database/sql"
import "fmt"
import "path/filepath"
import "encoding/json"
import "io/ioutil"
import _ "github.com/go-sql-driver/mysql"

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.ActionInvoker,           // Invoke the action.
	}

	// register startup functions with OnAppStart
	// ( order dependent )
	revel.OnAppStart(InitDB)
}

var DB *sql.DB
const LogsPerPage = 10

type dbConnectionInfo struct {
        Name string
        User string
        Pass string
}

type Settings struct {
    DB dbConnectionInfo
}

// revel uses a ConfPaths variable that includes both revel and app confs
var AppConfPath string

func loadSettings() (Settings, error) {
    var settings Settings
    data, err := ioutil.ReadFile(filepath.Join(AppConfPath, "settings.json"))
    if err != nil {
        return settings, err
    }
    err = json.Unmarshal(data, &settings)
    if err != nil {
        return settings, err
    }
    return settings, nil
}

func loadDB(DB dbConnectionInfo) (*sql.DB, error) {
    connection := fmt.Sprintf("%s:%s@/%s?charset=latin1", DB.User, DB.Pass, DB.Name)
	return sql.Open("mysql", connection)
}

func InitDB() {
    AppConfPath = filepath.Join(revel.BasePath, "conf")
	var err error
    var settings Settings
    settings, err = loadSettings()
	if err != nil {
		revel.ERROR.Println("Load settings", err)
	}

    DB, err = loadDB(settings.DB)

	if err != nil {
		revel.ERROR.Println("DB Error", err)
	}
	revel.INFO.Println("DB Connected")
}

// TODO turn this into revel.HeaderFilter
// should probably also have a filter for CSRF
// not sure if it can go in the same filter or not
var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	// Add some common security headers
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")

	fc[0](c, fc[1:]) // Execute the next filter stage.
}

