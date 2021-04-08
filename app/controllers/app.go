package controllers

import "github.com/revel/revel"

import "pingo/app"
import "pingo/app/helpers"
import "pingo/app/models"

import "net/http"
import "strings"
import "strconv"
import "time"

type App struct {
	*revel.Controller
}

type PagerItem struct {
	Ellipsis  bool
	IsCurrent bool
	Page      int
}

type Pager struct {
	HasPrev  bool
	PrevPage int
	HasNext  bool
	NextPage int
	Items    []PagerItem
}

func computeCookie(renderMap *map[string]interface{}) {
	cookie, err := GetRandomCookie()
	if err != nil {
		return
	}
	(*renderMap)["cookie"] = cookie
}

func processLogHighlights(plog *models.Plog, keywords []string) {
	t, err := helpers.ProcessLogText(plog.RawText, keywords)
	if err != nil {
		revel.ERROR.Println("Error when processing text of log", err)
	} else {
		plog.Text = t
	}
	t, err = helpers.ProcessLogTitle(plog.RawTitol, keywords)
	if err != nil {
		revel.ERROR.Println("Error when processing title of log", err)
	} else {
		plog.Titol = t
	}
}

func processLog(plog *models.Plog) {
	processLogHighlights(plog, []string{})
}

func processLogsHighlights(plogs *[]models.Plog, keywords []string) {
	for i := range *plogs {
		processLogHighlights(&((*plogs)[i]), keywords)
	}
}

func processLogs(plogs *[]models.Plog) {
	processLogsHighlights(plogs, []string{})
}

func (c App) FinishAndRender(template string) revel.Result {

	// Make sure "menuitem" is defined
	_, menuitem := c.ViewArgs["menuitem"]
	if !menuitem {
		c.ViewArgs["menuitem"] = ""
	}
	computeCookie(&c.ViewArgs)
	return c.RenderTemplate(template)
}

func buildPager(page int, numplogs int) Pager {
	// Pager
	numpages := numplogs / app.LogsPerPage
	if numplogs%app.LogsPerPage != 0 {
		numpages += 1
	}

	// Degenerated case
	if numpages == 0 {
		return Pager{}
	}

	links := make([]bool, numpages)
	links[0] = true
	links[page-1] = true
	links[len(links)-1] = true
	for i := page - 2; i <= page+2; i++ {
		j := i - 1
		if 0 < j && j < len(links) {
			links[j] = true
		}
	}

	var pagerItems []PagerItem
	prevIsEllipsis := false
	for i, v := range links {
		if v {
			if i+1 != page {
				pagerItems = append(pagerItems, PagerItem{false, false, i + 1})
			} else {
				pagerItems = append(pagerItems, PagerItem{false, true, i + 1})
			}
			prevIsEllipsis = false
		} else {
			if !prevIsEllipsis {
				pagerItems = append(pagerItems, PagerItem{true, false, 0})
				prevIsEllipsis = true
			}
		}
	}

	return Pager{page > 1, page - 1, page < numpages, page + 1, pagerItems}
}

func (c App) Index() revel.Result {
	return c.Menu(1)
}

func (c App) Menu(page int) revel.Result {
	if page <= 0 {
		page = 1
	}
	revel.INFO.Println("Requesting page", page)

	var numplogs int
	plogs, err := GetPlogBunch(page, &numplogs)
	if err != nil {
		// Abandon all hope here
		revel.ERROR.Println("Error when showing page", err)
		return c.RenderError(err)
	}

	pager := buildPager(page, numplogs)

	processLogs(&plogs)
	c.ViewArgs["plogs"] = plogs

	c.ViewArgs["pager"] = pager
	c.ViewArgs["menuitem"] = "menu"

	return c.FinishAndRender("menu.html")
}

