local WS_DISPATCH_MESSAGE = "rshell_ws_dispatch"

local function DispatchWebsocketMessage(msg, isBinary, src)
    os.queueEvent(WS_DISPATCH_MESSAGE, msg, isBinary, src)
end

return {
    DispatchWebsocketMessage = DispatchWebsocketMessage,
    WS_DISPATCH_MESSAGE = WS_DISPATCH_MESSAGE,
}
