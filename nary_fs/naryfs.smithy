metadata package = [ { namespace: "org.jordanrash.naryfs", crate: "naryfs" } ]

namespace org.jordanrash.naryfs

use org.wasmcloud.model#wasmbus
use org.wasmcloud.model#I64

@wasmbus(
    contractId: "jordanrash:terminal:fs",
    actorReceive: true )
service FsSubscriber {
  version: "0.1",
  operations: [ HandleFsMessage, HandleFsStatus ]
}

operation HandleFsMessage {
    input: FsMsg
    output: HandlerFsResponse
}

operation HandleFsStatus {
    output: HandlerFsStatus
}

string NodeType
string Error

structure HandlerFsStatus {
   commands: Commands, 
}

list Commands {
    member: Command
}

structure Command {
    name: String,
    usage: String
}

structure Node {
   id: String,
   type: String,
   key: String,
   children: Nodes,
   parent: String
}

structure HandlerFsResponse {
    error: Error,
    success: Boolean,
    currnodeid: String,
    response: String,
    abspath: String
}

structure NaryFs {
    root: Node,
    name: String
}

map Nodes {
    key: String,
    value: Node
}

structure FsMsg {
    action: String,
    fsname: String,
    session: String,
    payload: Payload
}

structure Payload {
    nodeid: String,
    eventNodeKey: String
}
