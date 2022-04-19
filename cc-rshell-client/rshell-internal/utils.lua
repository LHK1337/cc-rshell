---{ Chunk header }---
-- Size: 1 Byte
---- 2 bits: OP_CODE
---- 6 bits: BUFFER_ID
local MAX_BUFFER_ID = 0x3F

local CHUNK_OPCODE_NOT_CHUNKED = 0
local CHUNK_OPCODE_START_CHUNK = 1
local CHUNK_OPCODE_CONTINUE_CHUNK = 2
local CHUNK_OPCODE_END_CHUNK = 3

-- 512 seems to be the bigest working value in ccemux
local MAX_CHUNK_SIZE = 512 - 1 -- 1 bytes for header

local CurrentBufferID = 0

local function ws_chunkedSend(baseSend, data, isBinary)
    local bufID = CurrentBufferID
    CurrentBufferID = (CurrentBufferID + 1) % (MAX_BUFFER_ID + 1)

    local chunks = {}
    if #data >= MAX_CHUNK_SIZE then
        for i = 0, math.ceil(#data / MAX_CHUNK_SIZE) - 1 do
            local opcode
            if i == 0 then
                opcode = CHUNK_OPCODE_START_CHUNK
            elseif i == math.ceil(#data / MAX_CHUNK_SIZE) - 1 then
                opcode = CHUNK_OPCODE_END_CHUNK
            else
                opcode = CHUNK_OPCODE_CONTINUE_CHUNK
            end

            local chunkHeader = bit.bor(bit.blshift(opcode, 6), bufID)
            chunks[i + 1] = string.char(chunkHeader) .. string.sub(data, i * MAX_CHUNK_SIZE + 1, (i + 1) * MAX_CHUNK_SIZE)
        end
    else
        chunks[1] = string.char(CHUNK_OPCODE_NOT_CHUNKED) .. data
    end

    for _, c in ipairs(chunks) do
        baseSend(c, isBinary)
    end
end

local function yield()
    os.queueEvent("yield")
    os.pullEvent("yield")
end

local function termPrint(target, s)
    local old = term.current()
    term.redirect(target)
    print(s)
    term.redirect(old)
end

local function shallowcopy(orig)
    local orig_type = type(orig)
    local copy
    if orig_type == 'table' then
        copy = {}
        for orig_key, orig_value in pairs(orig) do
            copy[orig_key] = orig_value
        end
    else
        -- number, string, boolean, etc
        copy = orig
    end
    return copy
end

local function deepcopy(orig)
    yield()
    local orig_type = type(orig)
    local copy
    if orig_type == 'table' then
        copy = {}
        for orig_key, orig_value in next, orig, nil do
            copy[deepcopy(orig_key)] = deepcopy(orig_value)
        end
        setmetatable(copy, deepcopy(getmetatable(orig)))
    else
        -- number, string, boolean, etc
        copy = orig
    end
    return copy
end

local function dump(o)
    if type(o) == 'table' then
        local s = '{ '
        for k, v in pairs(o) do
            if type(k) ~= 'number' then
                k = '"' .. k .. '"'
            end
            s = s .. '[' .. k .. '] = ' .. dump(v) .. ','
        end
        return s .. '} '
    else
        return tostring(o)
    end
end

return {
    termPrint = termPrint,
    ws_chunkedSend = ws_chunkedSend,

    shallowcopy = shallowcopy,
    deepcopy = deepcopy,
    dump = dump,
}