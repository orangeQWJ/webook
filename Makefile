# 用来压力测试
.PHONY: wrk-login
wrk-login:
	wrk -t5 -d10s -c30 -s ./scripts/wrk/signup.lua http://localhost:8080/users/signup
#-t[x]：使用x个线程。
#-d[x]s：测试持续时间为x秒。
#-c[x]：使用x个并发连接。
wrk-login:
	wrk -t5 -d10s -c30 -s ./scripts/wrk/login.lua http://localhost:8080/users/login

wrk-profile:
	wrk -t4 -d10s -c4 -s ./scripts/wrk/profile.lua http://localhost:8080/users/profile
