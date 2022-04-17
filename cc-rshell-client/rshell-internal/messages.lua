local function BuildMessage(type, payload)
    payload["type"] = type
    return payload
end

local function BuildActivateMessage(localTerm)
    local keyCodes = {}
    for key, value in pairs(keys) do
        if type(value) == "number" then
            keyCodes[key] = value
        end
    end

    local nativeColors = {}
    for key, value in pairs(colors) do
        if type(value) == "number" then
            nativeColors[colors.toBlit(value)] = {
                label = key,
                colorID = value,
                colorCode = colors.packRGB(localTerm.getPaletteColor(value))
            }
        end
    end

    local label = os.getComputerLabel()
    if label == nil then
        label = ""
    end

    local payload = {
        id = os.getComputerID(),
        label = label,
        keyCodes = keyCodes,
        colors = nativeColors,
    }

    return BuildMessage("activate", payload)
end

local function BuildBufferUpdateMessage(procID, buffer)
    return BuildMessage("framebuffer", {
        procID = procID,
        buffer = buffer,
    })
end

return {
    BuildActivateMessage = BuildActivateMessage,
    BuildBufferUpdateMessage = BuildBufferUpdateMessage
}
