package ping

import "github.com/google/wire"

var Set = wire.NewSet(
	NewController,
)
