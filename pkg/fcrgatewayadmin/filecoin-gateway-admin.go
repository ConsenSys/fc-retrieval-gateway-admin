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

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	log "github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-gateway-admin/internal/control"
	"github.com/ConsenSys/fc-retrieval-gateway-admin/internal/settings"
)

// FilecoinRetrievalGatewayAdminClient holds information about the interaction of
// the Filecoin Retrieval Gateway Admin Client with Filecoin Retrieval Gateways.
type FilecoinRetrievalGatewayAdminClient struct {
	gatewayManager *control.GatewayManager
	// TODO have a list of gateway objects of all the current gateways being interacted with
}


// NewFilecoinRetrievalGatewayAdminClient initialise the Filecoin Retreival Client library
func NewFilecoinRetrievalGatewayAdminClient(conf Settings) *FilecoinRetrievalGatewayAdminClient {
	var c = FilecoinRetrievalGatewayAdminClient{}
	log.Info("Filecoin Retrieval Gateway Admin Client started")
	clientSettings := conf.(*settings.ClientGatewayAdminSettings)
	c.gatewayManager = control.NewGatewayManager(*clientSettings)
	return &c

}

// CreateKey creates a private key for a Gateway.
func CreateKey() (*fcrcrypto.KeyPair, error) {
	log.Info("Filecoin Retrieval Gateway Admin Client: RequestKeyCreation()")

	gatewayPrivateKey, err := fcrcrypto.GenerateRetrievalV1KeyPair()
	if err != nil {
		log.Error("Error creating Gateway Private Key: %s", err)
		return nil, err
	}

	return gatewayPrivateKey, nil
}

// InitializeGateway sends a private key to a Gateway along with a key version number.
func (c *FilecoinRetrievalGatewayAdminClient) InitializeGateway(gatewayDomain string, gatewayPrivateKey *fcrcrypto.KeyPair) error {
	log.Info("Filecoin Retrieval Gateway Admin Client: InitializeGateway()")
	return c.gatewayManager.InitializeGateway(gatewayDomain, gatewayPrivateKey)
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
