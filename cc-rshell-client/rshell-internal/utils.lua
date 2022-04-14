local function termPrint(target, s)
    local old = term.current()
    term.redirect(target)
    print(s)
    term.redirect(old)
end

return {
    termPrint = termPrint
}