local socketService = require("rshell-internal.socketService")
local utils = require("rshell-internal.utils")
local run = require("rshell-internal.runner")

math.randomseed(os.epoch(), os.clock())

local localTerm = term.current()
localTerm["print"] = function(s)
    utils.termPrint(localTerm, s)
end

while true do
    local ws = socketService.NewWebSocket(localTerm)
    socketService.WebSocketMainLoop(localTerm, ws)
end