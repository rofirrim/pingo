package controllers

import "github.com/revel/revel"
import "errors"

type TelegramBot struct {
	*revel.Controller
}

func (c TelegramBot) Entry() revel.Result {
        return c.RenderError(errors.New("Not implemented yet"))
}
