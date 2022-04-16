local function termPrint(target, s)
    local old = term.current()
    term.redirect(target)
    print(s)
    term.redirect(old)
end

local MESSAGE_NOT_CHUNKED_BYTE = 0
local MESSAGE_START_CHUNK_BYTE = 1
local MESSAGE_CONTINNUE_CHUNK_BYTE = 2
local MESSAGE_END_CHUNK_BYTE = 3

-- 512 seems to be the bigest working value in ccemux
local MAX_CHUNK_SIZE = 512 - 2 -- 2 bytes for header


local function ws_chunkedSend(baseSend, data, isBinary)
    local bufID = 0

    local chunks = {}
    if #data >= MAX_CHUNK_SIZE then
        for i = 0, math.ceil(#data / MAX_CHUNK_SIZE) - 1 do
            if i == 0 then
                chunks[i + 1] = string.char(MESSAGE_START_CHUNK_BYTE, bufID)
            elseif i == math.ceil(#data / MAX_CHUNK_SIZE) - 1 then
                chunks[i + 1] = string.char(MESSAGE_END_CHUNK_BYTE, bufID)
            else
                chunks[i + 1] = string.char(MESSAGE_CONTINNUE_CHUNK_BYTE, bufID)
            end

            chunks[i + 1] = chunks[i + 1] .. string.sub(data, i * MAX_CHUNK_SIZE + 1, (i + 1) * MAX_CHUNK_SIZE)
        end
    else
        chunks[1] = string.char(MESSAGE_NOT_CHUNKED_BYTE) .. data
    end

    for _, c in ipairs(chunks) do
        baseSend(c, isBinary)
    end
end

return {
    termPrint = termPrint,
    ws_chunkedSend = ws_chunkedSend
}