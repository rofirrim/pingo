package tests

import "github.com/revel/revel/testing"
import "pingo/app/models"
import "encoding/json"
import "strings"
import "fmt"
import "net/http"

type AppTest struct {
	testing.TestSuite
}

func (t *AppTest) Before() {
	println("Set up")
}

func (t *AppTest) TestThatIndexPageWorks() {
	t.Get("/")
	t.AssertOk()
	t.AssertContentType("text/html; charset=utf-8")
}

// JSON endpoints
func (t *AppTest) TestThatWeCanUpload() {
	// Upload
	query := `{ "AuthToken" : "auth-token-test",
			"Upload" : {
			"Text": "New text!",
			"Protagonista": 1,
			"Autor": 1,
			"Titol": "New title!",
			"Data" : 1513932522
			} }`
	t.Post("/json/upload", "application/json", strings.NewReader(query))
	// Check the POST succeeds and returns the expected content type
	t.AssertOk()
	t.AssertContentType("application/json; charset=utf-8")

	// Load the JSON result
	var result models.JSONUploadResult
	err := json.Unmarshal(t.ResponseBody, &result)

	// Check we can unmarshal the answer correctly
	t.Assert(err == nil)
	// Check the result contains some plog id
	t.Assert(result.IdPlog != 0)

	// Check the new post does exist
	t.Get(fmt.Sprintf("/%d", result.IdPlog))
	t.AssertOk()
	t.AssertContentType("text/html; charset=utf-8")
	// Check the content makes sense
	str := string(t.ResponseBody[:])
	t.Assert(strings.Contains(str, "New text!"))
	t.Assert(strings.Contains(str, "New title!"))
	t.Assert(strings.Contains(str, "Log del 22/12/2017 a les 08:48 enviat per PaRaP"))
}

func (t *AppTest) TestUploadFailsAuth() {
	// Upload
	query := `{ "AuthToken" : "invalid-token-test",
			"Upload" : {
			"Text": "New text!",
			"Protagonista": 1,
			"Autor": 1,
			"Titol": "New title!",
			"Data" : 1513932522
			} }`
	t.Post("/json/upload", "application/json", strings.NewReader(query))
	t.Assert(t.Response.StatusCode == http.StatusForbidden)
}

func (t *AppTest) TestUploadInvalidSyntax() {
    // Note the trailing comma after Data
	query := `{ "AuthToken" : "invalid-token-test",
			"Upload" : {
			"Text": "New text!",
			"Protagonista": 1,
			"Autor": 1,
			"Titol": "New title!",
			"Data" : 1513932522,
			} }`
	t.Post("/json/upload", "application/json", strings.NewReader(query))
	t.Assert(t.Response.StatusCode == http.StatusBadRequest)
}

func (t *AppTest) After() {
	println("Tear down")
}
