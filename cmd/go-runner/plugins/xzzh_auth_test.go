package plugins

import (
	inHTTP "github.com/Xiangzhuzhihui/apisix-go-plugin-runner/internal/http"
	"github.com/Xiangzhuzhihui/apisix-go-plugin-runner/internal/util"
	hrc "github.com/api7/ext-plugin-proto/go/A6/HTTPReqCall"
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"log"
	"testing"
)

func TestXzzhAuth_ParseConf(t *testing.T) {
	conf := "{\"url\":\"login.smartcloud-service-istio.svc.cluster.local:10801\",\"login_urls\":[\"/api/login/pwd/login\",\"/api/dictionary/kv/hotel/key/system_domain/v1\"],\"external_links\":[\"/api/es/\",\"/api/yw/\"],\"ignore_links\":[\"/api/es/\",\"/api/yw/\"]}"
	xzzhAuth := XzzhAuth{}
	xzzhAuthConf, _ := xzzhAuth.ParseConf([]byte(conf))
	log.Printf("测试结果： %s", xzzhAuthConf)
}

func TestGrpc(t *testing.T) {
	conn, err := grpc.Dial("login.smartcloud-service-istio.svc.cluster.local:10801", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("did not connect: %v", err)
	}
	log.Printf("connect: %v", conn)
}

func TestXzzhAuth_Filter(t *testing.T) {
	token := "Bearer:eyJhdWRJZCI6IjQ0NTI0MjIiLCJleHAiOiIxNzQwNzI3OTcxIiwiaWF0IjoiMTc0MDY0MTU3MSIsImp0aSI6Ii00OTkyMTgxNzMifQ"
	conf := "{\"url\":\"127.0.0.1:10801\",\"login_urls\":[\"/api/login/pwd/login\",\"/api/dictionary/kv/hotel/key/system_domain/v1\"],\"external_links\":[\"/api/es/\",\"/api/yw/\"],\"ignore_links\":[\"/api/es/\",\"/api/yw/\"],\"request_headers\":[\"client\"]}"
	xzzhAuth := XzzhAuth{}
	xzzhAuthConf, _ := xzzhAuth.ParseConf([]byte(conf))

	builder := flatbuffers.NewBuilder(1024)

	hrc.ReqStart(builder)
	hrc.ReqAddId(builder, 233)
	hrc.ReqAddConfToken(builder, 1)
	r := hrc.ReqEnd(builder)
	builder.Finish(r)
	out := builder.FinishedBytes()

	req := inHTTP.CreateRequest(out)
	req.Header().Set("Authorization", token)
	//req.Header().Set("client", "hk_manager_xcx")
	req.SetPath([]byte("/api/customer/merchant/current-user/simple/list/v1"))
	resp := inHTTP.CreateResponse(out)

	//直接使用myHandler，传入参数rr,req
	xzzhAuth.RequestFilter(xzzhAuthConf, nil, req)

	bu := util.GetBuilder()
	assert.True(t, resp.FetchChanges(bu))
	log.Printf(string(bu.Bytes))
}
