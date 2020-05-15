package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/rpc"
	"time"

	"github.com/weihongguo/gglmm"
	redis "github.com/weihongguo/gglmm-redis"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Example --
type Example struct {
	gglmm.Model
	IntValue    int     `json:"intValue"`
	FloatValue  float64 `json:"floatValue"`
	StringValue string  `json:"stringValue"`
}

// ResponseKey --
func (example Example) ResponseKey() [2]string {
	return [...]string{"example", "examples"}
}

// Cache --
func (example Example) Cache() bool {
	return true
}

// ExampleAction --
func ExampleAction(w http.ResponseWriter, r *http.Request) {
	gglmm.OkResponse().JSON(w)
}

// RPCExampleService --
type RPCExampleService struct {
	gormDB *gglmm.GormDB
}

// NewRPCExampleService --
func NewRPCExampleService() *RPCExampleService {
	return &RPCExampleService{
		gormDB: gglmm.DefaultGormDB(),
	}
}

// Actions --
func (service *RPCExampleService) Actions(cmd string, actions *[]*gglmm.RPCAction) error {
	*actions = append(*actions, []*gglmm.RPCAction{
		gglmm.NewRPCAction("Get", "string", "*Test"),
		gglmm.NewRPCAction("List", "gglmm.FilterRequest", "*[]Test"),
	}...)
	return nil
}

// Get --
func (service *RPCExampleService) Get(idRequest gglmm.IDRequest, example *Example) error {
	err := service.gormDB.Get(example, idRequest)
	if err != nil {
		return err
	}
	return nil
}

// List --
func (service *RPCExampleService) List(filterRequest gglmm.FilterRequest, examples *[]Example) error {
	service.gormDB.List(examples, filterRequest)
	return nil
}

func main() {
	gglmm.RegisterGormDB("mysql", "example:123456@(127.0.0.1:3306)/example?charset=utf8mb4&parseTime=true&loc=UTC", 10, 5, 600)
	defer gglmm.CloseGormDB()

	redisCacher := redis.NewCacher("tcp", "127.0.0.1:6379", 5, 10, 3, 30)
	defer redisCacher.Close()
	gglmm.RegisterCacher(redisCacher)

	gglmm.BasePath("/api")

	// 认证中间间
	// authenticationMiddlerware := gglmm.authenticationMiddlerware("example")

	exampleMiddleware := gglmm.Middleware{
		Name: "ExampleMiddleware",
		Func: middlewareFunc,
	}

	gglmm.UseTimeLogger(true)

	exampleService := gglmm.NewHTTPService(Example{}).
		HandleBeforeCreateFunc(beforeCreate).
		HandleBeforeUpdateFunc(beforeUpdate)
	gglmm.HandleHTTP("/example", exampleService).
		Action(exampleMiddleware, gglmm.ReadActions, gglmm.WriteActions, gglmm.DeleteActions)

	gglmm.HandleHTTPAction("/example_action", ExampleAction, "GET", "POST").
		Middleware(exampleMiddleware)

	gglmm.RegisterRPC(NewRPCExampleService())

	go testHTTP()
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

func beforeCreate(model interface{}) (interface{}, error) {
	log.Printf("%#v\n", model)
	example, ok := model.(*Example)
	if !ok {
		return nil, gglmm.ErrModelType
	}
	example.StringValue = "string"
	return example, nil
}

func beforeUpdate(model interface{}, id int64) (interface{}, int64, error) {
	example, ok := model.(*Example)
	if !ok {
		return nil, 0, gglmm.ErrModelType
	}
	example.StringValue = "string"
	return example, id, nil
}

func testHTTP() {
	time.Sleep(4 * time.Second)

	response, err := http.Get("http://localhost:10000/api/example/1")
	if err != nil {
		log.Println("http", err)
		return
	}
	defer response.Body.Close()

	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("ReadAll", err)
		return
	}

	fmt.Println()
	log.Println(string(result))
}

func testRPC() {
	time.Sleep(2 * time.Second)

	client, err := rpc.DialHTTP("tcp", ":10000")
	if err != nil {
		log.Println("rpc", err)
		return
	}

	fmt.Println()

	idRequest := gglmm.IDRequest{
		ID: 1,
	}
	example := Example{}
	err = client.Call("RPCExampleService.Get", idRequest, &example)
	if err != nil {
		log.Println("RPCExampleService.Get", err)
	} else {
		log.Printf("Get: \n%+v", example)
	}

	fmt.Println()

	filterRequest := gglmm.FilterRequest{}
	filterRequest.AddFilter("id", gglmm.FilterOperateGreaterEqual, 2)
	filterRequest.AddFilter("id", gglmm.FilterOperateLessThan, 4)
	examples := make([]Example, 0)
	err = client.Call("RPCExampleService.List", filterRequest, &examples)
	if err != nil {
		log.Println("RPCExampleService.List", err)
	} else {
		log.Printf("List: \n%+v", examples)
	}
}
