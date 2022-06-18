package api

import (
	service "ServiceManagement/pkg"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	_ "github.com/aws/aws-sdk-go-v2/aws"
	_ "github.com/aws/aws-sdk-go-v2/config"
	_ "github.com/aws/aws-sdk-go-v2/service/dynamodb"

	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	ecs20140526 "github.com/alibabacloud-go/ecs-20140526/v2/client"
	"github.com/alibabacloud-go/tea/tea"

	core "ServiceManagement/pkg/hw-ak-sk/core"
)

// 调用AWS云OPEN API查看镜像参数 NovaShowFlavor
func GetAwsFlavorDetails(instance_type, region_id string) error {
	key_pair, err := service.GetKey_db(1)
	key, secret := key_pair.Aws_a_key, key_pair.Aws_s_key
	s := core.Signer{
		Key:    key,
		Secret: secret,
	}

	r, _ := http.NewRequest("GET", "https://ecs."+region_id+".myhuaweicloud.com/v2.1/47e2aa1bb84e4cb385270f9a87303990/flavors/"+instance_type, ioutil.NopCloser(bytes.NewBuffer([]byte(""))))

	s.Sign(r)
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return errors.New(err.Error())
	}
	defer resp.Body.Close()
	// 200 OK
	fmt.Println("status: ", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("body:\n", string(body))
	return nil
}

// 调用阿里云OPEN API查看镜像参数 DescribeInstanceTypes
func GetAliFlavorDetails(instance_type, region_id string) error {
	key_pair, err := service.GetKey_db(1)
	key, secret := key_pair.Ali_a_key, key_pair.Ali_s_key

	config := &openapi.Config{
		AccessKeyId:     tea.String(key),
		AccessKeySecret: tea.String(secret),
	}
	config.Endpoint = tea.String("ecs." + region_id + ".aliyuncs.com")
	client, err := ecs20140526.NewClient(config)

	params := []*string{&instance_type}

	describeInstanceTypesRequest := &ecs20140526.DescribeInstanceTypesRequest{
		InstanceTypes: params,
	}
	resp, err := client.DescribeInstanceTypes(describeInstanceTypesRequest)
	if err != nil {
		return err
	}
	fmt.Println("resp.Body: ", resp.Body)
	fmt.Println("resp.Body type: ", reflect.TypeOf(resp.Body))

	return nil
}

// 调用华为云OPEN API查看镜像参数 NovaShowFlavor
func GetHwFlavorDetails(instance_type, region_id string) error {
	key_pair, err := service.GetKey_db(1)
	key, secret := key_pair.Hw_a_key, key_pair.Hw_s_key
	s := core.Signer{
		Key:    key,
		Secret: secret,
	}

	r, _ := http.NewRequest("GET", "https://ecs."+region_id+".myhuaweicloud.com/v2.1/47e2aa1bb84e4cb385270f9a87303990/flavors/"+instance_type, ioutil.NopCloser(bytes.NewBuffer([]byte(""))))

	s.Sign(r)
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return errors.New(err.Error())
	}
	defer resp.Body.Close()
	// 200 OK
	fmt.Println("status: ", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("body:\n", string(body))
	return nil
}
