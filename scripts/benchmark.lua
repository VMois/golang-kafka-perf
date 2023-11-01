math.randomseed(os.time())

local function random_type()
    local types = {"red", "green", "yellow"}
    return types[math.random(#types)]
end

local function random_client_id(length)
    local chars = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+-={}|[]`~'
    local randomString = ''

    charTable = {}
    for c in chars:gmatch"." do
        table.insert(charTable, c)
    end

    for i = 1, length do
        randomString = randomString .. charTable[math.random(1, #charTable)]
    end

    return randomString
end

threads = {} 

setup = function (thread)
    table.insert(threads,thread)
    thread:set("red_counter", 0)
    thread:set("green_counter", 0)
    thread:set("yellow_counter", 0)
end

request = function()
    local payload_type = random_type()
    local payload = '{"type":"' .. payload_type .. '", "client_id": "' .. random_client_id(10) .. '"}'
    return wrk.format("POST", "/metrics", headers, payload)
end

function response(status, headers, body)
    if status == 200 then
        if string.find(body, "red") then
            red_counter = red_counter + 1
        elseif string.find(body, "green") then
            green_counter = green_counter + 1
        elseif string.find(body, "yellow") then
            yellow_counter = yellow_counter + 1
        end
    end
end

function done(summary, latency, requests)
    total = {}
    total["red_counter"] = 0
    total["green_counter"] = 0
    total["yellow_counter"] = 0

    for i, thread in pairs(threads) do
        local red_counter = thread:get("red_counter")
        local green_counter = thread:get("green_counter")
        local yellow_counter = thread:get("yellow_counter")
        total["red_counter"] = total["red_counter"] + red_counter
        total["green_counter"] = total["green_counter"] + green_counter
        total["yellow_counter"] = total["yellow_counter"] + yellow_counter
    end

    print("Red Count: " .. total["red_counter"])
    print("Green Count: " .. total["green_counter"])
    print("Yellow Count: " .. total["yellow_counter"])
end