func (c App) ShowLog(id int) revel.Result {
	plog, err := GetPlog(id)
	if err != nil {
		revel.INFO.Println("Plog not found. Redirecting to menu", err)
		return c.Menu(1)
	} else {
		processLog(&plog)
		c.ViewArgs["plog"] = plog
		return c.FinishAndRender("single_log.html")
	}
}

func (c App) Top20() revel.Result {
	plogs, err := GetTop20Plogs()
	if err != nil {
		// Abandon all hope here
		revel.ERROR.Println("Error when showing page", err)
		return c.RenderError(err)
	}

	processLogs(&plogs)
	c.ViewArgs["plogs"] = plogs
	c.ViewArgs["menuitem"] = "especialitats"

	return c.FinishAndRender("top20.html")
}

func (c App) Random() revel.Result {
	plogs, err := GetRandomPlogs()
	if err != nil {
		// Abandon all hope here
		revel.ERROR.Println("Error when showing page", err)
		return c.RenderError(err)
	}

	processLogs(&plogs)
	c.ViewArgs["plogs"] = plogs
	c.ViewArgs["menuitem"] = "tapeta"

	return c.FinishAndRender("random.html")
}

type BlobBytes []byte

func (b BlobBytes) Apply(req *revel.Request, resp *revel.Response) {
	resp.WriteHeader(http.StatusOK, "image/png")
	resp.Out.Write(b)
}

func (c App) Avatar(id int) revel.Result {
	blob, err := GetBlobAvatar(id)
	if err != nil {
		// Abandon all hope here
		revel.ERROR.Println("Error when showing page", err)
		return c.RenderError(err)
	}

	return BlobBytes(blob)
}

func (c App) Search(page int) revel.Result {
	if page <= 0 {
		page = 1
	}
	keywords := strings.Split(c.Params.Get("s"), " ")
	if len(keywords) == 0 {
		return c.Index()
	}

	var numplogs int
	plogs, err := SearchPlogs(keywords, page, &numplogs)
	if err != nil {
		// Abandon all hope here
		revel.ERROR.Println("Error when showing page", err)
		return c.RenderError(err)
	}

	pager := buildPager(page, numplogs)

	processLogsHighlights(&plogs, keywords)
	c.ViewArgs["plogs"] = plogs

	c.ViewArgs["pager"] = pager

	c.ViewArgs["query"] = c.Params.Get("s")

	return c.FinishAndRender("search.html")
}

func (c App) EditLog(id int) revel.Result {
	plog, err := GetPlog(id)
	if err != nil {
		revel.INFO.Println("Plog not found. Redirecting to menu", err)
		return c.Menu(1)
	}
	processLog(&plog)

	users, err := GetUsers()
	if err != nil {
		revel.INFO.Println("Error obtaining users", err)
		return c.RenderError(err)
	}

	c.ViewArgs["plog"] = plog
	c.ViewArgs["users"] = users
	return c.FinishAndRender("edit_log.html")
}

func (c App) makePlogFromPOST(id int) (models.Plog, error) {
	var plog models.Plog
	var err error
	var user_id int

	plog.Id = id
	plog.RawText = c.Params.Form.Get("log_text")
	plog.Text = plog.RawText

	user_id, err = strconv.Atoi(c.Params.Form.Get("protagonista"))
	if err != nil {
		revel.INFO.Println("protagonista is not integer", err)
		return models.Plog{}, err
	}
	plog.Protagonista, err = GetUser(user_id)
	if err != nil {
		revel.INFO.Println("Error obtaining protagonista.", err)
		return models.Plog{}, err
	}

	var t time.Time
	t, err = time.Parse("2006-01-02", c.Params.Form.Get("dia"))
	if err != nil {
		revel.INFO.Println("Wrong day", err)
		return models.Plog{}, err
	}
	plog.DiaYMD = t.Format("2006-01-02")
	plog.Dia = t.Format("02/01/2006")

	t, err = time.Parse("15:04", c.Params.Form.Get("hora"))
	if err != nil {
		revel.INFO.Println("Wrong hour", err)
		return models.Plog{}, err
	}
	plog.Hora = t.Format("15:04")

	user_id, err = strconv.Atoi(c.Params.Form.Get("autor"))
	if err != nil {
		revel.INFO.Println("autor is not integer", err)
		return models.Plog{}, err
	}
	plog.Autor, err = GetUser(user_id)
	if err != nil {
		revel.INFO.Println("Error obtaining autor.", err)
		return models.Plog{}, err
	}

	plog.RawTitol = c.Params.Form.Get("titol")
	plog.Titol = plog.RawTitol

	return plog, nil
}

