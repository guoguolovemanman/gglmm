module gglmm-example-server

go 1.13

replace github.com/weihongguo/gglmm => ../../

replace github.com/weihongguo/gglmm-redis => ../../../gglmm-redis

require (
	github.com/gorilla/websocket v1.4.2
	github.com/jinzhu/gorm v1.9.12
	github.com/weihongguo/gglmm v0.0.0-20200225064623-73efc6160d28
	github.com/weihongguo/gglmm-redis v0.0.0-20200517090511-b7b885354c4d
)
