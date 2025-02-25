/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package plugins

import (
	"encoding/json"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strings"

	"context"
	pb "github.com/Xiangzhuzhihui/apisix-go-plugin-runner/gen/api/v1/loginService"
	pkgHTTP "github.com/Xiangzhuzhihui/apisix-go-plugin-runner/pkg/http"
	"github.com/Xiangzhuzhihui/apisix-go-plugin-runner/pkg/log"
	"github.com/Xiangzhuzhihui/apisix-go-plugin-runner/pkg/plugin"
	"google.golang.org/grpc"
	"time"
)

func init() {
	err := plugin.RegisterPlugin(&XzzhAuth{})
	if err != nil {
		log.Fatalf("failed to register plugin XzzhAuth: %s", err)
	}

	log.Infof("success to register plugin XzzhAuth")
}

// XzzhAuth 自定义认证鉴权插件
type XzzhAuth struct {
}

type XzzhResult struct {
	StatusCode int32  `json:"statusCode"`
	Msg        string `json:"msg"`
	ErrorCode  string `json:"errorCode"`
}

type XzzhAuthConf struct {
	Body string `json:"body"`
	Url  string `json:"url"`
	// 登录路径等不需要认证和鉴权的路径
	LoginUrls []string `json:"login_urls"`
	// 外部路径，即不要在路径中添加用户信息参数
	ExternalLinks []string `json:"external_links"`
	// 不需要鉴权和认证的路径
	IgnoreLinks []string `json:"ignore_links"`
}

func (p *XzzhAuth) Name() string {
	return "xzzhAuth"
}

func (p *XzzhAuth) ParseConf(in []byte) (interface{}, error) {
	conf := XzzhAuthConf{}
	err := json.Unmarshal(in, &conf)
	return conf, err
}

func (p *XzzhAuth) RequestFilter(conf interface{}, w http.ResponseWriter, r pkgHTTP.Request) {
	xzzhAuthConf := conf.(XzzhAuthConf)
	if len(xzzhAuthConf.Url) == 0 {
		return
	}
	method := strings.ToUpper(r.Method())
	path := strings.Split(string(r.Path()), "?")[0]
	for _, v := range xzzhAuthConf.LoginUrls {
		if v == path {
			return
		}
	}
	for _, v := range xzzhAuthConf.IgnoreLinks {
		if strings.HasPrefix(path, v) {
			return
		}
	}
	for _, v := range xzzhAuthConf.ExternalLinks {
		if strings.HasPrefix(path, v) {
			return
		}
	}
	token := r.Header().Get("Authorization")
	result := postLoginService(token, method, path, &xzzhAuthConf)
	if result.StatusCode == 2000 {
		args := r.Args()
		addParam := result.ApisixUserInfo.AddPathParam
		if addParam != nil {
			for k, v := range addParam {
				args.Set(k, v)
			}
		}
		header := r.Header()
		addHeader := result.ApisixUserInfo.AddHeaderParam
		if addHeader != nil {
			for k, v := range addHeader {
				header.Set(k, v)
			}
		}
		return
	}
	w.WriteHeader(200)
	re, jsonErr := json.Marshal(XzzhResult{StatusCode: result.GetStatusCode(), Msg: result.Msg, ErrorCode: result.GetErrorCode()})
	if jsonErr != nil {
		log.Errorf("Json Marshal Error2: %s", jsonErr)
	}
	_, err := w.Write(re)
	if err != nil {
		log.Errorf("failed to write: %s", err)
	}
}

func (p *XzzhAuth) ResponseFilter(conf interface{}, w pkgHTTP.Response) {

}

// Validate method: 请求方式
// path: 请求路径
func postLoginService(token string, method string, path string, conf *XzzhAuthConf) *pb.ResultApiSix {
	// Set up a connection to the server.
	conn, err := grpc.Dial(conf.Url, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Errorf("did not connect: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Errorf("failed to close: %s", err)
		}
	}(conn)
	c := pb.NewApiSixServerClient(conn)

	// Contact the server and print out its response.
	//ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	//defer cancel()
	clientDeadline := time.Now().Add(3 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), clientDeadline)
	defer cancel()
	result, err := c.DecodeAndVerifyV3(ctx, &pb.Jwt{Jwt: token, Path: method + ":" + path})
	if err != nil {
		//获取错误状态
		stats, ok := status.FromError(err)
		if ok {
			//判断是否为调用超时
			if stats.Code() == codes.DeadlineExceeded {
				log.Errorf("DecodeAndVerifyV3 timeout!")
			}
		}
		log.Errorf("could not greet: %v", err)
	}
	return result
}
