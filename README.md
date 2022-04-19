# cc-rshell - ComputerCraft RemoteShell

✨ The fancy SSH solution you did not ask for. ✨

## What?

A gateway between our universe and a minecraft world of your choice to access your
[ComputerCraft](https://github.com/cc-tweaked/CC-Tweaked) Computer's shell via SSH.

## Why?

Why not?

## How to?

- This relies on the `multishell` api, therefore
  an [advanced computer](http://www.computercraft.info/wiki/Advanced_Computer) is required

1. Copy everything from [`cc-rshell-client`](cc-rshell-client) to your minecraft computer's root
2. Change the server url to your url in [socketService.lua](cc-rshell-client/rshell-internal/socketService.lua)
3. Run the server `cc-rshell-server` somewhere reachable for your minecraft computer.
    1. _You may have to tweak your ComputerCraft config to allow websocket connection to your server_
4. Run the client `cc-rshell-client` with `rshell`
5. SSH into your minecraft machine with `ssh <COMPUTER_ID>@my_server.org -p<SERVER_PORT>`
   or `ssh <COMPUTER_LABEL>@my_server.org -p<SERVER_PORT>`

## ✨ What works? ✨

Everything except things that do not and things I did not encounter.

## 🐛 What does not work? 🐛

- Programs that rely on the keys `CTRL, SHIFT, ALT, SUPER/WIN`. (Key combos do work again)
    - Example: The default editor `edit` expects just `CTRL` to access the menu
    - This can be fixed by remapping these keys to other keys which appear to work just fine
    - My uneducated guess here is, that the terminal or at least the terminal library is not able to pick up those keys
      because they do not produce a character in the terminal.

- Dynamic terminal resizing
    - The terminal will use the local terminal size you connect with

- Mouse input
    - because I have not implemented it (yet)
    - But I think this _**should not**_ be that hard 🤞

## 🤞 What might work some time?🤞

- Actual login with a password or an SSH Key to get closer to a secure shell
- Prevention of computer id spoofing to get closer to a secure shell
- Dynamic shell resizing
- Mouse input

### 🤫 Absurd shenanigans 🤫

##### RMTS™ (Rapid Multishell Tab Switching)

Key presses reach their destination through `char`, `key` and `key_up` events and only the current selected tab in a
multishell receives these events. So what happens when you summon
_multiple_ shells by opening _multiple_ SSH connections?

##### MBMCP™ (Multi Buffer Message Chunking Protocol)

Imagine a world where web socket messages are limited by 512 bytes...

- See [this](cc-rshell-server/sockets/messages/messages.go) and [that](cc-rshell-client/rshell-internal/utils.lua)

# 💛 Special thanks to 💛

- @cc-tweaked, @dan200 and all contributors for this amazing minecraft mod
- @markstinson and @fperrad for implementing messagepack in lua
- @Lyqyd for the framebuffer api
- @gliderlabs for their handy ssh library
- @Gin-Gonic because gin is awesome
- @CCEmuX for their emulator