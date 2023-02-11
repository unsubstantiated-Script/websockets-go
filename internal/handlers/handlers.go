package handlers

import (
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sort"
)

//& check the address
//*int pointer that's a base type
//*p operator that returns the value of what p is pointing to

//Only accepts payload type
var wsChan = make(chan WsPayload)

//When someone connects to the connection, we'll add em to the map
var clients = make(map[WsConnection]string)

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

// WsConnection Websocket connection
type WsConnection struct {
	*websocket.Conn
}

// WsJsonResponse Structuring JSON for a response sent back from websocket
type WsJsonResponse struct {
	Action         string   `json:"action"`
	Message        string   `json:"message"`
	MessageType    string   `json:"message_type"`
	ConnectedUsers []string `json:"connected_users"`
}

type WsPayload struct {
	Action   string `json:"action"`
	Username string `json:"username"`
	Message  string `json:"message"`
	//How to leave something out of the json
	Conn WsConnection `json:"-"`
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

	//Creating a new Websocket connection and adding it to our list of clients
	conn := WsConnection{Conn: ws}
	clients[conn] = ""

	err = ws.WriteJSON(response)
	if err != nil {
		log.Println(err)
	}

	go ListenForWs(&conn)
}

// ListenForWs Putting our clients through a Go Routine
func ListenForWs(conn *WsConnection) {
	//When the main function here stops executing, execute this function. Mostly to fire off an error
	defer func() {
		//If there's a lockup, assign r to recover(). If r is equal to something, log the error. Will recover if things die.
		if r := recover(); r != nil {
			log.Println("Error", fmt.Sprintf("%v", r))
		}
	}()

	var payload WsPayload

	//Here we are listening for a payload
	for {
		err := conn.ReadJSON(&payload)

		if err != nil {
			//do nothing...
		} else {
			payload.Conn = *conn
			wsChan <- payload
		}
	}
}

func ListenToWsChannel() {
	var response WsJsonResponse

	for {
		//Reading from wsChan
		e := <-wsChan

		switch e.Action {
		case "username":
			//get a list of all users and send it back via broadcast
			clients[e.Conn] = e.Username
			users := getUserList()
			response.Action = "list_users"
			response.ConnectedUsers = users
			broadcastToAll(response)
		}

		//response.Action = "Got Here"
		//response.Message = fmt.Sprintf("Some Message, and action was %s", e.Action)
		//broadcastToAll(response)
	}
}

func getUserList() []string {

	//Building user list into a slice array
	var userList []string
	for _, x := range clients {
		userList = append(userList, x)
	}

	//Organizing the list in alphabetical order and returning
	sort.Strings(userList)
	return userList
}

func broadcastToAll(response WsJsonResponse) {
	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			log.Println("Websocket error")
			_ = client.Close()
			delete(clients, client)
		}
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
