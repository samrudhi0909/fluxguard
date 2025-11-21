-- redis keys: [1] the user_id
-- argv: [1] capacity, [2] rate, [3] current_timestamp, [4] requested_tokens (usually 1)

local key = KEYS[1]
local capacity = tonumber(ARGV[1])
local rate = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
local requested = tonumber(ARGV[4])

-- 1. Get current state (tokens left, last time we refilled)
local info = redis.call("HMGET", key, "tokens", "last_refilled")
local tokens = tonumber(info[1])
local last_refilled = tonumber(info[2])

-- 2. Initialize if this is a new user
if tokens == nil then
    tokens = capacity
    last_refilled = now
end

-- 3. Calculate Refill (Magic Hand)
-- How many seconds passed since last request?
local delta = math.max(0, now - last_refilled)
-- Add tokens based on time (e.g., 0.5 sec * 10 rate = 5 tokens)
local filled_tokens = math.min(capacity, tokens + (delta * rate))

-- 4. Can we afford this request?
local allowed = 0
if filled_tokens >= requested then
    allowed = 1
    filled_tokens = filled_tokens - requested
end

-- 5. Save state back to Redis
redis.call("HMSET", key, "tokens", filled_tokens, "last_refilled", now)
-- Set expiry (cleanup dead keys after 60 seconds of inactivity)
redis.call("EXPIRE", key, 60)

return allowed