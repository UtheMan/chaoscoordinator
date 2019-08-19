package azure

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

const (
	AzurePublicCloud = "AzurePublicCloud"
)

type AuthContext interface {
	Cloud() string
	ClientID() string
	ClientSecret() string
	TenantID() string
}

// authContext contains the necessary parameters for authentication to Azure resources.
type authContext struct {
	AzureCloud        string `json:"-"`
	AzureClientID     string `json:"clientId"`
	AzureClientSecret string `json:"clientSecret"`
	AzureTenantID     string `json:"tenantId"`
}

func InjectAuthorizer() (autorest.Authorizer, error) {
	config, err := provideConfiguration()
	if err != nil {
		return nil, err
	}
	return provideAuthorizer(config)
}

func provideConfiguration() (*authContext, error) {
	file, err := auth.GetSettingsFromFile()
	if err != nil {
		env, err := auth.GetSettingsFromEnvironment()
		if err != nil {
			return &authContext{}, err
		}
		config, err := env.GetClientCredentials()
		if err != nil {
			config, err = getMSICredentials()
			if err != nil {
				return &authContext{}, err
			}
		}
		return &authContext{
			AzureClientID:     config.ClientID,
			AzureClientSecret: config.ClientSecret,
			AzureTenantID:     config.TenantID,
			AzureCloud:        env.Environment.Name,
		}, nil
	}
	return &authContext{
		AzureClientID:     file.Values[auth.ClientID],
		AzureClientSecret: file.Values[auth.ClientSecret],
		AzureTenantID:     file.Values[auth.TenantID],
		AzureCloud:        AzurePublicCloud,
	}, nil
}

func provideAuthorizer(ac *authContext) (autorest.Authorizer, error) {
	env, err := azure.EnvironmentFromName(ac.AzureCloud)
	if err != nil {
		return nil, err
	}
	return provideResourceAuthorizer(env.ResourceManagerEndpoint)
}

func provideResourceAuthorizer(resource string) (autorest.Authorizer, error) {
	authorizer, err := auth.NewAuthorizerFromFileWithResource(resource)
	if err != nil {
		return auth.NewAuthorizerFromEnvironmentWithResource(resource)
	}
	return authorizer, nil
}

func getMSICredentials() (auth.ClientCredentialsConfig, error) {
	r, err := requestIdentityToken()
	if err != nil {
		return auth.ClientCredentialsConfig{}, err
	}
	return auth.ClientCredentialsConfig{ClientID: r.ClientID, TenantID: ""}, nil
}

type responseToken struct {
	AccessToken  string `json:"access_token"`
	ClientID     string `json:"client_id"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    string `json:"expires_in"`
	ExpiresOn    string `json:"expires_on"`
	NotBefore    string `json:"not_before"`
	Resource     string `json:"resource"`
	TokenType    string `json:"token_type"`
}

func requestIdentityToken() (*responseToken, error) {
	endpoint, err := url.Parse("http://169.254.169.254/metadata/identity/oauth2/token")
	if err != nil {
		return nil, err
	}

	parameters := url.Values{}
	parameters.Add("resource", "https://management.azure.com/")
	parameters.Add("api-version", "2018-02-01")
	endpoint.RawQuery = parameters.Encode()
	req, err := http.NewRequest("GET", endpoint.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Metadata", "true")

	client := &http.Client{}
	client.Timeout = 5 * time.Second
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	responseBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	var r responseToken
	err = json.Unmarshal(responseBytes, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}
