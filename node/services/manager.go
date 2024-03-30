package services

type IClientManager[ClientType any] interface {
	GetPrimaryClient() ClientType
	GetFallbackClient() ClientType
	IsPrimaryReady() bool
	IsFallbackReady() bool
	IsFallbackEnabled() bool
	GetClientTypeName() string

	// Internal functions
	setPrimaryReady(bool)
	setFallbackReady(bool)
}
