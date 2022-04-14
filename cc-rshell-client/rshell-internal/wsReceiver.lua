local URL = "localhost:8080"
local RECONNECT_ATTEMPTS = 10
local RECONNECT_TIMEOUT = 10

function dump(o)
    if type(o) == 'table' then
       local s = '{ '
       for k,v in pairs(o) do
          if type(k) ~= 'number' then k = '"'..k..'"' end
          s = s .. '['..k..'] = ' .. dump(v) .. ','
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

    serverNotification = function (localTerm, msg)
        if msg.message then
           localTerm.print(string.format("[*] Server: %s", msg.message))
        else
            localTerm.print("[!] Received invalid serverNotification message.")
        end
    end
}

local function _connectWebSocket(localTerm)
    for i=0, RECONNECT_ATTEMPTS do
        local ws = http.websocket("ws://"..URL)
        if ws then
            return ws
        end

        localTerm.print(string.format("[!] Failed to connect to %s. Retrying in %d seconds...", URL, RECONNECT_TIMEOUT))
        sleep(RECONNECT_TIMEOUT)
    end

    error(string.format("unable to reach %s after %d trys.", URL, RECONNECT_ATTEMPTS))
end

local function _handleMessage(rawMessage, localTerm)
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

            if not isBinary then
                _handleMessage(msg, localTerm)
            end
        end
    end
end

return {
    WebSocketReceiver=WebSocketReceiver
}