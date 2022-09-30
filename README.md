# Multi-Channel Chat App / File System Emulator

A fun little application that emulates a multiplayer "filesystem" and real time chat.

# Quickstart

If you have wasmCloud already running (if not, download [wash](https://github.com/wasmcloud/wash) and run `wash up --nats-connect-only`), you can get started quickly with the prebuilt OCI artifacts.

## Requirements
- `wash` >= v0.12.0
- NATs server running with included `nats.config` (to enable the servers websocket)
- [Vite](https://www.vitejs.dev)
- Redis server running locally, below assumes no authentication
> Recommend `docker run --name redis -p 6379:6379 --rm -d redis` for easy start


### Actors
UI - `ghcr.io/jordan-rash/wc-terminal/ui:0.1.0`  
File System - `ghcr.io/jordan-rash/wc-terminal/nary_fs:0.1.0`  
Broker - `ghcr.io/jordan-rash/wc-terminal/broker:0.1.0`  

### Providers 
HTTPServer - `wasmcloud.azurecr.io/httpserver:0.16.0`  
NATs Messaging - `wasmcloud.azurecr.io/nats_messaging:0.14.5`  
KV Redis - `wasmcloud.azurecr.io/kvredis:0.16.4`  

### Link Definations
UI <-> HTTPServer (no configuration)  
File System <-> KV Redis (no configuration)  
Broker <-> NATs Messaging (SUBSCRIPTION=session.*.command)  


# Build

To build the project, you should be able to `cd` into all three directories and run `make`.  Artifacts will be dropped in a `build/` directory

# Fun with NGS

The UI interacts with the rest of the platform via the [NATs websocket library](https://github.com/nats-io/nats.ws). Following the quickstart, you will have multiplayer, but only locally.  If you would like to harness the dope power of the NATs Global Supercluster ([NGS](https://app.ngs.global/)), then the configuration of your Broker/NATs Messaging link will be a little different.

You will have 4 settings: (in the wasmCloud UI, they will be comma seperated)
```
SUBSCRIPTION=session.*.command
URI=connect.ngs.global:4222
CLIENT_SEED=SUANKIRFN7CUSD93M6OSJNWAJPIQZ6GJP7AWDX6ITS4MF6T3TCI3RFBI64
CLIENT_JWT=eyJ0eXAiOiJKV1QiLCJhbGciOiJlZDI1NTE5LW5rZXkifQ.eyJqdGkiOiJaUktNWEFZSE01Q0NBVzVXNjJOMlFXUjZTVlEzRlpTM0tPNFRKWlFWWVNIVUU1Uk1RNzdBIiwiaWF0IjoxNjY0MzA1ODkzLCJpc3MiOiJBQ04zMjJEVDJBNjRaQ0FaRUVMSFNIS1o3QUNKUFYzSjRTRTdWUFMzNDdKNkpDUkY2SEtTWEtTUSIsIm5hbWUiOiJ0ZXJtaW5hbF93cyIsInN1YiI6IlVBRFNNNzcyRUMyV0lIMklRM1FTVk5PVFhMWVk1S1ZMVkJORFhURTVFWTU2TDIzTTVFVkxCWEs1IbmF0cyI6eyJwdWIiOnsiYWxsb3ciOlsiJEpTLkFQSS5DT05TVU1FUi5DUkVBVEUudGVybWluYWwiLCIkSlMuQVBJLlNUUkVBTS5OQU1FUyIsIi5BUEkuU1RSRUFNLk5BTUVTIiwiXHUwMDNlIiwiSlMuQVBJLlNUUkVBTS5OQU1FUyIsIl9JTkJPWC5cdTAwM2UiLCJfaW5ib3guXHUwMDNlIiwic2Vzc2lvbi4qIiwic2Vzc2lvbi4qLmNvbW1hbmQiXX0sInN1YiI6eyJhbGxvdyI6WyJcdTAwM2UiLCJfSU5CT1guKiIsIl9JTkJPWC5cdTAwM2UiLCJfaW5ib3guXHUwMDNlIiwic2Vzc2lvbi4qIl19LCJyZXNwIjp7Im1heCI6MSwidHRsIjowfSwic3VicyI6LTEsImRhdGEiOi0xLCJwYXlsb2FkIjotMSwidHlwZSI6InVzZXIiLCJ2ZXJzaW9uIjoyfX0.4JrX7cI9a6J7gOeX9fYOPMAbI1JksmIsQplUhIEi5O5aDV_CeBvlbYL8z75GEiMEU-VEXPrrBfbnfDxyDhxhBg
```

> Note: This Seed and JWT is an example, you will need to generate your own.  Go to `https://app.ngs.global`, and sign up for a free account.  Follow their documentation for getting connected locally.  Once that is completed, you can run the `nsc generate creds` command to print your `CLIENT_SEED` and `CLIENT_JWT` to screen

> Note: You will need to update the User credentials in the UI's `main.js` as well.  I couldnt figure out how to make a multi-line variable work in the environment files

# How to play
Once you navigate to the app, you will see a generic splash screen that will give you the option to start a session that will generate a random session code.  At this point, you can share that URL with a friend for multiplayer, or not, its still fun :-)

The default mode you will be placed in is command mode.  This will allow you to interact with the filesystem.  You can switch to chat mode by using the slash command `/chat`.  Once in chat mode, it works like any other chat UI.  You will be given a random user name, you can update that with the `/nick myname` slash command.  

_*TL;DR*_  
`/cmd` -> Go to command mode.  Interacts with filessystem  
`ls` -> List nodes at current path  
`mkdir` -> Create directory type node at current path  
`touch` -> Create file type node at current path  
`cd` -> Navigate up or down a single node at a time. Relative to current path  
`rm` -> Remove node at current path  

`/chat` -> Go to chat mode.  
`/nick <name>` -> Change your name displaed in chat mode  
