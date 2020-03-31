module gglmm-example

go 1.13

replace github.com/weihongguo/gglmm => ../

replace github.com/weihongguo/gglmm-redis => ../../gglmm-redis

require (
	github.com/jinzhu/gorm v1.9.12
	github.com/weihongguo/gglmm v0.0.0-20200323132608-c2c78309e5c8
	github.com/weihongguo/gglmm-redis v0.0.0-00010101000000-000000000000
)
