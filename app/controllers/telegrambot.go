package controllers

import "github.com/revel/revel"
import "errors"
import "encoding/json"
import "pinchito/app/telegram/messages"
import "pinchito/app/telegram/commands"
import "strings"
import "fmt"
import "io/ioutil"

type TelegramBot struct {
	*revel.Controller
}

// Handlers
func doTest(c *TelegramBot, update messages.Update) revel.Result {
	inputMessage := update.Message
	from := inputMessage.From

    if from.LastName == nil {
        var noLastName string = "<sense cognom>"
        from.LastName = &noLastName
    }

    if from.Username == nil {
        var noUserName string = "<sense username>"
        from.LastName = &noUserName
    }

	var sendMessage commands.SendMessage
	sendMessage.ChatId = fmt.Sprintf("%d", inputMessage.Chat.Id)
	sendMessage.Text = fmt.Sprintf("Ordre de prova. Enviada per \"%s %s\" (conegut com \"%s\")",
		from.FirstName, *from.LastName, *from.Username)

	return c.RenderJson(sendMessage)
}

type HandlerType func(c *TelegramBot, update messages.Update) revel.Result
type CommandPair struct {
	Command string
	Handler HandlerType
}
type CommandList []CommandPair

// Handler list
var commandList CommandList = []CommandPair{
	CommandPair{"/test", doTest},
}

// Generic entry
func (c TelegramBot) Entry() revel.Result {
    jsonData, err := ioutil.ReadAll(c.Request.Body)
    defer c.Request.Body.Close()
	if err != nil {
        revel.ERROR.Println("Error reading data:", err)
		return c.RenderError(err)
    }

	var update messages.Update
	err = json.Unmarshal(jsonData, &update)
	if err != nil {
        revel.ERROR.Println("Could not decode input JSON", err)
		return c.RenderError(err)
	}

	if update.Message == nil {
		revel.ERROR.Println("Missing message")
		return c.RenderError(errors.New("Missing message"))
	}

	if update.Message.Text == nil {
		revel.ERROR.Println("Missing text")
		return c.RenderError(errors.New("Missing text"))
	}

	// Dispatcher
	for _, v := range commandList {
		if strings.HasPrefix(*update.Message.Text, v.Command) {
			return v.Handler(&c, update)
		}
	}

    revel.ERROR.Println("Uknown command")
    return c.RenderError(errors.New("Unknown command"))
}
