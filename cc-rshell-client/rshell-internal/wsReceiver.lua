local URL = "localhost:8080/clients/socket"
local RECONNECT_ATTEMPTS = 1
local RECONNECT_TIMEOUT = 3

local mp = require("rshell-internal.MessagePack")
local msgFactory = require("rshell-internal.messages")
local utils = require("rshell-internal.utils")

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

local function _activateConnection(ws, localTerm)
    local activateMessage = msgFactory.BuildActivateMessage(localTerm)
    local rawMP = mp.pack(activateMessage)
    ws.send(rawMP, true)
end

local function _connectWebSocket(localTerm)
    for _ = 0, RECONNECT_ATTEMPTS do
        local ws = http.websocket("ws://" .. URL)
        if ws then
            -- wrap web socket send method to support message chunking
            local baseSend = ws.send
            ws.send = function(data, isBinary)
                utils.ws_chunkedSend(baseSend, data, isBinary)
            end

            _activateConnection(ws, localTerm)

            localTerm.print("[*] Connected and activated.")
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
