local URL = "localhost:8080/clients/socket"
local RECONNECT_ATTEMPTS = 1
local RECONNECT_TIMEOUT = 3

local mp = require("rshell-internal.MessagePack")
local msgFactory = require("rshell-internal.messages")

local MESSAGE_NOT_CHUNKED_BYTE = 0
local MESSAGE_START_CHUNK_BYTE = 1
local MESSAGE_CONTINNUE_CHUNK_BYTE = 2
local MESSAGE_END_CHUNK_BYTE = 3

-- 512 seems to be the bigest working value in ccemux
local MAX_CHUNK_SIZE = 512 - 2 -- 2 bytes for header

function dump(o)
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

local _msgTypeHandler = {
    event = function (localTerm, msg)
        if msg.event and msg.params then
            localTerm.print(string.format("[*] Event: %s, [%s]", msg.event, dump(msg.params)))
            os.queueEvent(msg.event, table.unpack(msg.params))
        else
            localTerm.print("[!] Received invalid event message.")
        end
    end,

    serverNotification = function(localTerm, msg)
        if msg.message then
            localTerm.print(string.format("[*] Server: %s", msg.message))
        else
            localTerm.print("[!] Received invalid serverNotification message.")
        end
    end
}

local function _activateConnection(ws)
    local activateMessage = msgFactory.BuildActivateMessage()
    local rawMP = mp.pack(activateMessage)

    local bufID = 0

    local chunks = {}
    if #rawMP >= MAX_CHUNK_SIZE then
        for i = 0, math.ceil(#rawMP / MAX_CHUNK_SIZE) - 1 do
            if i == 0 then
                chunks[i + 1] = string.char(MESSAGE_START_CHUNK_BYTE, bufID)
            elseif i == math.ceil(#rawMP / MAX_CHUNK_SIZE) - 1 then
                chunks[i + 1] = string.char(MESSAGE_END_CHUNK_BYTE, bufID)
            else
                chunks[i + 1] = string.char(MESSAGE_CONTINNUE_CHUNK_BYTE, bufID)
            end

            chunks[i + 1] = chunks[i + 1] .. string.sub(rawMP, i * MAX_CHUNK_SIZE + 1, (i + 1) * MAX_CHUNK_SIZE)
        end
    else
        chunks[1] = string.char(MESSAGE_NOT_CHUNKED_BYTE) .. rawMP
    end

    for _, c in ipairs(chunks) do
        ws.send(c, true)
    end
end

local function _connectWebSocket(localTerm)
    for _ = 0, RECONNECT_ATTEMPTS do
        local ws = http.websocket("ws://" .. URL)
        if ws then
            _activateConnection(ws)
            return ws
        end

        localTerm.print(string.format("[!] Failed to connect to %s. Retrying in %d seconds...", URL, RECONNECT_TIMEOUT))
        sleep(RECONNECT_TIMEOUT)
    end

    error(string.format("unable to reach %s after %d trys.", URL, RECONNECT_ATTEMPTS))
end

local function _handleMessageJSON(rawMessage, localTerm)
    local msg, err = textutils.unserialiseJSON(rawMessage)
    if msg == nil then
        localTerm.print(string.format("[!] Received invalid JSON. Error: %s", err))
    end

    if msg.type then
        if _msgTypeHandler[msg.type] == nil then
            localTerm.print(string.format("[!] Received unsupported message type (%s).", msg.type))
        else
            _msgTypeHandler[msg.type](localTerm, msg)
        end
    end
end

local function _handleMessageMessagePack(rawMessage, localTerm)
    local success, msg = pcall(mp.unpack, rawMessage)
    if not success then
        localTerm.print(string.format("[!] Received invalid MessagePack. Error: %s", msg))
        return
    end

    if msg.type then
        if _msgTypeHandler[msg.type] == nil then
            localTerm.print(string.format("[!] Received unsupported message type (%s).", msg.type))
        else
            _msgTypeHandler[msg.type](localTerm, msg)
        end
    end
end

function WebSocketReceiver(localTerm)
    if not http.checkURL("http://"..URL) then
        error(string.format("not allowed to connect %s.", URL))
    end

    while true do
        local ws = _connectWebSocket(localTerm)

        while true do
            local msg, isBinary = ws.receive()

            if msg == nil then
                localTerm.print("[!] Lost connection. Reconnecting...")
                break
            end

            if isBinary then
                _handleMessageMessagePack(msg, localTerm)
            else
                _handleMessageJSON(msg, localTerm)
            end
        end
    end
end

return {
    WebSocketReceiver=WebSocketReceiver
}
