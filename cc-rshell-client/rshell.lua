local rec = require("rshell-internal.wsReceiver")
local utils = require("rshell-internal.utils")
local run = require("rshell-internal.runner")

math.randomseed(os.epoch(), os.clock())

local localTerm = term.current()
localTerm["print"] = function(s)
    utils.termPrint(localTerm, s)
end

run.Runner(localTerm, "echo.lua")
rec.WebSocketReceiver(localTerm)
