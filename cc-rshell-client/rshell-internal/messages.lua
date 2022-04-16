local function BuildMessage(type, payload)
    payload["type"] = type
    return payload
end

local function BuildActivateMessage()
    local keyCodes = {}
    for key, value in pairs(keys) do
        if type(value) == "number" then
            keyCodes[key] = value
        end
    end

    local label = os.getComputerLabel()
    if label == nil then
        label = ""
    end

    local payload = {
        id = os.getComputerID(),
        label = label,
        keyCodes = keyCodes
    }

    return BuildMessage("activate", payload)
end

return {
    BuildActivateMessage = BuildActivateMessage
}