module gglmm-example-client-http

go 1.13

replace gglmm-example => ../../

replace github.com/weihongguo/gglmm => ../../../../gglmm

replace github.com/weihongguo/gglmm-auth => ../../../../gglmm-auth

require (
	gglmm-example v0.0.0-00010101000000-000000000000
	github.com/weihongguo/gglmm-auth v0.0.0-20200527134404-e0cbabd366f6
)
