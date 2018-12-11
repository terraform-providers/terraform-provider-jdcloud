package jdcloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	vpcApis "github.com/jdcloud-api/jdcloud-sdk-go/services/vpc/apis"
	vpcClient "github.com/jdcloud-api/jdcloud-sdk-go/services/vpc/client"
	"github.com/satori/go.uuid"
	"strings"
)

func validateStringInSlice(validSlice []string, ignoreCase bool) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, err []error) {
		v, ok := i.(string)
		if !ok {
			err = append(err, fmt.Errorf("type of %s must be string", k))
			return
		}

		for _, str := range validSlice {
			if v == str || (ignoreCase && strings.ToLower(v) == strings.ToLower(str)) {
				return
			}
		}

		err = append(err, fmt.Errorf("expected %s to be one of %v, got %s", k, validSlice, v))
		return
	}
}

func validateStringNoEmpty(v interface{}, k string) (s []string, errs []error) {

	value := v.(string)
	if len(value) < 1 {
		errs = append(errs, fmt.Errorf("%v can not be Empty characters.", k))
	}

	return
}

func diskClientTokenDefault() string {
	var clientToken string
	nonce, _ := uuid.NewV4()
	clientToken = nonce.String()
	clientToken = strings.Replace(clientToken, "-", "", -1)

	if len(clientToken) > 20 {
		clientToken = string([]byte(clientToken)[:20])
	}
	return clientToken
}

func verifyVPC(d *schema.ResourceData, m interface{}, vpc, subnet string) error {

	config := m.(*JDCloudConfig)
	subnetClient := vpcClient.NewVpcClient(config.Credential)

	req := vpcApis.NewDescribeSubnetRequest(config.Region, subnet)
	resp, err := subnetClient.DescribeSubnet(req)

	if err != nil {
		return fmt.Errorf("[ERROR] verifyVPC Failed, when creating RDS %s ", err.Error())
	}

	if resp.Error.Code != 0 {
		return fmt.Errorf("[ERROR] verifyVPC Failed, when creating RDS  code:%d staus:%s message:%s ", resp.Error.Code, resp.Error.Status, resp.Error.Message)
	}

	if resp.Result.Subnet.VpcId != vpc {
		return fmt.Errorf("[ERROR] verifyVPC Failed, vpc ID does not match")
	}
	return nil
}
