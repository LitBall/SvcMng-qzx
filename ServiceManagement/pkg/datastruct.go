package service

type User struct {
	Uid       int    `db:"uid"`
	Username  string `db:"username"`
	Password  string `db:"password"`
	Email     string `db:"email"`
	UserGroup string `db:"usergroup"`
}
type Resource struct {
	Rid          int    `db:"rid"`
	ResourceType string `db:"resource_type"`
	// SystemType   string `db:"system_type"`
	// VCPU         int    `db:"vCPU"`
	// Memory       string `db:"memory"`
	Bandwidth string `db:"bandwidth"`
	NumberOf  int    `db:"number_of"`
	LastTime  int    `db:"last_time"`
	// Status       int    `db:"status"`
	ImageId      string `db:"image_id"`
	InstanceType string `db:"instance_type"`
	InstanceName string `db:"instance_name"`
	RegionId     string `db:"region_id"`
}
type Template struct {
	Tid          int    `db:"tid"`
	ResourceType string `db:"resource_type"`
	// SystemType   string `db:"system_type"`
	// VCPU         int    `db:"vCPU"`
	// Memory       string `db:"memory"`
	Bandwidth    string `db:"bandwidth"`
	NumberOf     int    `db:"number_of"`
	LastTime     int    `db:"last_time"`
	ImageId      string `db:"image_id"`
	InstanceType string `db:"instance_type"`
	// InstanceName string `db:"instance_name"`
	// RegionId     string `db:"region_id"`
}
type AccessKey struct {
	Uid        int    `db:"uid"`
	Aws_a_key  string `db:"aws_a_key"`
	Aws_s_key  string `db:"aws_s_key"`
	Aws_region string `db:"aws_region"`
	Ali_a_key  string `db:"ali_a_key"`
	Ali_s_key  string `db:"ali_s_key"`
	Ali_region string `db:"ali_region"`
	Hw_a_key   string `db:"hw_a_key"`
	Hw_s_key   string `db:"hw_s_key"`
	Hw_region  string `db:"hw_region"`
}
type InstanceConfig struct {
	ResourceType string            `json:"resource_type"`
	Bandwidth    string            `json:"bandwidth"`
	NumberOf     string            `json:"number_of"`
	LastTime     string            `json:"last_time"`
	ImageId      string            `json:"image_id"`
	InstanceType string            `json:"instance_type"` // -->hw:flavor_id
	InstanceName string            `json:"instance_name"`
	RegionId     string            `json:"region_id"`
	HW_NetWork   map[string]string `json:"hw_network"`
	Description  string            `json:"description"`
	TagsName     string            `json:"tags_name"`
}

type Argument struct {
	name string
}
type TfVars struct {
	Aws_a_key  string `json:"aws_a_key"`
	Aws_s_key  string `json:"aws_s_key"`
	Aws_region string `json:"aws_region"`
	Ali_a_key  string `json:"ali_a_key"`
	Ali_s_key  string `json:"ali_s_key"`
	Ali_region string `json:"ali_region"`
	Hw_a_key   string `json:"hw_a_key"`
	Hw_s_key   string `json:"hw_s_key"`
	Hw_region  string `json:"hw_region"`
}

type TemplateResponse struct {
	Code int                   `json:"code"`
	Msg  string                `json:"msg"`
	Data map[string][]Template `json:"data"`
}
type ASkeyResponse struct {
	Code int                  `json:"code"`
	Msg  string               `json:"msg"`
	Data map[string]AccessKey `json:"data"`
}
type InstanceResponse struct {
	Code int      `json:"code"`
	Msg  string   `json:"msg"`
	Data Resource `json:"data"` // 测试使用aahahahahahah
}

type CloudResponse struct {
	Headers int    `json:"headers"`
	Body    string `json:"body"`
}
