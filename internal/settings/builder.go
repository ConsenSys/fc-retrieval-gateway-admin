package settings

// Copyright (C) 2020 ConsenSys Software Inc

// Filecoin Retrieval Gateway Admin Client Settings

import (
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	log "github.com/ConsenSys/fc-retrieval-common/pkg/logging"
)

// BuilderImpl holds the library configuration
type BuilderImpl struct {
	logLevel         string
	logTarget        string
	establishmentTTL int64

	blockchainPrivateKey *fcrcrypto.KeyPair

	gatewayAdminPrivateKey    *fcrcrypto.KeyPair
	gatewayAdminPrivateKeyVer *fcrcrypto.KeyVersion
	registerURL string
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

// SetRegisterURL sets the URL of the register service
func (f *BuilderImpl) SetRegisterURL(regURL string) {
	f.registerURL = regURL
}

// Build creates a settings object and initialises the logging system.
func (f *BuilderImpl) Build() *ClientGatewayAdminSettings {
	log.Init1(f.logLevel, f.logTarget)

	g := ClientGatewayAdminSettings{}
	g.establishmentTTL = f.establishmentTTL
	g.registerURL = f.registerURL

	if f.blockchainPrivateKey == nil {
		log.ErrorAndPanic("Settings: Blockchain Private Key not set")
	}
	g.blockchainPrivateKey = f.blockchainPrivateKey

	if f.gatewayAdminPrivateKey == nil {
		pKey, err := fcrcrypto.GenerateRetrievalV1KeyPair()
		if err != nil {
			log.ErrorAndPanic("Settings: Error while generating random retrieval key pair: %s" + err.Error())
		}
		g.gatewayAdminPrivateKey = pKey
		g.gatewayAdminPrivateKeyVer = fcrcrypto.DecodeKeyVersion(1)
	} else {
		g.gatewayAdminPrivateKey = f.gatewayAdminPrivateKey
		g.gatewayAdminPrivateKeyVer = f.gatewayAdminPrivateKeyVer
	}

	return &g
}
