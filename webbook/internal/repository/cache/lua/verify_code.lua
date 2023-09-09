local key = KEYS[1]
-- 用户输入的code
local expectedCode = ARGV[1]
local code = redis.call("get", key)
local cntKey = key .. ":cnt"

local cnt = tonumber(redis.call("get", cntKey))
if cnt <= 0 then
    -- 说明， 用户一直输错
    -- 或者已经用过了
    return -1
elseif expectedCode == code then
    -- 输入对了
    -- 用完了， 不能再使用了
    redis.call("set", cntKey, -1)
    return 0
else
    -- 用户输入错了，改验证码的可用次数减少1
    redis.call("decr", cntKey)
    return -2
end
