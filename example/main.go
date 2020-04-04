package main

import (
	"log"
	"net/http"
	"net/rpc"
	"time"

	"github.com/weihongguo/gglmm"
	redis "github.com/weihongguo/gglmm-redis"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Test --
type Test struct {
	gglmm.Model
	IntValue    int     `json:"intValue"`
	FloatValue  float64 `json:"floatValue"`
	StringValue string  `json:"stringValue"`
}

// Cache --
func (test Test) Cache() bool {
	return true
}

// RPCTestService --
type RPCTestService struct {
	repository *gglmm.GormRepository
}

// NewRPCTestService --
func NewRPCTestService() *RPCTestService {
	return &RPCTestService{
		repository: gglmm.DefaultGormRepository(),
	}
}

// Actions --
func (service *RPCTestService) Actions(cmd string, actionInfos *[]gglmm.RPCActionInfo) error {
	*actionInfos = append(*actionInfos, gglmm.RPCActionInfo{
		Name:     "Get",
		Request:  "string",
		Response: "*Test",
	})
	*actionInfos = append(*actionInfos, gglmm.RPCActionInfo{
		Name:     "List",
		Request:  "gglmm.FilterRequest",
		Response: "*[]Test",
	})
	return nil
}

// Get --
func (service *RPCTestService) Get(idRequest gglmm.IDRequest, test *Test) error {
	err := service.repository.Get(test, idRequest)
	if err != nil {
		return err
	}
	return nil
}

// List --
func (service *RPCTestService) List(filterRequest gglmm.FilterRequest, tests *[]Test) error {
	service.repository.List(tests, filterRequest)
	return nil
}

func main() {
	gglmm.RegisterGormRepository("mysql", "example:123456@(127.0.0.1:3306)/example?charset=utf8mb4&parseTime=true&loc=UTC", 10, 5, 600)
	defer gglmm.CloseGormRepository()

	redisCacher := redis.NewCacher("tcp", "127.0.0.1:6379", 5, 10, 3, 30)
	defer redisCacher.Close()
	gglmm.RegisterCacher(redisCacher)

	gglmm.RegisterBasePath("/api/example")

	testHTTPService := gglmm.NewHTTPService(Test{})
	testHTTPService.HandleBeforeStoreFunc(beforeStore)
	testHTTPService.HandleBeforeUpdateFunc(beforeUpdate)
	gglmm.RegisterHTTPHandler(testHTTPService, "/test").
		Middleware(gglmm.Middleware{
			Name: "example",
			Func: middlewareFunc, // 登录态中间件请参考gglmm-account
		}).
		RESTAction(gglmm.RESTAll)

	gglmm.RegisterRPCHandler(NewRPCTestService(), "RPCTestService")

	go testRPC()

	gglmm.ListenAndServe(":10000")
}

func middlewareFunc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("middleware before")
		next.ServeHTTP(w, r)
		log.Printf("middleware after")
		return
	})
}

func beforeStore(model interface{}) interface{} {
	log.Printf("%#v\n", model)
	test, ok := model.(*Test)
	if !ok {
		return nil
	}
	test.StringValue = "string"
	return test
}

func beforeUpdate(model interface{}, id int64) (interface{}, int64) {
	test, ok := model.(*Test)
	if !ok {
		return nil, 0
	}
	test.StringValue = "string"
	return test, id
}

func testRPC() {
	time.Sleep(2 * time.Second)
	client, err := rpc.DialHTTP("tcp", ":10000")
	if err != nil {
		log.Println("rpc", err)
		return
	}

	idRequest := gglmm.IDRequest{
		ID: 1,
	}
	test := Test{}
	err = client.Call("RPCTestService.Get", idRequest, &test)
	if err != nil {
		log.Println("RPCTestService.Get", err)
	} else {
		log.Printf("Get: %#v", test)
	}

	filterRequest := gglmm.FilterRequest{}
	tests := make([]Test, 0)
	err = client.Call("RPCTestService.List", filterRequest, &tests)
	if err != nil {
		log.Println("RPCTestService.List", err)
	} else {
		log.Printf("List: %#v", tests)
	}
}
