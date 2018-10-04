package cube_websocket_gateway

import (
	"fmt"
	"strconv"

	"github.com/akaumov/cube"
	"github.com/akaumov/cube-websocket-gateway/lib"
)

const Version = "1"

type BusSubject string
type Uri string

type Handler struct {
	cubeInstance           cube.Cube
	timeoutMs              uint64
	onlyAuthorizedRequests bool
	jwtSecret              string
	devMode                bool
	port                   int
	server                 *lib.Server
}

func (h *Handler) OnInitInstance() []cube.InputChannel {
	return []cube.InputChannel{}
}

func (h *Handler) OnStart(cubeInstance cube.Cube) error {
	fmt.Println("Starting http gateway...")

	h.cubeInstance = cubeInstance
	h.jwtSecret = cubeInstance.GetParam("jwtSecret")
	h.onlyAuthorizedRequests = cubeInstance.GetParam("onlyAuthorizedRequests") == "true"
	h.devMode = cubeInstance.GetParam("dev") == "true"

	portString := cubeInstance.GetParam("port")

	var err error
	port := 80

	if portString != "" {
		port, err = strconv.Atoi(portString)
		if err != nil {
			cubeInstance.LogError("Wrong timeout")
			return err
		}
	}

	h.port = port

	h.timeoutMs = 30000
	timeoutString := cubeInstance.GetParam("timeoutMs")

	if timeoutString != "" {
		timeoutMs, err := strconv.ParseUint(timeoutString, 10, 64)
		if err != nil {
			cubeInstance.LogError("Wrong timeout")
			return err
		}

		h.timeoutMs = timeoutMs
	}

	h.server = lib.NewServer(cubeInstance, h.devMode, h.onlyAuthorizedRequests, h.jwtSecret)
	go h.server.Start(cubeInstance)
	return nil
}

func (h *Handler) OnStop(c cube.Cube) {
}

func (h *Handler) OnReceiveMessage(instance cube.Cube, channel cube.Channel, message cube.Message) {
	fmt.Println("OnReceiveMessage: is not implemented")
	instance.LogError("OnReceiveMessage: is not implemented")
}

//From bus
func (h *Handler) OnReceiveRequest(instance cube.Cube, channel cube.Channel, request cube.Request) cube.Response {
	fmt.Println("OnReceiveRequest: is not implemented")
	instance.LogError("OnReceiveRequest: is not implemented")
	return cube.NewErrorResponse(
		"",
		"NotImplemented",
		"",
	)
}

// func (h *Handler) packRequest(userId *string, deviceId *string, request *http.Request) (*cube.Request, error) {
// 	var err error
// 	var body []byte

// 	if request.Body != nil {
// 		body, err = ioutil.ReadAll(request.Body)
// 		if err != nil {
// 			return nil, err
// 		}
// 		request.Body.Close()

// 		if body != nil && len(body) == 0 {
// 			body = nil
// 		}
// 	}

// 	headers := map[string][]string{}

// 	for key, value := range request.Header {
// 		headers[key] = value
// 	}

// 	params := js.RequestParams{
// 		DeviceId:   deviceId,
// 		UserId:     userId,
// 		Method:     request.Method,
// 		InputTime:  time.Now().UnixNano(),
// 		Host:       request.Host,
// 		RequestURI: request.RequestURI,
// 		Body:       body,
// 		RemoteAddr: request.RemoteAddr,
// 		Headers:    headers,
// 	}

// 	packedParams, err := json.Marshal(params)

// 	requestData := &cube.Request{
// 		Method: request.Method,
// 		Params: (*json.RawMessage)(&packedParams),
// 	}

// 	return requestData, nil
// }

var _ cube.HandlerInterface = (*Handler)(nil)
