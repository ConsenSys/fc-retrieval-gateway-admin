package control

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
	"sync"

	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrtcpcomms"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"

	"github.com/ConsenSys/fc-retrieval-gateway-admin/internal/contracts"
	"github.com/ConsenSys/fc-retrieval-gateway-admin/internal/gatewayapi"
	"github.com/ConsenSys/fc-retrieval-gateway-admin/internal/settings"
	log "github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
)

// GatewayManager managers the pool of gateways and the connections to them.
type GatewayManager struct {
	settings                    settings.ClientGatewayAdminSettings
	gatewayRegistrationContract *contracts.GatewayRegistrationContract
	gateway                     ActiveGateway
	gatewaysLock                sync.RWMutex
}

// ActiveGateway contains information for a single gateway
type ActiveGateway struct {
	info  contracts.GatewayInformation
	comms *gatewayapi.Comms
}


// NewGatewayManager returns the single instance of the gateway manager.
// The settings parameter must be used with the first call to this function.
// After that, the settings parameter is ignored.
func NewGatewayManager(conf settings.ClientGatewayAdminSettings) *GatewayManager {
	g := GatewayManager{}
	g.settings = conf
	g.gatewayRegistrationContract = contracts.GetGatewayRegistrationContract()


	// TODO what should be done with error that is returned possibly in the future?
	// TODO would it be better just to have gatewayManagerRunner panic after emitting a log?
	return &g
}



// InitializeGateway initialise a new gateway
func (g *GatewayManager) InitializeGateway(gatewayDomain string, gatewayKeyPair *fcrcrypto.KeyPair) error {
// TODO check whether gateway not initialized.
// TODO check whether contract indicates initialised

	// Get gateway key version
	gatewaykeyversion := fcrcrypto.InitialKeyVersion() 
	gatewaykeyversionuint := gatewaykeyversion.EncodeKeyVersion()
	// Get encoded version of the gateway's private key
	gatewayprivatekeystr := gatewayKeyPair.EncodePrivateKey()

	// Make a request message
	request, err := fcrmessages.EncodeAdminAcceptKeyChallenge(gatewayprivatekeystr, gatewaykeyversionuint)
	if err != nil {
		log.Error("Internal error in encoding AdminAcceptKeyChallenge message.")
		return nil
	}

	// Sign the request
	if request.SignMessage(func(msg interface{}) (string, error) {
		return fcrcrypto.SignMessage(g.settings.GatewayAdminPrivateKey(), g.settings.GatewayAdminPrivateKeyVer(), msg)

	}) != nil {
		log.Error("Error signing message for sending private key to gateway: %+v", err)
		return err
	}

	// TODO Temporary: The ConnectionPool should be a client-wide persistent struct
	conxPool := fcrtcpcomms.NewCommunicationPool()

	// TODO: Register the gateway with the connection pool.

	// Get the gateway's NodeID
	gatewaynodeid, err := nodeid.NewNodeIDFromPublicKey(gatewayKeyPair)
	if err != nil {
		log.Error("Error getting gateway's NodeID: %s", err)
		return err
	}

	// TODO: Persistence the gateway's keys and NodeID locally

	log.Info("Sending message to gateway: %v, message: %v", gatewaynodeid.ToString(), request.GetMessageBody())

	// Get conn for the right gateway
	channel, err := conxPool.GetConnForRequestingNode(gatewaynodeid)
	if err != nil {
		return err
	}
	conn := channel.Conn
	if err != nil {
		log.Error("Error getting a connection to gateway %v: %s", gatewaynodeid.ToString(), err)
		return err
	}
	err = fcrtcpcomms.SendTCPMessage(conn, request, settings.DefaultTCPInactivityTimeout)

	if err != nil {
		log.Error("Error sending private key to Gateway: %s", err)
		return err
	}

	// TODO: Receive response from gateway?

	return nil
}

// BlockGateway adds a host to disallowed list of gateways
func (g *GatewayManager) BlockGateway(hostName string) {
	// TODO
}

// UnblockGateway add a host to allowed list of gateways
func (g *GatewayManager) UnblockGateway(hostName string) {
	// TODO

}

// Shutdown stops go routines and closes sockets. This should be called as part
// of the graceful library shutdown
func (g *GatewayManager) Shutdown() {
	// TODO
}
