package settings

// Copyright (C) 2020 ConsenSys Software Inc

// Filecoin Retrieval Gateway Admin Client Settings

import (
	"github.com/ConsenSys/fc-retrieval-gateway-admin/config"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrcrypto"
	log "github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
)

// BuilderImpl holds the library configuration
type BuilderImpl struct {
	logLevel         string
	logTarget        string
	establishmentTTL int64

	blockchainPrivateKey *fcrcrypto.KeyPair

	gatewayAdminPrivateKey    *fcrcrypto.KeyPair
	gatewayAdminPrivateKeyVer *fcrcrypto.KeyVersion
}

// CreateSettings creates an object with the default settings.
func CreateSettings() *BuilderImpl {
	f := BuilderImpl{}
	f.logLevel = defaultLogLevel
	f.logTarget = defaultLogTarget
	f.establishmentTTL = defaultEstablishmentTTL
	return &f
}

// SetLogging sets the log level and target.
func (f *BuilderImpl) SetLogging(logLevel string, logTarget string) {
	f.logLevel = defaultLogLevel
	f.logTarget = defaultLogTarget
}

// SetEstablishmentTTL sets the time to live for the establishment message between client and gateway.
func (f *BuilderImpl) SetEstablishmentTTL(ttl int64) {
	f.establishmentTTL = ttl
}

// SetBlockchainPrivateKey sets the blockchain private key.
func (f *BuilderImpl) SetBlockchainPrivateKey(bcPkey *fcrcrypto.KeyPair) {
	f.blockchainPrivateKey = bcPkey
}

// SetGatewayAdminPrivateKey sets the private key used for authenticating to the gateway
func (f *BuilderImpl) SetGatewayAdminPrivateKey(rPkey *fcrcrypto.KeyPair, ver *fcrcrypto.KeyVersion) {
	f.gatewayAdminPrivateKey = rPkey
	f.gatewayAdminPrivateKeyVer = ver
}

// Build creates a settings object and initialises the logging system.
func (f *BuilderImpl) Build() *ClientGatewayAdminSettings {
	conf := config.NewConfig()
	log.Init(conf)

	g := ClientGatewayAdminSettings{}
	g.establishmentTTL = f.establishmentTTL

	if f.blockchainPrivateKey == nil {
		log.ErrorAndPanic("Settings: Blockchain Private Key not set")
	}
	g.blockchainPrivateKey = f.blockchainPrivateKey

	if f.gatewayAdminPrivateKey == nil {
		pKey, err := fcrcrypto.GenerateRetrievalV1KeyPair()
		if err != nil {
			log.ErrorAndPanic("Settings: Error while generating random retrieval key pair: %s" + err.Error())
		}
		f.gatewayAdminPrivateKey = pKey
		f.gatewayAdminPrivateKeyVer = fcrcrypto.DecodeKeyVersion(1)
	} else {
		g.gatewayAdminPrivateKey = f.gatewayAdminPrivateKey
		g.gatewayAdminPrivateKeyVer = f.gatewayAdminPrivateKeyVer
	}

	return &g
}
