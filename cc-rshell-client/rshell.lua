local socketService = require("rshell-internal.socketService")
local utils = require("rshell-internal.utils")
local run = require("rshell-internal.runner")

math.randomseed(os.epoch(), os.clock())

local localTerm = term.current()
localTerm["print"] = function(s)
    utils.termPrint(localTerm, s)
end

local procID = 0

while true do
    local ws = socketService.NewWebSocket(localTerm)
    run.Runner(localTerm, procID, "echo.lua")
    socketService.WebSocketMainLoop(localTerm, ws)
end