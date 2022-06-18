package service

var aws_pre_config string = ``

// 阿里前置配置：
/*
1. 使用data source查询可用区，规格，镜像和网络参数
2. 安全组规则等参数，复用已有的实例
*/
var ali_pre_config string = `data "alicloud_instances" "myali_instances" {
}`

// 华为前置配置：
/*
1. 使用data source查询可用区，规格，镜像和网络参数
2. 安全组规则等参数，复用已有的实例
*/
var hw_pre_config string = `data "huaweicloud_vpc_subnet" "net" {
	name = "subnet-default"
}  
data "huaweicloud_compute_instances" "myhw_instances" {
}`
