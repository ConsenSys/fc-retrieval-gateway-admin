package settings

// Copyright (C) 2020 ConsenSys Software Inc

// Filecoin Retrieval Gateway Admin Client Settings

import (
	"time"
)

const (
	// DefaultGatewayBindAdminAPI is the default gateway admin port
	DefaultGatewayBindAdminAPI = "9013"

	// DefaultTCPInactivityTimeout is the default TCP timeout
	DefaultTCPInactivityTimeout = 100 * time.Millisecond


	// DefaultEstablishmentTTL is the default Time To Live used with Client - Gateway estalishment messages.
	defaultEstablishmentTTL = int64(100)

	// DefaultLogLevel is the default amount of logging to show.
	defaultLogLevel = "trace"

	// DefaultLogTarget is the default output location of log output.
	defaultLogTarget = "STDOUT"
)
