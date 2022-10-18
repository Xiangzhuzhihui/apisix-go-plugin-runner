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
	token := "Bearer:eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiItMSIsImF1ZCI6Imxpc2hhbl9hZG1pbiIsImF1ZElkIjoiNDQ1MjQyMiIsImlzcyI6Inh6emguY29tIiwic3ViVHlwZSI6IjEiLCJleHAiOiIyMDIxLTA5LTAxIDA5OjA0OjM4IiwiaWF0IjoiMjAyMS0wOC0zMSAwOTowNDozOCIsImp0aSI6IjI3MjgxNDQ0MTAzNiJ9.Btou3uiXkEj1QxBr9MPumi35DwrnCkhDAQyPHiTTOageuE_no7fKiP1cBkjzsG-3ceYrBM_lL_rz_TwioW3JfK9yBIzXhzFRhx_eMzOGoBSNjTwtHE_GXnvkexxbcvgNzXmY4zPAIP9j37SkJYrHtgmpWnB9iGXPfgptH48VZnLL-g8wEFBoWlMgivsj9Zg3oK4-WLsjR36MaDxwQDrAD_Mhe1QLu8Vlbet6TQuPf08B3g47a8X9XQV9hUCvrj5QRLFscXpcw0DM3m50d3gE3EK3GJdm_1ULQwsBjfzA2iKNscg6zdNYmi_VgKzUDI6WmnsRUyshV9VADZUiLgrxLTdhHzayv-1XKVlk5JhpXaLZEkejMj_-S9__82laK_KKvuaebSuRRwAqjtmuCWJb9FVzK-C_zR0J6EyZH8DwNLQRJT5EdJ4akOGySzueFgay4ZEm-wH2Uacwop_O7JD6HviKaK_wJK2ajisZOtI_qypKQy-wmGvrQsnYHjo5JztFZSV80KKzS9_B0nv6kg9Q0B45WU_-4iodZZ93cYjB5wGvkW05qgnzcySXWVSsgpjBvOnyeKyysSSL5-z6WPxjGpBk924xBYvwOwavjTiyxtd6LVuoEFFTwNImF-AhCqyjIbElrRCQY5xAQb6j-AT8rpi-q8WdKXDCdpLdKIAfMMM"
	conf := "{\"url\":\"127.0.0.1:10801\",\"login_urls\":[\"/api/login/pwd/login\",\"/api/dictionary/kv/hotel/key/system_domain/v1\"],\"external_links\":[\"/api/es/\",\"/api/yw/\"],\"ignore_links\":[\"/api/es/\",\"/api/yw/\"]}"
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
	req.SetPath([]byte("/api/login/user/router/1/v2"))
	resp := inHTTP.CreateResponse(out)

	//直接使用myHandler，传入参数rr,req
	xzzhAuth.RequestFilter(xzzhAuthConf, resp, req)

	bu := util.GetBuilder()
	assert.True(t, resp.FetchChanges(bu))
	log.Printf(string(bu.Bytes))
}
