package gatewayapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	"net/http"

	log "github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/bitly/go-simplejson"
)

const (
	apiURLStart string = "http://"
	apiURLEnd   string = "/v1"

//	apiURLEnd string = "/client/establishment"
)

const (
	gatewayAdminClientAPIProtocolVersion     = 1
	gatewayAdminClientAPIProtocolSupportedHi = 1
)

// Can't have constant slices so create this at runtime.
// Order the API versions from most desirable to least desirable.
var gatewayAdminClientAPIProtocolSupported []int

// Comms holds the communications specific data
type Comms struct {
	apiURL string
	nodeID *nodeid.NodeID
}

// NewGatewayAPIComms creates a connection with a gateway
func NewGatewayAPIComms(host string, nodeID *nodeid.NodeID) (*Comms, error) {
	// Create the constant array.
	if gatewayAdminClientAPIProtocolSupported == nil {
		gatewayAdminClientAPIProtocolSupported = make([]int, 1)
		gatewayAdminClientAPIProtocolSupported[0] = gatewayAdminClientAPIProtocolSupportedHi
	}

	// Check that the host name is valid
	err := validateHostName(host)
	if err != nil {
		log.Error("Host name invalid: %s", err.Error())
		return nil, err
	}

	netComms := Comms{}
	netComms.apiURL = apiURLStart + host + apiURLEnd
	netComms.nodeID = nodeID
	return &netComms, nil
}

// GatewayCall calls the Gateway's REST API
func (n *Comms) gatewayCall(method int32, args map[string]interface{}) *simplejson.Json {
	args["protocol_version"] = int32(1)
	args["protocol_supported"] = []int32{1}
	args["message_type"] = method
	args["node_id"] = n.nodeID.ToString()
	mJSON, _ := json.Marshal(args)
	log.Info("JSON sent: %s", string(mJSON))
	contentReader := bytes.NewReader(mJSON)
	req, _ := http.NewRequest("POST", n.apiURL, contentReader)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, errs := client.Do(req)
	if errs != nil {
		log.ErrorAndPanic(errs.Error())
	}

	data, _ := ioutil.ReadAll(resp.Body)
	log.Info("response body: %s", string(data))

	js, err := simplejson.NewJson(data)
	if err != nil {
		log.ErrorAndPanic("Error decoding JSON: %s", err.Error())
	}

	return js
}

func validateHostName(host string) error {
	if len(host) == 0 {
		return errors.New("Error: Host name empty")
	}

	ra, err := net.ResolveIPAddr("ip4", host)
	if err != nil {
		return errors.New("Error: DNS Look-up failed for host: " + host)
	}
	log.Info("Resolved %s as %s\n", host, ra.String())
	return nil
}
