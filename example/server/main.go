package main

import (
	"log"
	"net/http"
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

// ExampleService --
type ExampleService struct {
	*gglmm.HTTPService
}

// NewExampleService --
func NewExampleService() *ExampleService {
	return &ExampleService{
		HTTPService: gglmm.NewHTTPService(Example{}),
	}
}

// ExampleAction --
func (service *ExampleService) ExampleAction(w http.ResponseWriter, r *http.Request) {
	gglmm.OkResponse().JSON(w)
}

// ExampleAction --
func ExampleAction(w http.ResponseWriter, r *http.Request) {
	gglmm.OkResponse().JSON(w)
}

var tick int = 0

// TestWSHandler --
func TestWSHandler(chanRequest <-chan *gglmm.WSMessage) <-chan *gglmm.WSMessage {
	chanResponse := make(chan *gglmm.WSMessage)

	go func() {
		ticker := time.NewTicker(6 * time.Second)
		defer func() {
			close(chanResponse)
			ticker.Stop()
		}()

		for {
			select {
			case message, ok := <-chanRequest:
				if message == nil || !ok {
					return
				}
				// log.Println("server handler receive ", string(message.Content))
				// 查询结果
				tick++
				if tick > 5 {
					log.Println("server handler normal")
					chanResponse <- gglmm.NewWSMessage([]byte("success"), false)
					return
				}
				if message.Close {
					return
				}
			case <-ticker.C:
				// 查询结果
				log.Println("server handler timout")
				chanResponse <- gglmm.NewWSMessage([]byte("fail"), true)
				return
			}
		}
	}()

	return chanResponse
}

// ExampleRPCService --
type ExampleRPCService struct {
	gormDB *gglmm.GormDB
}

// NewExampleRPCService --
func NewExampleRPCService() *ExampleRPCService {
	return &ExampleRPCService{
		gormDB: gglmm.DefaultGormDB(),
	}
}

// Actions --
func (service *ExampleRPCService) Actions(cmd string, actions *[]*gglmm.RPCAction) error {
	*actions = append(*actions, []*gglmm.RPCAction{
		gglmm.NewRPCAction("Get", "string", "*Test"),
		gglmm.NewRPCAction("List", "gglmm.FilterRequest", "*[]Test"),
	}...)
	return nil
}

// Get --
func (service *ExampleRPCService) Get(idRequest gglmm.IDRequest, example *Example) error {
	err := service.gormDB.Get(example, idRequest)
	if err != nil {
		return err
	}
	return nil
}

// List --
func (service *ExampleRPCService) List(filterRequest gglmm.FilterRequest, examples *[]Example) error {
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

	gglmm.UseTimeLogger(true)

	exampleService := NewExampleService()
	exampleService.
		HandleBeforeCreateFunc(beforeCreate).
		HandleBeforeUpdateFunc(beforeUpdate)

	readMiddleware := gglmm.Middleware{
		Name: "ReadMiddleware",
		Func: middlewareFunc,
	}
	gglmm.HandleHTTP("/example", exampleService).
		Action(readMiddleware, gglmm.ReadActions)

	writeMiddleware := gglmm.Middleware{
		Name: "WriteMiddleware",
		Func: middlewareFunc,
	}
	gglmm.HandleHTTP("/example", exampleService).
		Action(writeMiddleware, gglmm.WriteActions)

	deleteMiddleware := gglmm.Middleware{
		Name: "DeleteMiddleware",
		Func: middlewareFunc,
	}
	gglmm.HandleHTTP("/example", exampleService).
		Action(deleteMiddleware, gglmm.DeleteActions)

	exampleMiddleware := gglmm.Middleware{
		Name: "ExampleMiddleware",
		Func: middlewareFunc,
	}
	gglmm.HandleHTTPAction("/example_action", exampleService.ExampleAction, "GET").
		Middleware(exampleMiddleware)
	gglmm.HandleHTTPAction("/example_action", ExampleAction, "POST").
		Middleware(exampleMiddleware)

	gglmm.HandleWS("/ws/example", TestWSHandler)

	gglmm.RegisterRPC(NewExampleRPCService())

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
