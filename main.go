package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sajad-dev/go-framwork/App/websocket"
	"github.com/sajad-dev/go-framwork/Config/setting"
	"github.com/sajad-dev/go-framwork/Database/connection"
	"github.com/sajad-dev/go-framwork/Database/migration"

	"github.com/sajad-dev/go-framwork/Route/api"
	"github.com/sajad-dev/go-framwork/Route/command"
)

func main() {
	file, _ := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(file)

	connection.Connection()
	if len(os.Args) > 2 {
		command.Handel(os.Args)
		return
	}
	go websocket.Handel()

	migration.Handel()
	api.RouteRun()

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println(err)
	}

	if !setting.DEBUG {
		defer log.Panicln("END PROGRAM")
	}
}
