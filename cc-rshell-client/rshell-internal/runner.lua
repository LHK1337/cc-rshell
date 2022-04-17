local function Runner(localTerm, id, program, ...)
    multishell.setFocus(multishell.launch(_ENV, "rshell-runner.lua", id, program, ...))
end

return {
    Runner=Runner
}