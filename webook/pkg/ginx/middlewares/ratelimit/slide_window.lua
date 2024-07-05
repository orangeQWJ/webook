-- 1, 2, 3, 4, 5, 6, 7 这是你的元素
-- ZREMRANGEBYSCORE key1 0 6
-- 7 执行完之后

-- 限流对象,限制这个ip的访问次数
local key = KEYS[1]
-- 窗口大小 60 * 1000 毫秒
local window = tonumber(ARGV[1])
-- 阈值 100
local threshold = tonumber(ARGV[2])
--Unix 纪元（1970 年 1 月 1 日 00:00:00 UTC）以来经过的毫秒数。
local now = tonumber(ARGV[3])
-- 窗口的起始时间 一分钟前的时间节点
local min = now - window

--这行代码从 Redis 有序集合中移除时间戳小于 min 的所有记录。ZREMRANGEBYSCORE
--命令用于移除有序集合中指定分数范围内的成员。在这里，它移除所有时间戳小于 min
--的成员，即窗口外的旧记录
redis.call('ZREMRANGEBYSCORE', key, '-inf', min)

--ZREMRANGEBYSCORE 是 Redis 提供的一个命令，其全称是 "Remove Range By Score"（按分数移除范围）。
--这个命令用于从 Redis 的有序集合（Sorted Set）中删除分数在指定范围内的所有成员。
--有序集合（Sorted Set）是一种数据结构，其中每个成员都关联一个分数，成员按分数排序。
--`key` 是有序集合的键名。在这段脚本中，它通常代表限流对象，比如一个特定的 IP 地址。
--'-inf' 表示负无穷大。在分数范围中，它代表从最小的可能值开始，即所有小于 `min` 的分数都会被包含在范围内。
--`min` 是一个变量，表示时间窗口的起始时间。在这里，它是通过当前时间减去窗口大小计算出来的。它的单位是毫秒。
local cnt = redis.call('ZCOUNT', key, '-inf', '+inf')
-- 这行代码统计当前窗口内的请求数。ZCOUNT
-- 命令用于计算有序集合中指定分数范围内的成员数量。
-- 它计算集合中所有成员的数量，因为我们已经移除了窗口外的旧记录
-- local cnt = redis.call('ZCOUNT', key, min, '+inf')
if cnt >= threshold then
	-- 执行限流
	return "true"
else
	-- 把 score 和 member 都设置成 now 否则，将当前时间 now
	-- 作为分数和成员添加到有序集合中，记录这次请求。ZADD
	-- 命令用于向有序集合中添加成员，并设置分数。在这里，它将 now
	-- 作为分数和成员添加到集合中。
	-- 如果 key 不存在，ZADD 命令会自动创建这个有序集合。
	-- 接着，设置有序集合的过期时间为窗口大小，以确保集合不会无限增长。PEXPIRE
	-- 命令用于设置键的过期时间，以毫秒为单位。
	-- 当指定的过期时间到达时，Redis
	-- 会自动删除这个键及其对应的值。这意味着，如果你的有序集合 key
	-- 设置了过期时间，当过期时间到达时，整个有序集合会被自动删除，从而释放内存。

	redis.call('ZADD', key, now, now)
	redis.call('PEXPIRE', key, window)
	return "false"
end
