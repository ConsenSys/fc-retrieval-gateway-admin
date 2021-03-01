package settings

// Copyright (C) 2020 ConsenSys Software Inc

// Filecoin Retrieval Gateway Admin Client Settings

import (
	"time"
)

const (
	// DefaultGatewayGatewayPort default port for gateway to gatway communications.
	DefaultGatewayGatewayPort = "9010"
	// DefaultGatewayProviderPort default port for provider to gateway communications.
	DefaultGatewayProviderPort = "9011"
	// DefaultGatewayClientPort default port for client to gateway communications.
	DefaultGatewayClientPort = "9012"
	// DefaultGatewayAdminPort is the default gateway admin port
	DefaultGatewayAdminPort = "9013"

	// DefaultTCPInactivityTimeout is the default TCP timeout
	DefaultTCPInactivityTimeout = 100 * time.Millisecond


	// DefaultEstablishmentTTL is the default Time To Live used with Client - Gateway estalishment messages.
	defaultEstablishmentTTL = int64(100)

	// DefaultLogLevel is the default amount of logging to show.
	defaultLogLevel = "trace"

	// DefaultLogTarget is the default output location of log output.
	defaultLogTarget = "STDOUT"
)
