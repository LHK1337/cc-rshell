local mp = require "rshell-internal.MessagePack"
local messages = require "rshell-internal.messages"

local DEFAULT_RPC_TIMEOUT = 10
local RSHELL_RPC_RESULT_EVENT = "rshell_rpc_result"

local function RShellTermRedirect(ws, localTerm)
    local redirectID = math.random(0, 0xffffff)
    local r = {}

    local function sendRPC(name, params)
        local rawMP = mp.pack(messages.BuildRPCMessage(redirectID, "term." .. name, params))
        ws.send(rawMP, true)
    end

    local function awaitRPCResult(name, timeout)
        if timeout == nil then
            timeout = DEFAULT_RPC_TIMEOUT
        end

        local timeoutTimerID = os.startTimer(timeout)
        while true do
            local eventData = { os.pullEvent() }
            local event = eventData[1]
            local id = eventData[2]

            if event == "timer" and id == timeoutTimerID then
                localTerm.print(string.format("[!] RPC %s timed out for %d.", name, redirectID))
                return nil
            elseif event == RSHELL_RPC_RESULT_EVENT and id == redirectID and eventData[3] == name then
                return eventData[4]
            else
                os.queueEvent(table.unpack(eventData))
            end
        end
    end

    r["write"] = function(text)
        sendRPC("write", { text = text })
    end
    r["scroll"] = function(y)
        sendRPC("scroll", { y = y })
    end
    r["getCursorPos"] = function()
        sendRPC("getCursorPos", nil)
        local resultV = awaitRPCResult("getCursorPos")
        if resultV == nil then
            return 0, 0
        end
        return resultV["x"], resultV["y"]
    end
    r["setCursorPos"] = function(x, y)
        sendRPC("getCursorPos", {
            x = x, y = y
        })
    end
    r["getCursorBlink"] = function()
        sendRPC("getCursorBlink", nil)
        return awaitRPCResult("getCursorBlink")[1] == true
    end
    r["setCursorBlink"] = function(blink)
        sendRPC("setCursorBlink", {
            blink = blink
        })
    end
    r["getSize"] = function()
        sendRPC("getSize", nil)
        local resultV = awaitRPCResult("getSize")
        if resultV == nil then
            return 0, 0
        end
        return resultV["width"], resultV["height"]
    end
    r["clear"] = function()
        sendRPC("clear", nil)
    end
    r["clearLine"] = function()
        sendRPC("clearLine", nil)
    end
    r["getTextColor"] = function()
        sendRPC("getTextColor", nil)
        local resultV = awaitRPCResult("getTextColor")
        if resultV == nil then
            return localTerm.getTextColor()
        end
        return resultV["textColor"]
    end
    r["getTextColour"] = r["getTextColor"]
    r["setTextColor"] = function(textColor)
        sendRPC("setTextColor", {
            textColor = textColor
        })
    end
    r["setTextColour"] = r["setTextColor"]
    r["getBackgroundColor"] = function()
        sendRPC("getBackgroundColor", nil)
        local resultV = awaitRPCResult("getBackgroundColor")
        if resultV == nil then
            return localTerm.getBackgroundColor()
        end
        return resultV["backgroundColor"]
    end
    r["getBackgroundColour"] = r["getBackgroundColor"]
    r["setBackgroundColor"] = function(backgroundColor)
        sendRPC("setBackgroundColor", {
            backgroundColor = backgroundColor
        })

    end
    r["setBackgroundColour"] = r["setBackgroundColor"]
    r["isColor"] = function()
        sendRPC("isColor", nil)
        local resultV = awaitRPCResult("isColor")
        return resultV["isColor"] == true
    end
    r["isColour"] = r["isColor"]
    r["blit"] = function(text, textColou, backgroundColor)
        sendRPC("blit", {
            text = text,
            textColor = textColor,
            backgroundColor = backgroundColor,
        })
    end
    r["setPaletteColor"] = function(...)
        local colorID = arg[1]
        local colorCode
        if #arg == 4 then
            colorCode = colors.packRGB(arg[2], arg[3], arg[4])
        else
            colorCode = arg[2]
        end

        sendRPC("setPaletteColor", {
            colorID = colorID,
            colorCode = colorCode,
        })
    end
    r["setPaletteColour"] = r["setPaletteColor"]
    r["getPaletteColor"] = function(colorID)
        sendRPC("getPaletteColor", {
            colorID = colorID
        })
        local resultV = awaitRPCResult("getPaletteColor")
        if resultV == nil then
            return localTerm.getPaletteColor(colorID)
        end

        return colors.unpackRGB(resultV[1])
    end
    r["getPaletteColour"] = r["getPaletteColor"]
end

return {
    RShellTermRedirect = RShellTermRedirect
}
