print("Welcome to ECHO!")

while true do
    local r = read()
    print("User said: " .. r)
    if r == "q" then
        break
    end
end