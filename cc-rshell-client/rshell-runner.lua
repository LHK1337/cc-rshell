local utils = require("rshell-internal.utils")
local framebuffer = require("rshell-internal.framebuffer")
local mp = require("rshell-internal.MessagePack")
local socketSend = require("rshell-internal.socketSend")
local messages = require("rshell-internal.messages")

local argV = { ... }

local id = argV[1]

local localTerm = term.current()
localTerm["print"] = function(s)
    utils.termPrint(localTerm, s)
end

local w, h = term.getSize()
local fb = framebuffer.New(w, h, true, 0, 0, function(buffer, src)
    socketSend.DispatchWebsocketMessage(mp.pack(messages.BuildBufferUpdateMessage(id, buffer)), true, src)
end)
term.redirect(fb)
shell.run(table.unpack(argV, 2))
os.queueEvent("PROC_TABLE", "close", id)
socketSend.DispatchWebsocketMessage(mp.pack(messages.BuildBufferClosedMessage(id)), true)
