module gglmm-example

go 1.13

replace github.com/weihongguo/gglmm => ../

replace github.com/weihongguo/gglmm-redis => ../../gglmm-redis

require (
	github.com/gorilla/mux v1.7.4 // indirect
	github.com/jinzhu/gorm v1.9.12
	github.com/weihongguo/gglmm v0.0.0-00010101000000-000000000000
	github.com/weihongguo/gglmm-redis v0.0.0-20200517090511-b7b885354c4d
)
