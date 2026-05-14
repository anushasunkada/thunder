package oauth

import "testing"

func TestPublicInitializeIsExposed(t *testing.T) {
	var _ = Initialize
	var _ = InitializeWithDependencies
	var _ = RegisterRoutes
}
