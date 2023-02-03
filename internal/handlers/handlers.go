package handlers

import (
	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var views = jet.NewSet(
	//Setting up JET
	jet.NewOSFileSystemLoader("./html"),
	jet.InDevelopmentMode(),
)

//Setting up the first websocket a struct of sorts
var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Home handler renders the home page
func Home(w http.ResponseWriter, r *http.Request) {
	//Storing the err while instantiating. Pretty nifty.
	err := renderPage(w, "home.jet", nil)
	if err != nil {
		log.Println(err)
	}
}

// WsJsonResponse Structuring JSON for a response sent back from websocket
type WsJsonResponse struct {
	Action      string `json:"action"`
	Message     string `json:"message"`
	MessageType string `json:"message_type"`
}

// WsEndpoint Handler that handles our WebSocket Endpoint
func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	//Upgrading the connection to a websocket
	//Response writer, request, and header (we don't use the header here)
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("Client connected to endpoint")

	var response WsJsonResponse
	response.Message = `<em><small>Connected to the server boy! You got socket action going!</small></em>`

	err = ws.WriteJSON(response)
	if err != nil {
		log.Println(err)
	}
}

// renderPage renders a jet template
func renderPage(w http.ResponseWriter, tmpl string, data jet.VarMap) error {
	view, err := views.GetTemplate(tmpl)
	if err != nil {
		log.Println(err)
		return err
	}

	err = view.Execute(w, data, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
