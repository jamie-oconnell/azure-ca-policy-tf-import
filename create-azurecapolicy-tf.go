package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/zclconf/go-cty/cty"
)

func create_azurecapolicy(policy models.ConditionalAccessPolicy) {
	// create new empty hcl file object
	f := hclwrite.NewEmptyFile()

	tfFile, err := os.Create(fmt.Sprintf("generated/%s.tf", *policy.GetDisplayName()))
	if err != nil {
		fmt.Println(err)
		return
	}
	// initialize the body of the new file object
	rootBody := f.Body()

	// Create Azure AD Conditional Access Policy resource block
	resource_name := *policy.GetDisplayName()
	azureADPolicy := rootBody.AppendNewBlock("resource", []string{"azuread_conditional_access_policy", strings.ToLower(strings.ReplaceAll(resource_name, " ", "_"))})
	azureADPolicyBody := azureADPolicy.Body()

	// Set attributes for Azure AD Conditional Access Policy

	azureADPolicyBody.SetAttributeValue("display_name", cty.StringVal(*policy.GetDisplayName()))
	azureADPolicyBody.SetAttributeValue("state", cty.StringVal(policy.GetState().String()))

	azureADPolicyBody.AppendNewline()

	// Add conditions block
	conditionsBlock := azureADPolicyBody.AppendNewBlock("conditions", nil)
	conditionsBlockBody := conditionsBlock.Body()

	// Set conditions attributes
	var clientAppTypeStrings []string
	for _, app := range policy.GetConditions().GetClientAppTypes() {
		clientAppTypeStrings = append(clientAppTypeStrings, app.String())
	}
	setIfNotEmpty(conditionsBlockBody, "client_app_types", clientAppTypeStrings)
	var signInRiskLevelStrings []string
	for _, risk := range policy.GetConditions().GetSignInRiskLevels() {
		signInRiskLevelStrings = append(signInRiskLevelStrings, risk.String())
	}
	setIfNotEmpty(conditionsBlockBody, "sign_in_risk_levels", signInRiskLevelStrings)
	var userRiskLevelStrings []string
	for _, risk := range policy.GetConditions().GetUserRiskLevels() {
		userRiskLevelStrings = append(userRiskLevelStrings, risk.String())
	}
	setIfNotEmpty(conditionsBlockBody, "user_risk_levels", userRiskLevelStrings)
	var servicePrincipalRiskLevelStrings []string
	for _, risk := range policy.GetConditions().GetServicePrincipalRiskLevels() {
		servicePrincipalRiskLevelStrings = append(servicePrincipalRiskLevelStrings, risk.String())
	}
	setIfNotEmpty(conditionsBlockBody, "service_principal_risk_levels", servicePrincipalRiskLevelStrings)

	conditionsBlockBody.AppendNewline()

	// Add applications block
	if policy.GetConditions().GetApplications() != nil {
		appsBlock := conditionsBlockBody.AppendNewBlock("applications", nil)
		appsBlockBody := appsBlock.Body()
		setIfNotEmpty(appsBlockBody, "included_applications", policy.GetConditions().GetApplications().GetIncludeApplications())
		setIfNotEmpty(appsBlockBody, "excluded_applications", policy.GetConditions().GetApplications().GetExcludeApplications())
		setIfNotEmpty(appsBlockBody, "included_user_actions", policy.GetConditions().GetApplications().GetIncludeUserActions())
		conditionsBlockBody.AppendNewline()
	}

	// Add client_applications block
	if policy.GetConditions().GetClientApplications() != nil {
		clientAppsBlock := conditionsBlockBody.AppendNewBlock("client_applications", nil)
		clientAppsBlockBody := clientAppsBlock.Body()
		setIfNotEmpty(clientAppsBlockBody, "included_service_principals", policy.GetConditions().GetClientApplications().GetIncludeServicePrincipals())
		setIfNotEmpty(clientAppsBlockBody, "excluded_service_principals", policy.GetConditions().GetClientApplications().GetExcludeServicePrincipals())
		conditionsBlockBody.AppendNewline()
	}

	// Add devices block
	if policy.GetConditions().GetDevices() != nil {
		devicesBlock := conditionsBlockBody.AppendNewBlock("devices", nil)
		devicesBlockBody := devicesBlock.Body()
		filterBlock := devicesBlockBody.AppendNewBlock("filter", nil)
		filterBlockBody := filterBlock.Body()

		if policy.GetConditions().GetDevices().GetDeviceFilter() != nil {
			filterBlockBody.SetAttributeValue("mode", cty.StringVal(policy.GetConditions().GetDevices().GetDeviceFilter().GetMode().String()))
			filterBlockBody.SetAttributeValue("rule", cty.StringVal(*policy.GetConditions().GetDevices().GetDeviceFilter().GetRule()))
		}
		conditionsBlockBody.AppendNewline()
	}

	// Add locations block
	if policy.GetConditions() != nil && policy.GetConditions().GetLocations() != nil {
		locationsBlock := conditionsBlockBody.AppendNewBlock("locations", nil)
		locationsBlockBody := locationsBlock.Body()
		setIfNotEmpty(locationsBlockBody, "included_locations", policy.GetConditions().GetLocations().GetIncludeLocations())
		setIfNotEmpty(locationsBlockBody, "excluded_locations", policy.GetConditions().GetLocations().GetExcludeLocations())
		conditionsBlockBody.AppendNewline()
	}

	// Add platforms block
	if policy.GetConditions().GetPlatforms() != nil {
		platformsBlock := conditionsBlockBody.AppendNewBlock("platforms", nil)
		platformsBlockBody := platformsBlock.Body()
		var includePlatforms []string
		for _, platform := range policy.GetConditions().GetPlatforms().GetIncludePlatforms() {
			includePlatforms = append(includePlatforms, platform.String())
		}
		setIfNotEmpty(platformsBlockBody, "included_platforms", includePlatforms)
		var excludePlatforms []string
		for _, platform := range policy.GetConditions().GetPlatforms().GetExcludePlatforms() {
			excludePlatforms = append(excludePlatforms, platform.String())
		}
		setIfNotEmpty(platformsBlockBody, "excluded_platforms", excludePlatforms)
	}

	// Add users block
	usersBlock := conditionsBlockBody.AppendNewBlock("users", nil)
	usersBlockBody := usersBlock.Body()
	setIfNotEmpty(usersBlockBody, "included_users", policy.GetConditions().GetUsers().GetIncludeUsers())
	setIfNotEmpty(usersBlockBody, "excluded_users", policy.GetConditions().GetUsers().GetExcludeUsers())
	setIfNotEmpty(usersBlockBody, "included_groups", policy.GetConditions().GetUsers().GetIncludeGroups())
	setIfNotEmpty(usersBlockBody, "excluded_groups", policy.GetConditions().GetUsers().GetExcludeGroups())
	setIfNotEmpty(usersBlockBody, "included_roles", policy.GetConditions().GetUsers().GetIncludeRoles())
	setIfNotEmpty(usersBlockBody, "excluded_roles", policy.GetConditions().GetUsers().GetExcludeRoles())

	azureADPolicyBody.AppendNewline()

	// Add grant_controls block
	if policy.GetGrantControls() != nil {
		grantControlsBlock := azureADPolicyBody.AppendNewBlock("grant_controls", nil)
		grantControlsBlockBody := grantControlsBlock.Body()
		grantControlsBlockBody.SetAttributeValue("operator", cty.StringVal(*policy.GetGrantControls().GetOperator()))
		var builtInControls []string
		for _, control := range policy.GetGrantControls().GetBuiltInControls() {
			builtInControls = append(builtInControls, control.String())
		}
		setIfNotEmpty(grantControlsBlockBody, "built_in_controls", builtInControls)
		setIfNotEmpty(grantControlsBlockBody, "custom_authentication_factors", policy.GetGrantControls().GetCustomAuthenticationFactors())
		setIfNotEmpty(grantControlsBlockBody, "terms_of_use", policy.GetGrantControls().GetTermsOfUse())
		if policy.GetGrantControls().GetAuthenticationStrength() != nil {
			grantControlsBlockBody.SetAttributeValue("authentication_strength_policy_id", cty.StringVal(*policy.GetGrantControls().GetAuthenticationStrength().GetId()))
		}
		azureADPolicyBody.AppendNewline()
	}

	// Add session_controls block
	if policy.GetSessionControls() != nil {
		sessionControlsBlock := azureADPolicyBody.AppendNewBlock("session_controls", nil)
		sessionControlsBlockBody := sessionControlsBlock.Body()
		if applicationEnforcedRestrictions := policy.GetSessionControls().GetApplicationEnforcedRestrictions(); applicationEnforcedRestrictions != nil {
			sessionControlsBlockBody.SetAttributeValue("application_enforced_restrictions_enabled", cty.BoolVal(*applicationEnforcedRestrictions.GetIsEnabled()))
		}
		if disableResilienceDefaults := policy.GetSessionControls().GetDisableResilienceDefaults(); disableResilienceDefaults != nil {
			sessionControlsBlockBody.SetAttributeValue("disable_resilience_defaults", cty.BoolVal(*disableResilienceDefaults))
		}
		if frequency := policy.GetSessionControls().GetSignInFrequency(); frequency != nil && frequency.GetValue() != nil {
			sessionControlsBlockBody.SetAttributeValue("sign_in_frequency", cty.NumberIntVal(int64(*frequency.GetValue())))
		}
		if frequencyPeriod := policy.GetSessionControls().GetSignInFrequency(); frequencyPeriod != nil && frequencyPeriod.GetTypeEscaped() != nil {
			sessionControlsBlockBody.SetAttributeValue("sign_in_frequency_period", cty.StringVal(frequencyPeriod.GetTypeEscaped().String()))
		}
		if frequencyAuthType := policy.GetSessionControls().GetSignInFrequency(); frequencyAuthType != nil && frequencyAuthType.GetAuthenticationType() != nil {
			sessionControlsBlockBody.SetAttributeValue("sign_in_frequency_authentication_type", cty.StringVal(frequencyAuthType.GetAuthenticationType().String()))
		}
		if frequencyInterval := policy.GetSessionControls().GetSignInFrequency(); frequencyInterval != nil && frequencyInterval.GetFrequencyInterval() != nil {
			sessionControlsBlockBody.SetAttributeValue("sign_in_frequency_interval", cty.StringVal(frequencyInterval.GetFrequencyInterval().String()))
		}
		if cloudAppSecurity := policy.GetSessionControls().GetCloudAppSecurity(); cloudAppSecurity != nil && cloudAppSecurity.GetCloudAppSecurityType() != nil {
			sessionControlsBlockBody.SetAttributeValue("cloud_app_security_policy", cty.StringVal(cloudAppSecurity.GetCloudAppSecurityType().String()))
		}
		if persistentBrowser := policy.GetSessionControls().GetPersistentBrowser(); persistentBrowser != nil {
			sessionControlsBlockBody.SetAttributeValue("persistent_browser_mode", cty.StringVal(persistentBrowser.GetMode().String()))
		}

	}

	fmt.Printf("Created terraform file for policy: %s \n", *policy.GetDisplayName())
	tfFile.Write(f.Bytes())
}

func setIfNotEmpty(body *hclwrite.Body, attributeName string, values []string) {
	if len(values) > 0 {
		attrList := make([]cty.Value, len(values))
		for i, v := range values {
			attrList[i] = cty.StringVal(v)
		}
		body.SetAttributeValue(attributeName, cty.ListVal(attrList))
	}
}
