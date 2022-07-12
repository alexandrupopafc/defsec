package api_gateway

import (
	"fmt"

	v2 "github.com/aquasecurity/defsec/pkg/providers/aws/apigateway/v2"

	"github.com/aquasecurity/defsec/internal/types"
	api "github.com/aws/aws-sdk-go-v2/service/apigatewayv2"
	agTypes "github.com/aws/aws-sdk-go-v2/service/apigatewayv2/types"
)

func (a *adapter) getDomainNamesV2() ([]v2.DomainName, error) {

	var adapted []v2.DomainName

	a.Tracker().SetServiceLabel("Discovering v2 domain names...")

	var input api.GetDomainNamesInput
	var apiDomainNames []agTypes.DomainName
	for {
		output, err := a.clientV2.GetDomainNames(a.Context(), &input)
		if err != nil {
			return nil, err
		}
		apiDomainNames = append(apiDomainNames, output.Items...)
		a.Tracker().SetTotalResources(len(apiDomainNames))
		if output.NextToken == nil {
			break
		}
		input.NextToken = output.NextToken
	}

	a.Tracker().SetServiceLabel("Adapting v2 domain names...")

	for _, apiDomain := range apiDomainNames {
		adapted = append(adapted, a.adaptDomainNameV2(apiDomain))
		a.Tracker().IncrementResource()
	}

	return adapted, nil

}

func (a *adapter) adaptDomainNameV2(domain agTypes.DomainName) v2.DomainName {
	metadata := a.CreateMetadata(fmt.Sprintf("/domainnames/%s", *domain.DomainName))
	securityPolicy := "TLS_1_2"
	for _, policy := range domain.DomainNameConfigurations {
		if string(policy.SecurityPolicy) != "TLS_1_2" {
			securityPolicy = string(policy.SecurityPolicy)
		}
	}
	return v2.DomainName{
		Metadata:       metadata,
		Name:           types.String(*domain.DomainName, metadata),
		SecurityPolicy: types.String(securityPolicy, metadata),
	}
}
