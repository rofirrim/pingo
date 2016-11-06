package controllers

import "pinchito/app"
import "github.com/revel/revel"

type App struct {
	*revel.Controller
}

type PagerItem struct {
    Ellipsis bool;
    Current bool;
    Page int;
}

func buildPager(page int, numplogs int) []PagerItem {
    // Pager
    maxpages := numplogs / app.LogsPerPage
    if numplogs % app.LogsPerPage != 0 {
        maxpages += 1
    }

    links := make([]bool, maxpages)
    links[0] = true
    links[page - 1] = true
    links[len(links) - 1] = true
    for i := page - 2; i <= page + 2; i++ {
        j := i - 1
        if 0 < j && j < len(links) {
            links[j] = true
        }
    }

    var pager []PagerItem
    prevIsEllipsis := false
    for i, v := range links {
        if v {
            if (i + 1 != page) {
                pager = append(pager, PagerItem{false, false, i + 1})
            } else {
                pager = append(pager, PagerItem{false, true, i + 1})
            }
            prevIsEllipsis = false
        } else {
            if !prevIsEllipsis {
                pager = append(pager, PagerItem{true, false, 0})
                prevIsEllipsis = true
            }
        }
    }

    return pager
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
	return c.RenderTemplate("menu.html")
}

func (c App) ShowLog(id int) revel.Result {
	plog, err := GetPlog(id)
	if err != nil {
		revel.INFO.Println("Plog not found. Redirecting to menu", err)
		return c.Menu(1)
	} else {
        c.RenderArgs["plog"] = plog
		return c.RenderTemplate("single_log.html")
	}
}
