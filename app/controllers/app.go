package controllers

import "pinchito/app"
import "github.com/revel/revel"
import "net/http"

type App struct {
	*revel.Controller
}

type PagerItem struct {
    Ellipsis bool;
    IsCurrent bool;
    Page int;
}

type Pager struct
{
    HasPrev bool;
    PrevPage int;
    HasNext bool;
    NextPage int;
    Items []PagerItem;
}

func computeCookie(renderMap *map[string]interface{}) {
    cookie, err := GetRandomCookie()
    if err != nil {
        return
    }
    (*renderMap)["cookie"] = cookie
}

func (c App) FinishAndRender(template string) revel.Result {

    computeCookie(&c.RenderArgs)
    return c.RenderTemplate(template)
}

func buildPager(page int, numplogs int) Pager {
	// Pager
	numpages := numplogs / app.LogsPerPage
	if numplogs%app.LogsPerPage != 0 {
		numpages += 1
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

    c.RenderArgs["plogs"] = plogs
    c.RenderArgs["pager"] = pager
	return c.FinishAndRender("menu.html")
}

func (c App) ShowLog(id int) revel.Result {
	plog, err := GetPlog(id)
	if err != nil {
		revel.INFO.Println("Plog not found. Redirecting to menu", err)
		return c.Menu(1)
	} else {
        c.RenderArgs["plog"] = plog
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
    c.RenderArgs["plogs"] = plogs
    return c.FinishAndRender("top20.html")
}

func (c App) Random() revel.Result {
    plogs, err := GetRandomPlogs()
    if err != nil {
        // Abandon all hope here
		revel.ERROR.Println("Error when showing page", err)
        return c.RenderError(err)
    }
    c.RenderArgs["plogs"] = plogs
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