func (c App) SubmitEditLog(id int) revel.Result {
	plog, err := c.makePlogFromPOST(id)
	if err != nil {
		revel.INFO.Println("Error making plog from form POST", err)
		return c.RenderError(err)
	}

	err = UpdatePlog(plog)
	if err != nil {
		revel.INFO.Println("Error updating DB", err)
		return c.RenderError(err)
	}

	users, err := GetUsers()
	if err != nil {
		revel.INFO.Println("Error obtaining users", err)
		return c.Menu(1)
	}

	// Reload plog from the DB.
	plog, err = GetPlog(id)
	processLog(&plog)

	c.ViewArgs["plog"] = plog
	c.ViewArgs["users"] = users
	return c.FinishAndRender("edit_log.html")
}

// --------------
// JSON endpoints
// --------------

func (c App) ShowLogJSON(id int) revel.Result {
	plog, err := GetPlog(id)

	if err != nil {
		return c.NotFound("Log not found")
	} else {
		return c.RenderJSON(plog)
	}
}

func (c App) RandomJSON() revel.Result {
	plogs, err := GetRandomPlogs()
	if err != nil {
		revel.ERROR.Println("Error when showing page", err)
		return c.RenderError(err)
	}

	return c.RenderJSON(plogs[0])
}

func (c App) SearchJSON(page int) revel.Result {
	keywords := strings.Split(c.Params.Get("s"), " ")
	if len(keywords) == 0 {
		return c.NotFound("Log not found")
	}
	var numplogs int
	plogs, err := SearchPlogs(keywords, page, &numplogs)
	if err != nil {
		c.RenderError(err)
	}
	if len(plogs) == 0 {
		return c.NotFound("Log not found")
	}
	return c.RenderJSON(plogs)
}

func (c App) UploadJSON() revel.Result {
	var uploadJSON models.JSONUploadOp
	err := c.Params.BindJSON(&uploadJSON)
	if err != nil {
		revel.ERROR.Println("Error when binding JSON info", err)
		c.Response.Status = http.StatusBadRequest
		return c.RenderJSON(models.JSONUploadResult{false, 0, err.Error()})
	}

	// Check the shared secret
	if uploadJSON.AuthToken != app.AuthToken || uploadJSON.AuthToken == "" {
		return c.Forbidden("Invalid AuthToken")
	}

	plogData := uploadJSON.Upload

	// Check input
	_, err = GetUser(plogData.Protagonista)
	// FIXME: Refactor if we need to add more checks like the ones below.
	if err != nil {
		revel.ERROR.Println("Error retrieving Protagonista", err)
		c.Response.Status = http.StatusUnprocessableEntity
		return c.RenderJSON(models.JSONUploadResult{false, 0, err.Error()})
	}
	_, err = GetUser(plogData.Autor)
	if err != nil {
		revel.ERROR.Println("Error retrieving Autor", err)
		c.Response.Status = http.StatusUnprocessableEntity
		return c.RenderJSON(models.JSONUploadResult{false, 0, err.Error()})
	}

	// Process input
	var idPlog int
	idPlog, err = UploadPlog(plogData)

	if err != nil {
		c.Response.Status = http.StatusInternalServerError
		return c.RenderJSON(models.JSONUploadResult{false, 0, err.Error()})
	}

	return c.RenderJSON(models.JSONUploadResult{true, idPlog, ""})
}
