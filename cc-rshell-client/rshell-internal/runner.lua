local redirect = require "rshell-internal.redirect"

local function Runner(localTerm, program)
    term.redirect(redirect.RShellTermRedirect(ws, localTerm))
    multishell.setFocus(multishell.launch({}, program))
end

return {
    Runner=Runner
}