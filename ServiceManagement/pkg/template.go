package service

// region 引入 done
var default_main_aws string = `terraform {
	required_providers {
	  aws = {
		source  = "hashicorp/aws" 
	  }
	}
}
provider "aws" {
	region = "${var.region}"
	access_key = "${var.access_key}"
	secret_key = "${var.secret_key}"
}`
var default_main_ali string = `terraform {
	required_providers {
	  alicoud = {
		  source = "aliyun/alicloud"
	  }
	}
}
provider "alicloud" {
	region     = "${var.region}"
	access_key = "${var.access_key}"
	secret_key = "${var.secret_key}"
}`
var default_main_hw string = `terraform {
	required_providers {
	  huaweicloud = {
		  source = "huaweicloud/huaweicloud"
	  }
	}
}
provider "huaweicloud" {
	region     = "${var.region}"
	access_key = "${var.access_key}"
	secret_key = "${var.secret_key}"
}`
var default_var string = `variable "access_key" {
    type = string
    default = ""
  }
variable "secret_key" {
    type = string
    default = ""
}
variable "region" {
    type = string
    default = ""
}`
var default_tfvars string = `{
    "access_key": "",
    "secret_key": "",
	"region": ""
}`
