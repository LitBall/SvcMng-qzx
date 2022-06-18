package api

/*
- 封装网络请求响应方法，实现http请求
*/
import (
	service "ServiceManagement/pkg"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "fmt"
	"log"
	"net/http"
	"strconv"
)

// 跨域中间件
func cors(w *http.ResponseWriter, r *http.Request) bool {
	(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	(*w).Header().Set("Access-Control-Allow-Headers", "DNT,X-Mx-ReqToken,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Authorization")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, TRACE, CONNECT, OPTIONS")
	(*w).Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")

	// 如果将要返回的是json
	(*w).Header().Set("Content-Type", "application/json")

	// 如果是跨域预检请求，那就别再继续执行。
	if r.Method == "OPTIONS" {
		(*w).WriteHeader(200)
		return true
	}
	return false
}

// GET请求，响应模板数据
func getTemplateListHandler(w http.ResponseWriter, req *http.Request) {
	template, err := service.LoadTemplate()
	if err != nil {
		fmt.Println("get failed: ", err)
		return
	}
	// 数据封装至json
	resp_str := service.TemplateResponse{
		Code: 1,
		Msg:  "模板加载成功~",
		Data: map[string][]service.Template{
			"template": template,
		},
	}
	resp_json, err := json.Marshal(resp_str)
	if err != nil {
		fmt.Println("marshal failed: ", err)
		return
	}
	// 设置请求头、写入响应
	if isPreflight := cors(&w, req); isPreflight == true {
		return
	}
	w.Write(resp_json)
}

// GET请求，查阅模板详情
func getTemplateDetailHandler(w http.ResponseWriter, req *http.Request) {
	// var (
	// 	code int
	// 	msg  string
	// 	data map[string]service.AccessKey
	// )
	// // 解析get query数据
	// query := req.URL.Query()
	// 	instance_type := query.Get("instance_type")
	// 	region_id := query.Get("region_id")

	// 	// 处理数据

	// 	code = 1
	// 	msg = "密钥更新成功~"
	// 	data = map[string]service.AccessKey{
	// 		"accesskey": key,
	// 	}
	// Response:
	// 	// 数据封装至json
	// 	resp_str := service.ASkeyResponse{
	// 		Code: code,
	// 		Msg:  msg,
	// 		Data: data,
	// 	}
	// 	resp_json, err := json.Marshal(resp_str)
	// 	if err != nil {
	// 		fmt.Println("marshal failed: ", err)
	// 		return
	// 	}
	// 	// 设置请求头、写入响应
	// 	if isPreflight := cors(&w, req); isPreflight == true {
	// 		return
	// 	}
	// 	w.Write(resp_json)
}

// POST请求，创建实例
func addCloudInstanceHandler(w http.ResponseWriter, req *http.Request) {
	var (
		resource service.Resource
		code     int
		msg      string
		para     service.InstanceConfig
	)
	// 设置请求头、写入响应
	if isPreflight := cors(&w, req); isPreflight == true {
		return
	}
	// 解析数据
	query := req.URL.Query()
	tid := query.Get("tid")
	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(&para); err != nil {
		fmt.Println("decode failed: ", err)
		return
	}
	if tid != "" {
		// 基于模板
		tid, err := strconv.Atoi(tid)
		fmt.Println("tid: ", tid)
		if err != nil {
			fmt.Println("strconv failed: ", err)
			return
		}
		// 处理数据
		res, err := service.CreateInstanceBaseTpl(tid, para)
		if err != nil {
			code = -1
			msg = "部署失败了T.T！"
			resource = service.Resource{}
			goto Response
		} else {
			code = 1
			msg = "实例部署成功~"
			resource = res
		}

	} else {
		// 处理数据
		res, err := service.CreateInstance(para)
		resource = res
		if err != nil {
			fmt.Println("create failed: ", err)
			code = -1
			msg = "部署失败了T.T！"
			resource = service.Resource{}
			goto Response
		} else {
			code = 1
			msg = "实例部署成功~"
			resource = res
		}
	}
	// goto-label: Response
Response:
	// 数据封装至json
	resp_str := service.InstanceResponse{
		Code: code,
		Msg:  msg,
		Data: resource,
	}
	resp_json, err := json.Marshal(resp_str)
	if err != nil {
		fmt.Println("marshal failed: ", err)
		return
	}

	w.Write(resp_json)
}

// GET请求，查阅当前信息密钥
func getKeyHandler(w http.ResponseWriter, req *http.Request) {
	var (
		code int
		msg  string
		key  service.AccessKey
		data map[string]service.AccessKey
	)
	query := req.URL.Query()
	uid_str := query.Get("uid")
	uid, err := strconv.Atoi(uid_str)
	fmt.Println("uid: ", uid)
	if err != nil {
		fmt.Println("strconv failed: ", err)
		code = -1
		msg = "密钥加载失败了T.T！"
		data = map[string]service.AccessKey{}
		goto Response
	}

	key, err = service.GetKey_db(uid)
	if err != nil {
		fmt.Println("get failed: ", err)
		code = -1
		data = map[string]service.AccessKey{}
		if err == sql.ErrNoRows {
			msg = "没有找到该用户的密钥！T.T"
		} else {
			msg = "密钥加载失败了T.T！"
		}
		goto Response
	}
	code = 1
	msg = "密钥加载成功~"
	data = map[string]service.AccessKey{
		"accesskey": key,
	}

Response:
	// 数据封装至json
	resp_str := service.ASkeyResponse{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	resp_json, err := json.Marshal(resp_str)
	if err != nil {
		fmt.Println("marshal failed: ", err)
		return
	}
	// 设置请求头、写入响应
	if isPreflight := cors(&w, req); isPreflight == true {
		return
	}
	w.Write(resp_json)
}

// POST请求，更新密钥
func updateKeyHandler(w http.ResponseWriter, req *http.Request) {
	var (
		para service.TfVars
		code int
		msg  string
		key  service.AccessKey
		data map[string]service.AccessKey
	)
	query := req.URL.Query()
	uid_str := query.Get("uid")
	uid, err := strconv.Atoi(uid_str)
	fmt.Println("uid: ", uid)
	if err != nil {
		fmt.Println("strconv failed: ", err)
		code = -1
		msg = "密钥更新失败了T.T！"
		data = map[string]service.AccessKey{}
		goto Response
	} else {
		// 解析post数据
		decoder := json.NewDecoder(req.Body)
		if err := decoder.Decode(&para); err != nil {
			fmt.Println("decode failed: ", err)
			code = -1
			msg = "密钥更新失败了T.T！"
			data = map[string]service.AccessKey{}
			goto Response
		}
		// 处理数据
		if err := service.UpdateKey(para, uid); err != nil {
			fmt.Println("update failed: ", err)
			code = -1
			msg = "密钥更新失败了T.T！"
			data = map[string]service.AccessKey{}
			goto Response
		}
		key, err = service.GetKey_db(uid)
		if err != nil {
			fmt.Println("get failed: ", err)
			code = -1
			msg = "密钥更新失败了T.T！"
			data = map[string]service.AccessKey{}
			goto Response
		}
		code = 1
		msg = "密钥更新成功~"
		data = map[string]service.AccessKey{
			"accesskey": key,
		}
	}

Response:
	// 数据封装至json
	resp_str := service.ASkeyResponse{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	resp_json, err := json.Marshal(resp_str)
	if err != nil {
		fmt.Println("marshal failed: ", err)
		return
	}
	// 设置请求头、写入响应
	if isPreflight := cors(&w, req); isPreflight == true {
		return
	}
	w.Write(resp_json)
}

func LoginHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("请选择操作：\n1. 登录\n2. 注册"))
}
func Start() {
	http.HandleFunc("/ServiceManagement/getTemplateList", getTemplateListHandler)
	http.HandleFunc("/ServiceManagement/getKey", getKeyHandler)
	http.HandleFunc("/ServiceManagement/updateKey", updateKeyHandler)
	http.HandleFunc("/ServiceManagement/addCloudInstance", addCloudInstanceHandler)

	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
