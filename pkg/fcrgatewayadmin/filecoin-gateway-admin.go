package fcrgatewayadmin

/*
 * Copyright 2020 ConsenSys Software Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

import (
	"container/list"

	"github.com/ConsenSys/fc-retrieval-gateway-admin/internal/control"
	"github.com/ConsenSys/fc-retrieval-gateway-admin/internal/settings"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrtcpcomms"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	log "github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
)

// FilecoinRetrievalGatewayAdminClient holds information about the interaction of
// the Filecoin Retrieval Gateway Admin Client with Filecoin Retrieval Gateways.
type FilecoinRetrievalGatewayAdminClient struct {
	gatewayManager *control.GatewayManager
	// TODO have a list of gateway objects of all the current gateways being interacted with
}

var singleInstance *FilecoinRetrievalGatewayAdminClient
var initialised = false

// InitFilecoinRetrievalGatewayAdminClient initialise the Filecoin Retreival Client library
func InitFilecoinRetrievalGatewayAdminClient(settings Settings) *FilecoinRetrievalGatewayAdminClient {
	if initialised {
		log.ErrorAndPanic("Attempt to init Filecoin Retrieval Gateway Admin Client a second time")
	}
	var c = FilecoinRetrievalGatewayAdminClient{}
	c.startUp(settings)
	singleInstance = &c
	initialised = true
	return singleInstance

}

// GetFilecoinRetrievalGatewayAdminClient creates a Filecoin Retrieval Gateway Admin Client
func GetFilecoinRetrievalGatewayAdminClient() *FilecoinRetrievalGatewayAdminClient {
	if !initialised {
		log.ErrorAndPanic("Filecoin Retrieval Gateway Admin Client not initialised")
	}

	return singleInstance
}

func (c *FilecoinRetrievalGatewayAdminClient) startUp(conf Settings) {
	log.Info("Filecoin Retrieval Gateway Admin Client started")
	clientSettings := conf.(*settings.ClientGatewayAdminSettings)
	c.gatewayManager = control.GetGatewayManager(*clientSettings)
}

// CreateKey creates a private key for a Gateway.
func CreateKey() (*fcrcrypto.KeyPair, *fcrcrypto.KeyVersion, error) {
	log.Info("Filecoin Retrieval Gateway Admin Client: RequestKeyCreation()")

	gatewayPrivateKey, err := fcrcrypto.GenerateRetrievalV1KeyPair()
	if err != nil {
		logging.Error("Error creating Gateway Private Key: %s", err)
		return nil, nil, err
	}

	keyversion := fcrcrypto.InitialKeyVersion()

	return gatewayPrivateKey, keyversion, nil
}

// SendKeyToGateway sends a private key to a Gateway along with a key version number.
func SendKeyToGateway(gateway string, gatewayprivatekey *fcrcrypto.KeyPair) error {
	log.Info("Filecoin Retrieval Gateway Admin Client: SendKeyToGateway()")

	// Get gateway key version
	var gatewaykeyversion *fcrcrypto.KeyVersion
	gatewaykeyversion = fcrcrypto.InitialKeyVersion() // TODO: Shouldn't this be the /current/ version??
	gatewaykeyversionuint := gatewaykeyversion.EncodeKeyVersion()
	// Get encoded version of the gateway's private key
	gatewayprivatekeystr := gatewayprivatekey.EncodePrivateKey()

	// Make a request message
	settingsBuilder := CreateSettings()
	conf := settingsBuilder.Build()
	clientprivatekey := (*conf).RetrievalPrivateKey()
	clientprivatekeyversion := fcrcrypto.InitialKeyVersion() // TODO: Shouldn't this be the /current/ version??
	request, err := fcrmessages.EncodeAdminAcceptKeyChallenge(gatewayprivatekeystr, gatewaykeyversionuint)
	if err != nil {
		logging.Error("Internal error in encoding AdminAcceptKeyChallenge message.")
		return nil
	}

	// Sign the request
	if request.SignMessage(func(msg interface{}) (string, error) {
		return fcrcrypto.SignMessage(clientprivatekey, clientprivatekeyversion, msg)

	}) != nil {
		logging.Error("Error signing message for sending private key to gateway: %+v", err)
		return err
	}

	// TODO Temporary: The ConnectionPool should be a client-wide persistent struct
	conxPool := fcrtcpcomms.NewCommunicationPool()

	// - get the public key of the gateway and its NodeID
	gatewaypublickey, err := gatewayprivatekey.EncodePublicKey()
	if err != nil {
		logging.Error("Error encoding gateway's public key: %s", err)
		return err
	}
	gatewaynodeid := nodeid.NewNodeIDFromString(gatewaypublickey)

	log.Info("Sending message to gateway: %v, message: %v", gatewaynodeid.ToString(), request.GetMessageBody())

	// Get conn for the right gateway
	err = fcrtcpcomms.SendTCPMessage(conn, request, DefaultTCPInactivityTimeout)

	if err != nil {
		logging.Error("Error sending private key to Gateway: %s", err)
		return err
	}

	// TODO: How to receive response from gateway?

	return nil
}

// InitialiseClientReputation requests a Gateway to initialise a client's reputation to the default value.
func InitialiseClientReputation(clientID *nodeid.NodeID) bool {
	log.Info("Filecoin Retrieval Gateway Admin Client: InitialiseClientReputation(clientID: %s", clientID)
	// TODO DHW
	log.Info("InitialiseClientReputation(clientID: %s) failed to initialise reputation.", clientID)
	return false
}

// SetClientReputation requests a Gateway to set a client's reputation to a specified value.
func SetClientReputation(clientID *nodeid.NodeID, rep int64) bool {
	log.Info("Filecoin Retrieval Gateway Admin Client: SetClientReputation(clientID: %s, reputation: %d", clientID, rep)
	// TODO DHW
	log.Info("SetClientReputation(clientID: %s, reputation: %d) failed to set reputation.", clientID, rep)
	return false
}

// GetCIDOffersList requests a Gateway's current list of CID Offers.
func GetCIDOffersList() *list.List {
	log.Info("Filecoin Retrieval Gateway Admin Client: GetCIDOffersList()")
	// TODO
	log.Info("GetCIDOffersList() failed to find any CID Offers.")
	emptyList := list.New()
	return emptyList
}

// Shutdown releases all resources used by the library
func (c *FilecoinRetrievalGatewayAdminClient) Shutdown() {
	log.Info("Filecoin Retrieval Gateway Admin Client shutting down")
	c.gatewayManager.Shutdown()
}
