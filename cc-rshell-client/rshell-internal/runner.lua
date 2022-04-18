local procTable = {}

local function ProcTableAdd(procID, msID)
    procTable[procID] = msID
end

local function ProcTableRemove(procID)
    for pID, msID in pairs(procTable) do
        if msID > procTable[procID] then
            procTable[pID] = procTable[pID] - 1
        end
    end

    procTable[procID] = nil
end

local function Runner(procID, bufH, bufW, program, ...)
    local msID = multishell.launch(_ENV, "rshell-runner.lua", procID, bufH, bufW, program, ...)
    multishell.setFocus(msID)
    ProcTableAdd(procID, msID)
end

local function Focus(procID)
    if procTable[procID] then
        multishell.setFocus(procTable[procID])
        return true
    end
    return false
end

return {
    Runner = Runner,
    Focus = Focus,
    ProcTableAdd = ProcTableAdd,
    ProcTableRemove = ProcTableRemove,
}
