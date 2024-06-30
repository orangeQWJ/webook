-- 你的验证码在 Redis上的key
-- phone_code:login:130xxxxxx
local key = KEYS[1]

-- 验证次数,我们一个验证码,最多重复三次. 这个记录还可以验证几次
local cntkey = key .. ":cnt"
-- phone_code:login:130xxxxxx:cnt
-- 你的验证码 12356
local val = ARGV[1]

local ttl = tonumber(redis.call("ttl", key))
-- ttl == -1 未设置过期时间
-- ttl == -2 不存在这个键
-- ttl: 有效期剩余秒数

if ttl == -1 then
	-- 系统错误,有人手动设置了这个key,但是没有设置过期时间
	return -2
elseif ttl == -2 or ttl < 120 then
	-- key 不存在(有效期过了,自动清除)或者有效期剩余不足9分钟
	-- 重新设置验证码
	redis.call("set", key, val)
	redis.call("expire", key, 180)
	redis.call("set", cntkey, 3)
	redis.call("expire", cntkey, 180)
	return 0
else
	-- 发送太频繁
	return -1
end
