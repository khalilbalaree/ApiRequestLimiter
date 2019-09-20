package limiter

var luaText = `
local key = KEYS[1] 
local currNanoSec = tonumber(ARGV[1])
local numRequest = tonumber(ARGV[2])
local maxPermits = tonumber(ARGV[3])
local rate = tonumber(ARGV[4])

local lastNanoSec = tonumber(redis.call("HGET", key, "lastNanoSec"))
local currPermits = tonumber(redis.call("HGET", key, "currPermits"))

if (lastNanoSec == nil) then
    redis.call("HSET", key, "lastNanoSec", currNanoSec)
    redis.call("HSET", key, "currPermits", maxPermits-numRequest)
    return true
end

local reservePermits = math.floor((currNanoSec-lastNanoSec)/math.pow(10,9)*rate)
local current = math.min(reservePermits+currPermits, maxPermits)

redis.call("HSET", key, "lastNanoSec", currNanoSec)

local remaining = current - numRequest
if (remaining >= 0) then
    redis.call("HSET", key, "currPermits", remaining)
    return true
end

return false
`
