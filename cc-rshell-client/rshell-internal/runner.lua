
local function Runner(localTerm, program)
    multishell.setFocus(multishell.launch({}, program))
end

return {
    Runner=Runner
}