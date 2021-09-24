package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/gopherjs/gopherjs/js"
	client "github.com/revzim/amoeba-client"
	"github.com/revzim/amoeba/crypt"
	vue "github.com/revzim/gopherjs-vue"
	"honnef.co/go/js/dom"
)

type (
	Model struct {
		*js.Object
		MountEl      string                   `js:"mountEl"`
		ID           string                   `js:"id"`
		Token        string                   `js:"token"`
		SendBtnLabel string                   `js:"sendBtnLabel"`
		Username     string                   `js:"username"`
		InputMessage string                   `js:"inputMsg"`
		Messages     []map[string]interface{} `js:"messages"`
	}
)

const (
	VueAppMountElement = "#app"
	SendBtnLabel       = "SEND"
)

var (
	AmoebaAddress string = "localhost:80/ws"
	ServerToken   string
	ServerID      string = "test"
	AmoebaClient  *client.Connector
	c             = crypt.New([]byte(""))
)

func (m *Model) SendMessage() {
	msg := m.InputMessage
	mapData := map[string]interface{}{
		"name":    m.Username,
		"content": msg,
	}
	b, _ := json.Marshal(mapData)
	msgBytes, _ := c.Encrypt(b)
	log.Println("sending: ", string(b))
	AmoebaClient.Request("room.message", msgBytes, nil)
	m.InputMessage = ""
}

func DecryptPacket(msg []byte) (map[string]interface{}, error) {
	dataBytes, err := c.Decrypt(msg)
	if err != nil {
		return nil, err
	}
	var data interface{}
	err = json.Unmarshal(dataBytes, &data)
	if err != nil {
		log.Println("err parsing json bytes", string(msg), string(dataBytes))
		return nil, err
	}
	// log.Println(data)
	if data == nil {
		return nil, errors.New("recvd nil packet data")
	}
	return data.(map[string]interface{}), nil
}

func (m *Model) InitAmoebaClient(addr string) {
	AmoebaClient = client.NewConnector()

	err := AmoebaClient.InitReqHandshake("0.6.0", "amoeba-client", nil, map[string]interface{}{"name": "dude"})
	if err != nil {
		log.Fatal("Amoeba Handshake err: ", err)
	}
	err = AmoebaClient.InitHandshakeACK(1)
	if err != nil {
		panic(err)
	}
	// connected := false
	onMessage := func(data []byte) {
		decPkt, _ := DecryptPacket(data)
		log.Println("onMessage", decPkt) // string(data))
		if decPkt != nil {
			m.AddMessage(decPkt)
		}

	}
	AmoebaClient.Connected(func() {
		log.Printf("connected to server at: %s\n", addr)
		// connected = true
		err = AmoebaClient.Request("room.join", nil, func(data []byte) {
			decPkt, _ := DecryptPacket(data)
			log.Println("onJoinRoom", decPkt)
			// ON JOIN
			log.Println("on join!", data)
			dataCode := decPkt["code"].(float64)
			if dataCode == 0 {
				m.Username = decPkt["username"].(string)
				// msg := Encrypt(data["result"].(string))
				msg := decPkt["result"].(string)
				msgData := map[string]interface{}{
					"name":    "system",
					"content": msg,
				}
				m.AddMessage(msgData)

				// REGISTER ON MESSAGE
				AmoebaClient.On("onMessage", onMessage)

			}
		})
		if err != nil {
			panic(err)
		}
	})
	go func() {
		err := AmoebaClient.RunJS(addr, 10)
		if err != nil {
			panic(err)
		}
	}()
	AmoebaClient.On("onNewUser", func(data []byte) {
		decPkt, _ := DecryptPacket(data)
		log.Println("onNewUser", decPkt) // string(data))

		msgData := map[string]interface{}{
			"name":    "system",
			"content": decPkt["content"].(string),
		}
		m.AddMessage(msgData)

	})
	AmoebaClient.On("onMembers", func(data []byte) {
		decPkt, _ := DecryptPacket(data)
		log.Println("onMembers", decPkt) // string(data))

		var content string
		if decPkt["members"] != nil {
			content = fmt.Sprintf("active members: %v", decPkt["members"].([]interface{}))
		} else {
			content = "welcome"
		}
		msg := content
		msgData := map[string]interface{}{
			"name":    "system",
			"content": msg,
		}
		m.AddMessage(msgData)

	})
	defer AmoebaClient.Close()
}

func (m *Model) ConnectToServer(id, token string) {
	serverAddr := fmt.Sprintf("ws://%s?id=%s&token=%s", AmoebaAddress, id, token)
	log.Printf("attempting to connect to: %s...\n", serverAddr)
	m.InitAmoebaClient(serverAddr)
}

func (m *Model) AddMessage(data map[string]interface{}) {
	m.Messages = append(m.Messages, data)
	chatBoxElem := dom.GetWindow().Document().QuerySelector("#chat-box").Underlying()
	dom.GetWindow().SetTimeout(func() {
		chatBoxElem.Set("scrollTop", chatBoxElem.Get("scrollHeight"))
	}, 100)
}

func InitVuetify() *js.Object {
	Vuetify := js.Global.Get("Vuetify")

	vue.Use(Vuetify)

	vuetifyCFG := js.Global.Get("Object").New()
	vuetifyCFG.Set("theme", map[string]interface{}{
		"dark": true,
	})

	vtfy := Vuetify.New(vuetifyCFG)

	return vtfy
}

func InitVueOpts(m *Model, fn func(vm *vue.ViewModel)) *vue.Option {
	o := vue.NewOption()

	o.SetDataWithMethods(m)

	o.AddComputed("test", func(vm *vue.ViewModel) interface{} {
		return strings.ToUpper(vm.Data.Get("test").String())
	})

	o = o.Mixin(js.M{
		"vuetify": InitVuetify(),
		"mounted": js.MakeFunc(func(this *js.Object, arguments []*js.Object) interface{} {
			vm := &vue.ViewModel{
				Object: this,
			}
			fn(vm)
			return nil
		}),
	})

	return o
}

func main() {
	m := &Model{
		Object: js.Global.Get("Object").New(),
	}

	m.MountEl = VueAppMountElement
	m.SendBtnLabel = SendBtnLabel
	m.Messages = []map[string]interface{}{}
	m.Username = ""
	m.InputMessage = ""

	o := InitVueOpts(m, func(vm *vue.ViewModel) {
		// respToken := dom.GetWindow().Prompt("Input token", "server token")
		tknString := ""
		// log.Println("GOT TOKEN STRING!", tknString)
		m.ConnectToServer("test", tknString)
		// go func() {
		// 	resp, err := http.Get(fmt.Sprintf("http://%s/auth/token", dom.GetWindow().Location().Host))
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	defer resp.Body.Close()
		// 	var b []byte
		// 	b, err = ioutil.ReadAll(resp.Body)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	respMap := make(map[string]interface{})
		// 	err = json.Unmarshal(b, &respMap)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	// log.Println(respMap)
		// 	if respMap["data"] != nil {
		// 		tknString := ""
		// 		// log.Println("GOT TOKEN STRING!", tknString)
		// 		m.ConnectToServer("test", tknString)
		// 	}
		// }()

	})

	v := o.NewViewModel()

	v.Mount(VueAppMountElement)

}
