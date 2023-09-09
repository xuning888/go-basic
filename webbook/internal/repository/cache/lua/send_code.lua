-- 验证码在redis中的key phone_code:biz:phone
local key = KEYS[1]
-- 验证次数, 我们一个验证码最多重复三次， 当前这个验证码记录还可以用来验证几次
local cntKey = key .. ":cnt"
-- 验证码
local code = ARGV[1]
-- 过期时间
local ttl = tonumber(redis.call("ttl", key))
if ttl == -1 then
    -- key 存在， 但是没有过期时间
    -- 系统错误, 设置了这个key， 但是没有给这个key设置超时时间
    return -2
elseif ttl == -2 or ttl < 540 then
    -- 540 = 600 - 60 一个验证码的过期时间是10分钟，一分钟只能发送一次验证码
    -- 验证码没有过期，并且验证码的过期时间还有9分钟的就不允许再次发送
    redis.call("set", key, code)
    redis.call("expire", key, 600)
    redis.call("set", cntKey, 3)
    redis.call("expire", cntKey, 600)
    return 0
else
    -- 频繁发送
    return -1
end
