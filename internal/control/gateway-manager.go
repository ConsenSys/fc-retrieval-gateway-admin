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
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrtcpcomms"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/register"

	log "github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway-admin/internal/gatewayapi"
	"github.com/ConsenSys/fc-retrieval-gateway-admin/internal/settings"
)

// GatewayManager managers the pool of gateways and the connections to them.
type GatewayManager struct {
	settings                    settings.ClientGatewayAdminSettings
	gateway                     ActiveGateway
	gatewaysLock                sync.RWMutex
	registeredMap								map[string]register.RegisteredNode
	conxPool 									  *fcrtcpcomms.CommunicationPool
}

// ActiveGateway contains information for a single gateway
type ActiveGateway struct {
	comms *gatewayapi.Comms
}

// NewGatewayManager creates a gateway manager.
func NewGatewayManager(conf settings.ClientGatewayAdminSettings) *GatewayManager {
	g := GatewayManager{}
	g.settings = conf
	g.registeredMap = make(map[string]register.RegisteredNode)
	g.conxPool = fcrtcpcomms.NewCommunicationPool(g.registeredMap, &sync.RWMutex{})
	return &g
}

// InitializeGateway initialise a new gateway
func (g *GatewayManager) InitializeGateway(
	gatewayDomain, gatewayGatewayPort, gatewayProviderPort, gatewayClientPort, gatewayAdminPort string, 
	region string,
	gatewayRootKeyPair *fcrcrypto.KeyPair, gatewayRetrievalKeyPair *fcrcrypto.KeyPair) error {
	// TODO check whether gateway not initialized.
	// TODO check whether contract indicates initialised

	// Get gateway key version
	gatewayRetrievalKeyVersion := fcrcrypto.InitialKeyVersion()
	gatewayRetrievalKeyVersionUint := gatewayRetrievalKeyVersion.EncodeKeyVersion()
	// Get encoded version of the gateway's private key
	gatewayRetrievalPrivateKeyStr := gatewayRetrievalKeyPair.EncodePrivateKey()

	// Make a request message
	request, err := fcrmessages.EncodeAdminAcceptKeyChallenge(gatewayRetrievalPrivateKeyStr, gatewayRetrievalKeyVersionUint)
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

	// Get the gateway's NodeID
	gatewayNodeID, err := nodeid.NewNodeIDFromPublicKey(gatewayRootKeyPair)
	if err != nil {
		log.Error("Error getting gateway's NodeID: %s", err)
		return err
	}

	log.Info("Sending message to gateway: %v, message: %s", gatewayNodeID.ToString(), request.DumpMessage())

	conn, err := g.getConnection(gatewayNodeID, gatewayDomain, gatewayAdminPort) //"gateway:9013"
	if err != nil {
		return err
	}
	err = fcrtcpcomms.SendTCPMessage(conn, request, settings.DefaultTCPInactivityTimeout)
	if err != nil {
		log.Error("Error sending private key to Gateway: %s", err)
		return err
	}

	// Process the response from the gateway.
	response, err := fcrtcpcomms.ReadTCPMessage(conn, time.Second*1)
	log.Info("Response message: %+v", response)
	if response.MessageType != fcrmessages.AdminAcceptKeyResponseType {
		// TODO other types of messages such as protocol version negotiation need to be handled.
		return fmt.Errorf("Unexpected message in response to set-up Gateway message: %d", response.MessageType)
	}
	keyAccepted, err := fcrmessages.DecodeAdminAcceptKeyResponse(response)
	if err != nil {
		return err
	}
	if !keyAccepted {
		return fmt.Errorf("Key not accepted for unspecified reason")
	}



	gatewayRootPubKey, err := gatewayRootKeyPair.EncodePublicKey()
	if err != nil {
		return err
	}
	gatewayRetrievalPubKey, err := gatewayRetrievalKeyPair.EncodePublicKey()
	if err != nil {
		return err
	}

	// Register Gateway
	gatewayReg := register.GatewayRegister{
			NodeID:              gatewayNodeID.ToString(),
			Address:             gatewayDomain,
			RootSigningKey:      gatewayRootPubKey,
			SigningKey:          gatewayRetrievalPubKey,
			NetworkInfoGateway:  gatewayGatewayPort,
			NetworkInfoProvider: gatewayProviderPort,
			NetworkInfoClient:   gatewayClientPort,
			NetworkInfoAdmin:    gatewayAdminPort,
			RegionCode:          region,
		}
		gatewayReg.RegisterGateway(g.settings.RegisterURL())

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


func (g *GatewayManager) getConnection(gatewayNodeID *nodeid.NodeID, domain string, port string) (net.Conn, error) {
	domainAndPort := domain + ":" + port

	// Add new gateway to the connection pool.
	g.registeredMap[gatewayNodeID.ToString()] = &register.GatewayRegister{
		NodeID:             gatewayNodeID.ToString(),
		NetworkInfoGateway: domainAndPort,
	}

	// Get conn for the right gateway
	channel, err := g.conxPool.GetConnForRequestingNode(gatewayNodeID, fcrtcpcomms.AccessFromGateway)
	if err != nil {
		return nil, err
	}
	conn := channel.Conn
	if err != nil {
		log.Error("Error getting a connection to gateway %v: %s", gatewayNodeID.ToString(), err)
		return nil, err
	}
	return conn, nil

}