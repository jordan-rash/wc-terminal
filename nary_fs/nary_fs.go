package main

import (
	"encoding/base64"
	"errors"

	"github.com/wasmcloud/actor-tinygo"
	kv "github.com/wasmcloud/interfaces/keyvalue/tinygo"
	logging "github.com/wasmcloud/interfaces/logging/tinygo"
	rand "github.com/wasmcloud/interfaces/numbergen/tinygo"
)

var log = logging.NewProviderLogging()

const (
	NODE_TYPE_FILE      NodeType = "file"
	NODE_TYPE_DIRECTORY NodeType = "directory"
)

func main() {
	me := NaryFsActor{}
	actor.RegisterHandlers(
		FsSubscriberHandler(&me),
	)
}

type NaryFsActor struct {
	fs NaryFs
}

func newNaryFs(ctx *actor.Context, name string) *NaryFs {
	root := Node{Id: "0", Key: "/", Parent: "", Children: nil}
	randNumGen := rand.NewProviderNumberGen()
	guid, _ := randNumGen.GenerateGuid(ctx)

	docs := Node{
		Id:       guid,
		Key:      "docs",
		Type:     string(NODE_TYPE_DIRECTORY),
		Parent:   "0",
		Children: new(Nodes),
	}
	guid, _ = randNumGen.GenerateGuid(ctx)
	pics := Node{
		Id:       guid,
		Key:      "pics",
		Type:     string(NODE_TYPE_DIRECTORY),
		Parent:   "0",
		Children: new(Nodes),
	}

	root.Children = &Nodes{docs.Id: docs, pics.Id: pics}

	return &NaryFs{
		Root: &root,
		Name: name,
	}
}

// when broker gets a ? message
func (c *NaryFsActor) HandleFsStatus(
	ctx *actor.Context,
) (*HandlerFsStatus, error) {
	cmds := Commands{
		{Name: "ls", Usage: ""},
		{Name: "mkdir", Usage: ""},
		{Name: "touch", Usage: ""},
		{Name: "rm", Usage: ""},
		{Name: "mv", Usage: ""},
	}
	ret := &HandlerFsStatus{
		Commands: &cmds,
	}

	return ret, nil
}

// when broker gets a ls,mkdir,touch,rm,mv message
func (c *NaryFsActor) HandleFsMessage(
	ctx *actor.Context,
	msg FsMsg,
) (*HandlerFsResponse, error) {
	log.WriteLog(
		ctx,
		logging.LogEntry{Level: "debug", Text: "Running FS actor"},
	)

	randNumGen := rand.NewProviderNumberGen()

	fsKey := msg.Session + "." + msg.Fsname

	kvclient := kv.NewProviderKeyValue()
	exists, err := kvclient.Contains(ctx, fsKey)
	if err != nil {
		return nil, err
	}

	if exists {
		log.WriteLog(
			ctx,
			logging.LogEntry{Level: "debug", Text: "Using existing redis key"},
		)
		resp, err := kvclient.Get(ctx, fsKey)
		if err != nil {
			return nil, err
		}
		sd, err := base64.StdEncoding.DecodeString(resp.Value)
		if err != nil {
			return nil, err
		}
		fs, _ := DecodeFS(sd)

		c.fs = *fs
	} else {
		log.WriteLog(ctx, logging.LogEntry{Level: "debug", Text: "Creating new redis key"})
		fs := newNaryFs(ctx, msg.Fsname)
		c.fs = *fs
	}

	// do action that inside msg.Body
	resp := HandlerFsResponse{
		Currnodeid: "0",
		Success:    false,
		Error:      "",
		Response:   "{\"data\":[]}",
		Abspath:    "/",
	}
	c.fs.Root.PrintLog(ctx)

	node, err := c.fs.Root.FindNode(ctx, msg.Payload.Nodeid)
	if err != nil {
		resp.Error = Error(err.Error())
		return &resp, nil
	}
	resp.Currnodeid = node.Id

	switch msg.Action {
	case "ls":
		resp.Response = c.fs.String(ctx, node)
		resp.Abspath = node.GetAbsPath(ctx, c.fs.Root)
		resp.Success = true
	case "cd":
		if msg.Payload.EventNodeKey == ".." {
			if node.Parent == "0" {
				resp.Currnodeid = "0"
				resp.Abspath = "/"
			} else {
				node, err := c.fs.Root.FindNode(ctx, node.Parent)
				if err != nil {
					return nil, err
				}
				resp.Currnodeid = node.Id
				resp.Abspath = node.GetAbsPath(ctx, c.fs.Root)
			}
			resp.Success = true
			break
		}
		for _, n := range *node.Children {
			if n.Key == msg.Payload.EventNodeKey &&
				n.Type == string(NODE_TYPE_DIRECTORY) {
				log.WriteLog(
					ctx,
					logging.LogEntry{
						Level: "info",
						Text:  "CDing into " + n.Key + "|" + n.Id,
					},
				)
				resp.Currnodeid = n.Id
				resp.Abspath = n.GetAbsPath(ctx, c.fs.Root)
				resp.Success = true
				break
			}
		}

		if !resp.Success {
			resp.Error = "Did not find directory"
		}
	case "mkdir":
		for _, n := range *node.Children {
			if n.Key == msg.Payload.EventNodeKey {
				resp.Error = "Node name already exists"
				break
			}
		}
		if resp.Error == "" {
			guid, err := randNumGen.GenerateGuid(ctx)
			if err != nil {
				return nil, err
			}
			node.AddNode(guid, msg.Payload.EventNodeKey, NODE_TYPE_DIRECTORY)
			resp.Abspath = node.GetAbsPath(ctx, c.fs.Root)
			resp.Success = true
		}
	case "touch":
		for _, n := range *node.Children {
			if n.Key == msg.Payload.EventNodeKey {
				resp.Error = "Node name already exists"
				break
			}
		}
		if resp.Error == "" {
			guid, err := randNumGen.GenerateGuid(ctx)
			if err != nil {
				return nil, err
			}
			node.AddNode(guid, msg.Payload.EventNodeKey, NODE_TYPE_FILE)
			resp.Abspath = node.GetAbsPath(ctx, c.fs.Root)
			resp.Success = true
		}
	case "rm":
		for _, n := range *node.Children {
			if n.Key == msg.Payload.EventNodeKey {
				if n.Type == string(NODE_TYPE_DIRECTORY) &&
					len(*n.Children) > 0 {
					resp.Error = "Can not delete a directory with children"
					break
				}
				node.DeleteNode(n.Id)
				resp.Abspath = node.GetAbsPath(ctx, c.fs.Root)
				resp.Success = true
				break
			}
		}
		if !resp.Success {
			resp.Error = "Was not able to delete node"
		}
	//case "mv":
	//newPNode, err := c.fs.Root.GetNode(msg.Payload.NewNodeParent)
	//if err != nil {
	//return nil, err
	//}
	//node.MoveNode(msg.Payload.Nodeid, newPNode)
	default:
		resp.Error = "Invalid action"
	}

	// encode and save fs to key store
	efs := EncodeFS(c.fs)
	se := base64.StdEncoding.EncodeToString(efs)
	k := kv.SetRequest{
		Key:   fsKey,
		Value: se,
	}

	err = kvclient.Set(ctx, k)
	if err != nil {
		return nil, err
	}

	// return success msg
	return &resp, nil
}

