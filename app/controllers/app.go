package controllers

import "github.com/revel/revel"

import "pinchito/app"
import "pinchito/app/helpers"
import "pinchito/app/models"

import "net/http"
import "strings"

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
	t, err := helpers.ProcessLogText(plog.Text, keywords)
	if err != nil {
		revel.ERROR.Println("Error when processing text of log", err)
	} else {
		plog.Text = t
	}
	t, err = helpers.ProcessLogTitle(plog.Titol, keywords)
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
