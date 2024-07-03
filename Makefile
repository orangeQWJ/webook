# 用来压力测试
.PHONY: wrk-login mock wrk-signup
wrk-signup:
	wrk -t5 -d10s -c30 -s ./scripts/wrk/signup.lua http://localhost:8080/users/signup
#-t[x]：使用x个线程。
#-d[x]s：测试持续时间为x秒。
#-c[x]：使用x个并发连接。
wrk-login:
	wrk -t5 -d10s -c30 -s ./scripts/wrk/login.lua http://localhost:8080/users/login

wrk-profile:
	wrk -t4 -d10s -c4 -s ./scripts/wrk/profile.lua http://localhost:8080/users/profile

mock:
	@mockgen -source=/Users/orange/code/go/src/webook/webook/internal/service/code.go -package=svcmocks -destination=/Users/orange/code/go/src/webook/webook/internal/service/mocks/code.mock.go
	@mockgen -source=/Users/orange/code/go/src/webook/webook/internal/service/user.go -package=svcmocks -destination=/Users/orange/code/go/src/webook/webook/internal/service/mocks/user.mock.go
	@mockgen -source=/Users/orange/code/go/src/webook/webook/internal/repository/user.go -package=repomocks -destination=/Users/orange/code/go/src/webook/webook/internal/repository/mocks/user.mock.go
	@mockgen -source=/Users/orange/code/go/src/webook/webook/internal/repository/code.go -package=repomocks -destination=/Users/orange/code/go/src/webook/webook/internal/repository/mocks/code.mock.go
	@go mod tidy
