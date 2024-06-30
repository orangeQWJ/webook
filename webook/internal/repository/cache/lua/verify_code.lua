-- phone_code:login:130xxxxxx
local key = KEYS[1]
-- 用户提交的验证码
local expectedCode = ARGV[1]
-- 实际的验证码
local code = redis.call("get", key)

-- 验证码还能验证几次
local cntKey = key .. ":cnt"
local cnt = tonumber(redis.call("get", cntKey))

if cnt <= 0 then
	-- 用户多次输错,不能再尝试了
	-- 验证码已经失效
	return -1
elseif expectedCode == code then
	-- 输入对了, 使此验证码失效
	--redis.call("del", key)
	redis.call("set", cntKey, -1)
	return 0
else
	-- 用户输错了,但是还可以再试几次
	-- 可验证次数减一, 当次数减至0,验证码失效
	redis.call("set", cntKey, cnt - 1)
	return -2
end
