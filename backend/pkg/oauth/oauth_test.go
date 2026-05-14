package oauth

import "testing"

func TestPublicInitializeIsExposed(t *testing.T) {
	var _ = InitializeWithDependencies
	var _ = RegisterRoutes
}
