/*
 * Copyright (c) 2025, WSO2 LLC. (https://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package executor

import (
	"fmt"

	"github.com/thunder-id/thunderid/internal/attributecache"
	"github.com/thunder-id/thunderid/internal/authn/assert"
	"github.com/thunder-id/thunderid/internal/authn/consent"
	"github.com/thunder-id/thunderid/internal/authn/github"
	"github.com/thunder-id/thunderid/internal/authn/google"
	"github.com/thunder-id/thunderid/internal/authn/magiclink"
	"github.com/thunder-id/thunderid/internal/authn/oauth"
	"github.com/thunder-id/thunderid/internal/authn/oidc"
	"github.com/thunder-id/thunderid/internal/authn/otp"
	"github.com/thunder-id/thunderid/internal/authn/passkey"
	authnprovidermgr "github.com/thunder-id/thunderid/internal/authnprovider/manager"
	"github.com/thunder-id/thunderid/internal/authz"
	"github.com/thunder-id/thunderid/internal/entityprovider"
	"github.com/thunder-id/thunderid/internal/entitytype"
	"github.com/thunder-id/thunderid/internal/flow/common"
	"github.com/thunder-id/thunderid/internal/flow/core"
	"github.com/thunder-id/thunderid/internal/group"
	"github.com/thunder-id/thunderid/internal/idp"
	"github.com/thunder-id/thunderid/internal/notification"
	"github.com/thunder-id/thunderid/internal/ou"
	"github.com/thunder-id/thunderid/internal/role"
	"github.com/thunder-id/thunderid/internal/system/email"
	"github.com/thunder-id/thunderid/internal/system/jose/jwt"
	"github.com/thunder-id/thunderid/internal/system/log"
	"github.com/thunder-id/thunderid/internal/system/template"
)

// RegisterDeps holds service dependencies required to construct built-in executors.
type RegisterDeps struct {
	FlowFactory           core.FlowFactoryInterface
	OUService             ou.OrganizationUnitServiceInterface
	IDPService            idp.IDPServiceInterface
	NotifSenderSvc        notification.NotificationSenderServiceInterface
	JWTService            jwt.JWTServiceInterface
	AuthAssertGen         assert.AuthAssertGeneratorInterface
	ConsentEnforcer       consent.ConsentEnforcerServiceInterface
	AuthnProvider         authnprovidermgr.AuthnProviderManagerInterface
	OTPService            otp.OTPAuthnServiceInterface
	PasskeyService        passkey.PasskeyServiceInterface
	MagicLinkService      magiclink.MagicLinkAuthnServiceInterface
	AuthZService          authz.AuthorizationServiceInterface
	EntityTypeService     entitytype.EntityTypeServiceInterface
	GroupService          group.GroupServiceInterface
	RoleService           role.RoleServiceInterface
	RoleAssignmentService role.RoleAssignmentServiceInterface
	EntityProvider        entityprovider.EntityProviderInterface
	AttributeCacheSvc     attributecache.AttributeCacheServiceInterface
	EmailClient           email.EmailClientInterface
	TemplateService       template.TemplateServiceInterface
	OAuthSvc              oauth.OAuthAuthnServiceInterface
	OIDCSvc               oidc.OIDCAuthnServiceInterface
	GithubSvc             github.GithubOAuthAuthnServiceInterface
	GoogleSvc             google.GoogleOIDCAuthnServiceInterface
}

type builtInExecutorRegistrar func(ExecutorRegistryInterface, RegisterDeps) error

// builtInExecutorNames is the canonical ordered list of built-in executors.
// When adding an executor: define ExecutorName* in constants.go, append here, and add
// a matching entry in builtInExecutorRegistrars.
var builtInExecutorNames = []string{
	ExecutorNameBasicAuth,
	ExecutorNameSMSAuth,
	ExecutorNamePasskeyAuth,
	ExecutorNameMagicLinkAuth,
	ExecutorNameOAuth,
	ExecutorNameOIDCAuth,
	ExecutorNameGitHubAuth,
	ExecutorNameGoogleAuth,
	ExecutorNameProvisioning,
	ExecutorNameOUCreation,
	ExecutorNameAttributeCollect,
	ExecutorNameAuthAssert,
	ExecutorNameAuthorization,
	ExecutorNameHTTPRequest,
	ExecutorNameUserTypeResolver,
	ExecutorNameInviteExecutor,
	ExecutorNameEmailExecutor,
	ExecutorNameCredentialSetter,
	ExecutorNamePermissionValidator,
	ExecutorNameIdentifying,
	ExecutorNameConsent,
	ExecutorNameOUResolver,
	ExecutorNameAttributeUniquenessValidator,
	ExecutorNameSMSExecutor,
	ExecutorNameFederatedAuthResolver,
}

//nolint:gochecknoglobals // package-level catalog validated in init
var builtInExecutorRegistrars = map[string]builtInExecutorRegistrar{
	ExecutorNameBasicAuth: func(reg ExecutorRegistryInterface, deps RegisterDeps) error {
		reg.RegisterExecutor(ExecutorNameBasicAuth, newBasicAuthExecutor(
			deps.FlowFactory, deps.EntityProvider, deps.AuthnProvider))
		return nil
	},
	ExecutorNameSMSAuth: func(reg ExecutorRegistryInterface, deps RegisterDeps) error {
		reg.RegisterExecutor(ExecutorNameSMSAuth, newSMSOTPAuthExecutor(
			deps.FlowFactory, deps.OTPService, deps.AuthnProvider, deps.EntityProvider))
		return nil
	},
	ExecutorNamePasskeyAuth: func(reg ExecutorRegistryInterface, deps RegisterDeps) error {
		reg.RegisterExecutor(ExecutorNamePasskeyAuth, newPasskeyAuthExecutor(
			deps.FlowFactory, deps.PasskeyService, deps.AuthnProvider, deps.EntityProvider))
		return nil
	},
	ExecutorNameMagicLinkAuth: func(reg ExecutorRegistryInterface, deps RegisterDeps) error {
		reg.RegisterExecutor(ExecutorNameMagicLinkAuth, newMagicLinkAuthExecutor(
			deps.FlowFactory, deps.MagicLinkService, deps.AuthnProvider, deps.EntityProvider))
		return nil
	},
	ExecutorNameOAuth: func(reg ExecutorRegistryInterface, deps RegisterDeps) error {
		reg.RegisterExecutor(ExecutorNameOAuth, newOAuthExecutor(
			"", []common.Input{}, []common.Input{}, deps.FlowFactory, deps.IDPService, deps.EntityTypeService,
			deps.OAuthSvc, deps.AuthnProvider, idp.IDPTypeOAuth))
		return nil
	},
	ExecutorNameOIDCAuth: func(reg ExecutorRegistryInterface, deps RegisterDeps) error {
		reg.RegisterExecutor(ExecutorNameOIDCAuth, newOIDCAuthExecutor(
			"", []common.Input{}, []common.Input{}, deps.FlowFactory, deps.IDPService, deps.EntityTypeService,
			deps.OIDCSvc, deps.AuthnProvider, idp.IDPTypeOIDC))
		return nil
	},
	ExecutorNameGitHubAuth: func(reg ExecutorRegistryInterface, deps RegisterDeps) error {
		reg.RegisterExecutor(ExecutorNameGitHubAuth, newGithubOAuthExecutor(
			deps.FlowFactory, deps.IDPService, deps.EntityTypeService, deps.GithubSvc, deps.AuthnProvider))
		return nil
	},
	ExecutorNameGoogleAuth: func(reg ExecutorRegistryInterface, deps RegisterDeps) error {
		reg.RegisterExecutor(ExecutorNameGoogleAuth, newGoogleOIDCAuthExecutor(
			deps.FlowFactory, deps.IDPService, deps.EntityTypeService, deps.GoogleSvc, deps.AuthnProvider))
		return nil
	},
	ExecutorNameProvisioning: func(reg ExecutorRegistryInterface, deps RegisterDeps) error {
		reg.RegisterExecutor(ExecutorNameProvisioning, newProvisioningExecutor(
			deps.FlowFactory, deps.GroupService, deps.RoleService, deps.RoleAssignmentService,
			deps.EntityProvider, deps.EntityTypeService))
		return nil
	},
	ExecutorNameOUCreation: func(reg ExecutorRegistryInterface, deps RegisterDeps) error {
		reg.RegisterExecutor(ExecutorNameOUCreation, newOUExecutor(deps.FlowFactory, deps.OUService))
		return nil
	},
	ExecutorNameAttributeCollect: func(reg ExecutorRegistryInterface, deps RegisterDeps) error {
		reg.RegisterExecutor(ExecutorNameAttributeCollect, newAttributeCollector(deps.FlowFactory, deps.EntityProvider))
		return nil
	},
	ExecutorNameAuthAssert: func(reg ExecutorRegistryInterface, deps RegisterDeps) error {
		reg.RegisterExecutor(ExecutorNameAuthAssert, newAuthAssertExecutor(deps.FlowFactory, deps.JWTService,
			deps.OUService, deps.AuthAssertGen, deps.AuthnProvider, deps.EntityProvider,
			deps.AttributeCacheSvc, deps.RoleService))
		return nil
	},
	ExecutorNameAuthorization: func(reg ExecutorRegistryInterface, deps RegisterDeps) error {
		reg.RegisterExecutor(ExecutorNameAuthorization, newAuthorizationExecutor(
			deps.FlowFactory, deps.AuthZService, deps.EntityProvider))
		return nil
	},
	ExecutorNameHTTPRequest: func(reg ExecutorRegistryInterface, deps RegisterDeps) error {
		reg.RegisterExecutor(ExecutorNameHTTPRequest, newHTTPRequestExecutor(deps.FlowFactory, deps.OUService))
		return nil
	},
	ExecutorNameUserTypeResolver: func(reg ExecutorRegistryInterface, deps RegisterDeps) error {
		reg.RegisterExecutor(ExecutorNameUserTypeResolver, newUserTypeResolver(
			deps.FlowFactory, deps.EntityTypeService, deps.OUService))
		return nil
	},
	ExecutorNameInviteExecutor: func(reg ExecutorRegistryInterface, deps RegisterDeps) error {
		reg.RegisterExecutor(ExecutorNameInviteExecutor, newInviteExecutor(deps.FlowFactory))
		return nil
	},
	ExecutorNameEmailExecutor: func(reg ExecutorRegistryInterface, deps RegisterDeps) error {
		reg.RegisterExecutor(ExecutorNameEmailExecutor, newEmailExecutor(
			deps.FlowFactory, deps.EmailClient, deps.TemplateService, deps.EntityProvider))
		return nil
	},
	ExecutorNameCredentialSetter: func(reg ExecutorRegistryInterface, deps RegisterDeps) error {
		reg.RegisterExecutor(ExecutorNameCredentialSetter, newCredentialSetter(deps.FlowFactory, deps.EntityProvider))
		return nil
	},
	ExecutorNamePermissionValidator: func(reg ExecutorRegistryInterface, deps RegisterDeps) error {
		reg.RegisterExecutor(ExecutorNamePermissionValidator, newPermissionValidator(deps.FlowFactory))
		return nil
	},
	ExecutorNameIdentifying: func(reg ExecutorRegistryInterface, deps RegisterDeps) error {
		reg.RegisterExecutor(ExecutorNameIdentifying, newIdentifyingExecutor(
			"", []common.Input{{Identifier: userAttributeUsername, Type: "string", Required: true}}, []common.Input{},
			deps.FlowFactory, deps.EntityProvider))
		return nil
	},
	ExecutorNameConsent: func(reg ExecutorRegistryInterface, deps RegisterDeps) error {
		reg.RegisterExecutor(ExecutorNameConsent, newConsentExecutor(
			deps.FlowFactory, deps.ConsentEnforcer, deps.AuthnProvider))
		return nil
	},
	ExecutorNameOUResolver: func(reg ExecutorRegistryInterface, deps RegisterDeps) error {
		reg.RegisterExecutor(ExecutorNameOUResolver, newOUResolverExecutor(deps.FlowFactory, deps.OUService))
		return nil
	},
	ExecutorNameAttributeUniquenessValidator: func(reg ExecutorRegistryInterface, deps RegisterDeps) error {
		reg.RegisterExecutor(ExecutorNameAttributeUniquenessValidator, newAttributeUniquenessValidator(
			deps.FlowFactory, deps.EntityTypeService, deps.EntityProvider))
		return nil
	},
	ExecutorNameSMSExecutor: func(reg ExecutorRegistryInterface, deps RegisterDeps) error {
		reg.RegisterExecutor(ExecutorNameSMSExecutor, newSMSExecutor(
			deps.FlowFactory, deps.NotifSenderSvc, deps.TemplateService))
		return nil
	},
	ExecutorNameFederatedAuthResolver: func(reg ExecutorRegistryInterface, deps RegisterDeps) error {
		reg.RegisterExecutor(ExecutorNameFederatedAuthResolver, newFederatedAuthResolverExecutor(deps.FlowFactory))
		return nil
	},
}

var builtInExecutorNameSet = func() map[string]struct{} {
	set := make(map[string]struct{}, len(builtInExecutorNames))
	for _, name := range builtInExecutorNames {
		set[name] = struct{}{}
	}
	return set
}()

// A broken catalog panics and the process never starts.
// It does not run on every RegisterExecutor call; only once per process when the package loads.
func init() {
	if err := validateBuiltInExecutorCatalog(); err != nil {
		panic("invalid built-in executor catalog: " + err.Error())
	}
}

func validateBuiltInExecutorCatalog() error {
	if len(builtInExecutorNames) != len(builtInExecutorRegistrars) {
		return fmt.Errorf("builtInExecutorNames has %d entries but builtInExecutorRegistrars has %d",
			len(builtInExecutorNames), len(builtInExecutorRegistrars))
	}
	seen := make(map[string]struct{}, len(builtInExecutorNames))
	for _, name := range builtInExecutorNames {
		if _, ok := builtInExecutorRegistrars[name]; !ok {
			return fmt.Errorf("builtInExecutorNames includes %q with no registrar", name)
		}
		if _, dup := seen[name]; dup {
			return fmt.Errorf("duplicate name %q in builtInExecutorNames", name)
		}
		seen[name] = struct{}{}
	}
	for name := range builtInExecutorRegistrars {
		if _, ok := seen[name]; !ok {
			return fmt.Errorf("builtInExecutorRegistrars includes %q missing from builtInExecutorNames", name)
		}
	}
	return nil
}

// defaultBuiltInExecutorNames returns the names of all built-in executors.
func defaultBuiltInExecutorNames() []string {
	names := make([]string, len(builtInExecutorNames))
	copy(names, builtInExecutorNames)
	return names
}

// registerBuiltInExecutors registers the requested built-in executors on reg.
// When names is empty, all built-in executors are registered.
func registerBuiltInExecutors(reg ExecutorRegistryInterface, deps RegisterDeps, names []string) error {
	resolved, err := resolveBuiltInExecutorNames(names)
	if err != nil {
		return err
	}
	for _, name := range resolved {
		if err := registerBuiltInExecutor(reg, deps, name); err != nil {
			return err
		}
	}
	logger := log.GetLogger().With(log.String(log.LoggerKeyComponentName, "ExecutorRegistry"))
	logger.Info("Registered built-in flow executors",
		log.Int("count", len(resolved)),
		log.Any("executors", resolved))
	return nil
}

func resolveBuiltInExecutorNames(names []string) ([]string, error) {
	if len(names) == 0 {
		return defaultBuiltInExecutorNames(), nil
	}
	names = dedupeExecutorNames(names)
	for _, name := range names {
		if _, ok := builtInExecutorNameSet[name]; !ok {
			return nil, fmt.Errorf("unknown built-in executor: %q", name)
		}
	}
	return names, nil
}

func dedupeExecutorNames(names []string) []string {
	seen := make(map[string]struct{}, len(names))
	out := make([]string, 0, len(names))
	for _, name := range names {
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		out = append(out, name)
	}
	return out
}

func registerBuiltInExecutor(reg ExecutorRegistryInterface, deps RegisterDeps, name string) error {
	register, ok := builtInExecutorRegistrars[name]
	if !ok {
		return fmt.Errorf("unhandled built-in executor: %q", name)
	}
	if err := register(reg, deps); err != nil {
		return err
	}
	if !reg.IsRegistered(name) {
		return fmt.Errorf("failed to register built-in executor: %q", name)
	}
	return nil
}