func (n *Node) FindNode(ctx *actor.Context, nodeid string) (*Node, error) {
	log.WriteLog(
		ctx,
		logging.LogEntry{
			Level: "info",
			Text:  "Looking for " + nodeid + " on " + n.Id,
		},
	)
	if n.Id == nodeid {
		return n, nil
	}

	for _, c := range *n.Children {
		nn, _ := c.FindNode(ctx, nodeid)
		if nn != nil {
			return nn, nil
		}
	}

	return nil, errors.New("Node ID not found")
}

func (n *Node) AddNode(id, nodekey string, nodetype NodeType) error {
	nn := Node{
		Id:       id,
		Type:     string(nodetype),
		Key:      nodekey,
		Children: new(Nodes),
		Parent:   n.Id,
	}

	children := *n.Children
	children[nn.Id] = nn

	n.Children = &children
	return nil
}

func (n *Node) DeleteNode(cNodeId string) error {
	children := *n.Children
	cNode, ok := children[cNodeId]
	if !ok {
		return errors.New("Child node does not exist")
	}

	if cNode.Type != string(NODE_TYPE_FILE) {
		return errors.New("Must use force flag to delete node")
	}

	delete(children, cNode.Id)

	n.Children = &children
	return nil
}

func (n *Node) MoveNode(cNodeId string, newNodeParent *Node) error {
	children := *n.Children
	cNode, ok := children[cNodeId]
	if !ok {
		return errors.New("Child node does not exist")
	}
	n.DeleteNode(cNodeId)

	npChildren := *newNodeParent.Children
	npChildren[cNode.Id] = cNode

	newNodeParent.Children = &npChildren

	return nil
}

func (fs NaryFs) String(ctx *actor.Context, node *Node) string {
	if node == nil {
		return "{}"
	}

	ret := "{\"data\":["
	nodes := *node.Children
	count := 0
	for _, n := range nodes {
		ret += "{\"ty\":\"" + n.Type + "\","
		ret += "\"name\":\"" + n.Key + "\"}"
		count++
		if count < len(nodes) {
			ret += ","
		}
	}
	ret += "]}"

	log.WriteLog(ctx, logging.LogEntry{Level: "debug", Text: ret})
	return ret
}

func (n Node) PrintLog(ctx *actor.Context) {
	log.WriteLog(ctx, logging.LogEntry{Level: "info", Text: n.Key + "|" + n.Id})
	for _, c := range *n.Children {
		c.PrintLog(ctx)
	}
}

func (n Node) GetAbsPath(ctx *actor.Context, root *Node) string {
	if n.Id == "0" {
		return "/"
	}

	nn, _ := root.FindNode(ctx, n.Parent)
	return nn.GetAbsPath(ctx, root) + n.Key + "/"
}
