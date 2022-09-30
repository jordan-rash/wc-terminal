package main

import (
	"github.com/wasmcloud/actor-tinygo"           //nolint
	cbor "github.com/wasmcloud/tinygo-cbor"       //nolint
	msgpack "github.com/wasmcloud/tinygo-msgpack" //nolint
)

type Command struct {
	Name  string
	Usage string
}

// MEncode serializes a Command using msgpack
func (o *Command) MEncode(encoder msgpack.Writer) error {
	encoder.WriteMapSize(2)
	encoder.WriteString("name")
	encoder.WriteString(o.Name)
	encoder.WriteString("usage")
	encoder.WriteString(o.Usage)

	return encoder.CheckError()
}

// MDecodeCommand deserializes a Command using msgpack
func MDecodeCommand(d *msgpack.Decoder) (Command, error) {
	var val Command
	isNil, err := d.IsNextNil()
	if err != nil || isNil {
		return val, err
	}
	size, err := d.ReadMapSize()
	if err != nil {
		return val, err
	}
	for i := uint32(0); i < size; i++ {
		field, err := d.ReadString()
		if err != nil {
			return val, err
		}
		switch field {
		case "name":
			val.Name, err = d.ReadString()
		case "usage":
			val.Usage, err = d.ReadString()
		default:
			err = d.Skip()
		}
		if err != nil {
			return val, err
		}
	}
	return val, nil
}

// CEncode serializes a Command using cbor
func (o *Command) CEncode(encoder cbor.Writer) error {
	encoder.WriteMapSize(2)
	encoder.WriteString("name")
	encoder.WriteString(o.Name)
	encoder.WriteString("usage")
	encoder.WriteString(o.Usage)

	return encoder.CheckError()
}

// CDecodeCommand deserializes a Command using cbor
func CDecodeCommand(d *cbor.Decoder) (Command, error) {
	var val Command
	isNil, err := d.IsNextNil()
	if err != nil || isNil {
		return val, err
	}
	size, indef, err := d.ReadMapSize()
	if err != nil && indef {
		err = cbor.NewReadError("indefinite maps not supported")
	}
	if err != nil {
		return val, err
	}
	for i := uint32(0); i < size; i++ {
		field, err := d.ReadString()
		if err != nil {
			return val, err
		}
		switch field {
		case "name":
			val.Name, err = d.ReadString()
		case "usage":
			val.Usage, err = d.ReadString()
		default:
			err = d.Skip()
		}
		if err != nil {
			return val, err
		}
	}
	return val, nil
}

type Commands []Command

// MEncode serializes a Commands using msgpack
func (o *Commands) MEncode(encoder msgpack.Writer) error {

	encoder.WriteArraySize(uint32(len(*o)))
	for _, item_o := range *o {
		item_o.MEncode(encoder)
	}

	return encoder.CheckError()
}

// MDecodeCommands deserializes a Commands using msgpack
func MDecodeCommands(d *msgpack.Decoder) (Commands, error) {
	isNil, err := d.IsNextNil()
	if err != nil || isNil {
		return make([]Command, 0), err
	}
	size, err := d.ReadArraySize()
	if err != nil {
		return make([]Command, 0), err
	}
	val := make([]Command, size)
	for i := uint32(0); i < size; i++ {
		item, err := MDecodeCommand(d)
		if err != nil {
			return val, err
		}
		val = append(val, item)
	}
	return val, nil
}

// CEncode serializes a Commands using cbor
func (o *Commands) CEncode(encoder cbor.Writer) error {

	encoder.WriteArraySize(uint32(len(*o)))
	for _, item_o := range *o {
		item_o.CEncode(encoder)
	}

	return encoder.CheckError()
}

// CDecodeCommands deserializes a Commands using cbor
func CDecodeCommands(d *cbor.Decoder) (Commands, error) {
	isNil, err := d.IsNextNil()
	if err != nil || isNil {
		return make([]Command, 0), err
	}
	size, indef, err := d.ReadArraySize()
	if err != nil && indef {
		err = cbor.NewReadError("indefinite arrays not supported")
	}
	if err != nil {
		return make([]Command, 0), err
	}
	val := make([]Command, size)
	for i := uint32(0); i < size; i++ {
		item, err := CDecodeCommand(d)
		if err != nil {
			return val, err
		}
		val = append(val, item)
	}
	return val, nil
}

type Error string

// MEncode serializes a Error using msgpack
func (o *Error) MEncode(encoder msgpack.Writer) error {
	encoder.WriteString(string(*o))
	return encoder.CheckError()
}

// MDecodeError deserializes a Error using msgpack
func MDecodeError(d *msgpack.Decoder) (Error, error) {
	val, err := d.ReadString()
	if err != nil {
		return "", err
	}
	return Error(val), nil
}

// CEncode serializes a Error using cbor
func (o *Error) CEncode(encoder cbor.Writer) error {
	encoder.WriteString(string(*o))
	return encoder.CheckError()
}

// CDecodeError deserializes a Error using cbor
func CDecodeError(d *cbor.Decoder) (Error, error) {
	val, err := d.ReadString()
	if err != nil {
		return "", err
	}
	return Error(val), nil
}

type FsMsg struct {
	Action  string
	Fsname  string
	Payload *Payload
	Session string
}

// MEncode serializes a FsMsg using msgpack
func (o *FsMsg) MEncode(encoder msgpack.Writer) error {
	encoder.WriteMapSize(4)
	encoder.WriteString("action")
	encoder.WriteString(o.Action)
	encoder.WriteString("fsname")
	encoder.WriteString(o.Fsname)
	encoder.WriteString("payload")
	if o.Payload == nil {
		encoder.WriteNil()
	} else {
		o.Payload.MEncode(encoder)
	}
	encoder.WriteString("session")
	encoder.WriteString(o.Session)

	return encoder.CheckError()
}

// MDecodeFsMsg deserializes a FsMsg using msgpack
func MDecodeFsMsg(d *msgpack.Decoder) (FsMsg, error) {
	var val FsMsg
	isNil, err := d.IsNextNil()
	if err != nil || isNil {
		return val, err
	}
	size, err := d.ReadMapSize()
	if err != nil {
		return val, err
	}
	for i := uint32(0); i < size; i++ {
		field, err := d.ReadString()
		if err != nil {
			return val, err
		}
		switch field {
		case "action":
			val.Action, err = d.ReadString()
		case "fsname":
			val.Fsname, err = d.ReadString()
		case "payload":
			fval, err := MDecodePayload(d)
			if err != nil {
				return val, err
			}
			val.Payload = &fval
		case "session":
			val.Session, err = d.ReadString()
		default:
			err = d.Skip()
		}
		if err != nil {
			return val, err
		}
	}
	return val, nil
}

// CEncode serializes a FsMsg using cbor
func (o *FsMsg) CEncode(encoder cbor.Writer) error {
	encoder.WriteMapSize(4)
	encoder.WriteString("action")
	encoder.WriteString(o.Action)
	encoder.WriteString("fsname")
	encoder.WriteString(o.Fsname)
	encoder.WriteString("payload")
	if o.Payload == nil {
		encoder.WriteNil()
	} else {
		o.Payload.CEncode(encoder)
	}
	encoder.WriteString("session")
	encoder.WriteString(o.Session)

	return encoder.CheckError()
}

// CDecodeFsMsg deserializes a FsMsg using cbor
func CDecodeFsMsg(d *cbor.Decoder) (FsMsg, error) {
	var val FsMsg
	isNil, err := d.IsNextNil()
	if err != nil || isNil {
		return val, err
	}
	size, indef, err := d.ReadMapSize()
	if err != nil && indef {
		err = cbor.NewReadError("indefinite maps not supported")
	}
	if err != nil {
		return val, err
	}
	for i := uint32(0); i < size; i++ {
		field, err := d.ReadString()
		if err != nil {
			return val, err
		}
		switch field {
		case "action":
			val.Action, err = d.ReadString()
		case "fsname":
			val.Fsname, err = d.ReadString()
		case "payload":
			fval, err := CDecodePayload(d)
			if err != nil {
				return val, err
			}
			val.Payload = &fval
		case "session":
			val.Session, err = d.ReadString()
		default:
			err = d.Skip()
		}
		if err != nil {
			return val, err
		}
	}
	return val, nil
}

type HandlerFsResponse struct {
	Abspath    string
	Currnodeid string
	Error      Error
	Response   string
	Success    bool
}

// MEncode serializes a HandlerFsResponse using msgpack
func (o *HandlerFsResponse) MEncode(encoder msgpack.Writer) error {
	encoder.WriteMapSize(5)
	encoder.WriteString("abspath")
	encoder.WriteString(o.Abspath)
	encoder.WriteString("currnodeid")
	encoder.WriteString(o.Currnodeid)
	encoder.WriteString("error")
	o.Error.MEncode(encoder)
	encoder.WriteString("response")
	encoder.WriteString(o.Response)
	encoder.WriteString("success")
	encoder.WriteBool(o.Success)

	return encoder.CheckError()
}

// MDecodeHandlerFsResponse deserializes a HandlerFsResponse using msgpack
func MDecodeHandlerFsResponse(d *msgpack.Decoder) (HandlerFsResponse, error) {
	var val HandlerFsResponse
	isNil, err := d.IsNextNil()
	if err != nil || isNil {
		return val, err
	}
	size, err := d.ReadMapSize()
	if err != nil {
		return val, err
	}
	for i := uint32(0); i < size; i++ {
		field, err := d.ReadString()
		if err != nil {
			return val, err
		}
		switch field {
		case "abspath":
			val.Abspath, err = d.ReadString()
		case "currnodeid":
			val.Currnodeid, err = d.ReadString()
		case "error":
			val.Error, err = MDecodeError(d)
		case "response":
			val.Response, err = d.ReadString()
		case "success":
			val.Success, err = d.ReadBool()
		default:
			err = d.Skip()
		}
		if err != nil {
			return val, err
		}
	}
	return val, nil
}

// CEncode serializes a HandlerFsResponse using cbor
func (o *HandlerFsResponse) CEncode(encoder cbor.Writer) error {
	encoder.WriteMapSize(5)
	encoder.WriteString("abspath")
	encoder.WriteString(o.Abspath)
	encoder.WriteString("currnodeid")
	encoder.WriteString(o.Currnodeid)
	encoder.WriteString("error")
	o.Error.CEncode(encoder)
	encoder.WriteString("response")
	encoder.WriteString(o.Response)
	encoder.WriteString("success")
	encoder.WriteBool(o.Success)

	return encoder.CheckError()
}

// CDecodeHandlerFsResponse deserializes a HandlerFsResponse using cbor
func CDecodeHandlerFsResponse(d *cbor.Decoder) (HandlerFsResponse, error) {
	var val HandlerFsResponse
	isNil, err := d.IsNextNil()
	if err != nil || isNil {
		return val, err
	}
	size, indef, err := d.ReadMapSize()
	if err != nil && indef {
		err = cbor.NewReadError("indefinite maps not supported")
	}
	if err != nil {
		return val, err
	}
	for i := uint32(0); i < size; i++ {
		field, err := d.ReadString()
		if err != nil {
			return val, err
		}
		switch field {
		case "abspath":
			val.Abspath, err = d.ReadString()
		case "currnodeid":
			val.Currnodeid, err = d.ReadString()
		case "error":
			val.Error, err = CDecodeError(d)
		case "response":
			val.Response, err = d.ReadString()
		case "success":
			val.Success, err = d.ReadBool()
		default:
			err = d.Skip()
		}
		if err != nil {
			return val, err
		}
	}
	return val, nil
}

type HandlerFsStatus struct {
	Commands *Commands
}

// MEncode serializes a HandlerFsStatus using msgpack
func (o *HandlerFsStatus) MEncode(encoder msgpack.Writer) error {
	encoder.WriteMapSize(1)
	encoder.WriteString("commands")
	if o.Commands == nil {
		encoder.WriteNil()
	} else {
		o.Commands.MEncode(encoder)
	}

	return encoder.CheckError()
}

// MDecodeHandlerFsStatus deserializes a HandlerFsStatus using msgpack
func MDecodeHandlerFsStatus(d *msgpack.Decoder) (HandlerFsStatus, error) {
	var val HandlerFsStatus
	isNil, err := d.IsNextNil()
	if err != nil || isNil {
		return val, err
	}
	size, err := d.ReadMapSize()
	if err != nil {
		return val, err
	}
	for i := uint32(0); i < size; i++ {
		field, err := d.ReadString()
		if err != nil {
			return val, err
		}
		switch field {
		case "commands":
			fval, err := MDecodeCommands(d)
			if err != nil {
				return val, err
			}
			val.Commands = &fval
		default:
			err = d.Skip()
		}
		if err != nil {
			return val, err
		}
	}
	return val, nil
}

// CEncode serializes a HandlerFsStatus using cbor
func (o *HandlerFsStatus) CEncode(encoder cbor.Writer) error {
	encoder.WriteMapSize(1)
	encoder.WriteString("commands")
	if o.Commands == nil {
		encoder.WriteNil()
	} else {
		o.Commands.CEncode(encoder)
	}

	return encoder.CheckError()
}

// CDecodeHandlerFsStatus deserializes a HandlerFsStatus using cbor
func CDecodeHandlerFsStatus(d *cbor.Decoder) (HandlerFsStatus, error) {
	var val HandlerFsStatus
	isNil, err := d.IsNextNil()
	if err != nil || isNil {
		return val, err
	}
	size, indef, err := d.ReadMapSize()
	if err != nil && indef {
		err = cbor.NewReadError("indefinite maps not supported")
	}
	if err != nil {
		return val, err
	}
	for i := uint32(0); i < size; i++ {
		field, err := d.ReadString()
		if err != nil {
			return val, err
		}
		switch field {
		case "commands":
			fval, err := CDecodeCommands(d)
			if err != nil {
				return val, err
			}
			val.Commands = &fval
		default:
			err = d.Skip()
		}
		if err != nil {
			return val, err
		}
	}
	return val, nil
}

type NaryFs struct {
	Name string
	Root *Node
}

// MEncode serializes a NaryFs using msgpack
func (o *NaryFs) MEncode(encoder msgpack.Writer) error {
	encoder.WriteMapSize(2)
	encoder.WriteString("name")
	encoder.WriteString(o.Name)
	encoder.WriteString("root")
	if o.Root == nil {
		encoder.WriteNil()
	} else {
		o.Root.MEncode(encoder)
	}

	return encoder.CheckError()
}

// MDecodeNaryFs deserializes a NaryFs using msgpack
func MDecodeNaryFs(d *msgpack.Decoder) (NaryFs, error) {
	var val NaryFs
	isNil, err := d.IsNextNil()
	if err != nil || isNil {
		return val, err
	}
	size, err := d.ReadMapSize()
	if err != nil {
		return val, err
	}
	for i := uint32(0); i < size; i++ {
		field, err := d.ReadString()
		if err != nil {
			return val, err
		}
		switch field {
		case "name":
			val.Name, err = d.ReadString()
		case "root":
			fval, err := MDecodeNode(d)
			if err != nil {
				return val, err
			}
			val.Root = &fval
		default:
			err = d.Skip()
		}
		if err != nil {
			return val, err
		}
	}
	return val, nil
}

// CEncode serializes a NaryFs using cbor
func (o *NaryFs) CEncode(encoder cbor.Writer) error {
	encoder.WriteMapSize(2)
	encoder.WriteString("name")
	encoder.WriteString(o.Name)
	encoder.WriteString("root")
	if o.Root == nil {
		encoder.WriteNil()
	} else {
		o.Root.CEncode(encoder)
	}

	return encoder.CheckError()
}

// CDecodeNaryFs deserializes a NaryFs using cbor
func CDecodeNaryFs(d *cbor.Decoder) (NaryFs, error) {
	var val NaryFs
	isNil, err := d.IsNextNil()
	if err != nil || isNil {
		return val, err
	}
	size, indef, err := d.ReadMapSize()
	if err != nil && indef {
		err = cbor.NewReadError("indefinite maps not supported")
	}
	if err != nil {
		return val, err
	}
	for i := uint32(0); i < size; i++ {
		field, err := d.ReadString()
		if err != nil {
			return val, err
		}
		switch field {
		case "name":
			val.Name, err = d.ReadString()
		case "root":
			fval, err := CDecodeNode(d)
			if err != nil {
				return val, err
			}
			val.Root = &fval
		default:
			err = d.Skip()
		}
		if err != nil {
			return val, err
		}
	}
	return val, nil
}

type Node struct {
	Children *Nodes
	Id       string
	Key      string
	Parent   string
	Type     string
}

// MEncode serializes a Node using msgpack
func (o *Node) MEncode(encoder msgpack.Writer) error {
	encoder.WriteMapSize(5)
	encoder.WriteString("children")
	if o.Children == nil {
		encoder.WriteNil()
	} else {
		o.Children.MEncode(encoder)
	}
	encoder.WriteString("id")
	encoder.WriteString(o.Id)
	encoder.WriteString("key")
	encoder.WriteString(o.Key)
	encoder.WriteString("parent")
	encoder.WriteString(o.Parent)
	encoder.WriteString("type")
	encoder.WriteString(o.Type)

	return encoder.CheckError()
}

// MDecodeNode deserializes a Node using msgpack
func MDecodeNode(d *msgpack.Decoder) (Node, error) {
	var val Node
	isNil, err := d.IsNextNil()
	if err != nil || isNil {
		return val, err
	}
	size, err := d.ReadMapSize()
	if err != nil {
		return val, err
	}
	for i := uint32(0); i < size; i++ {
		field, err := d.ReadString()
		if err != nil {
			return val, err
		}
		switch field {
		case "children":
			fval, err := MDecodeNodes(d)
			if err != nil {
				return val, err
			}
			val.Children = &fval
		case "id":
			val.Id, err = d.ReadString()
		case "key":
			val.Key, err = d.ReadString()
		case "parent":
			val.Parent, err = d.ReadString()
		case "type":
			val.Type, err = d.ReadString()
		default:
			err = d.Skip()
		}
		if err != nil {
			return val, err
		}
	}
	return val, nil
}

// CEncode serializes a Node using cbor
func (o *Node) CEncode(encoder cbor.Writer) error {
	encoder.WriteMapSize(5)
	encoder.WriteString("children")
	if o.Children == nil {
		encoder.WriteNil()
	} else {
		o.Children.CEncode(encoder)
	}
	encoder.WriteString("id")
	encoder.WriteString(o.Id)
	encoder.WriteString("key")
	encoder.WriteString(o.Key)
	encoder.WriteString("parent")
	encoder.WriteString(o.Parent)
	encoder.WriteString("type")
	encoder.WriteString(o.Type)

	return encoder.CheckError()
}

// CDecodeNode deserializes a Node using cbor
func CDecodeNode(d *cbor.Decoder) (Node, error) {
	var val Node
	isNil, err := d.IsNextNil()
	if err != nil || isNil {
		return val, err
	}
	size, indef, err := d.ReadMapSize()
	if err != nil && indef {
		err = cbor.NewReadError("indefinite maps not supported")
	}
	if err != nil {
		return val, err
	}
	for i := uint32(0); i < size; i++ {
		field, err := d.ReadString()
		if err != nil {
			return val, err
		}
		switch field {
		case "children":
			fval, err := CDecodeNodes(d)
			if err != nil {
				return val, err
			}
			val.Children = &fval
		case "id":
			val.Id, err = d.ReadString()
		case "key":
			val.Key, err = d.ReadString()
		case "parent":
			val.Parent, err = d.ReadString()
		case "type":
			val.Type, err = d.ReadString()
		default:
			err = d.Skip()
		}
		if err != nil {
			return val, err
		}
	}
	return val, nil
}

type NodeType string

// MEncode serializes a NodeType using msgpack
func (o *NodeType) MEncode(encoder msgpack.Writer) error {
	encoder.WriteString(string(*o))
	return encoder.CheckError()
}

// MDecodeNodeType deserializes a NodeType using msgpack
func MDecodeNodeType(d *msgpack.Decoder) (NodeType, error) {
	val, err := d.ReadString()
	if err != nil {
		return "", err
	}
	return NodeType(val), nil
}

// CEncode serializes a NodeType using cbor
func (o *NodeType) CEncode(encoder cbor.Writer) error {
	encoder.WriteString(string(*o))
	return encoder.CheckError()
}

// CDecodeNodeType deserializes a NodeType using cbor
func CDecodeNodeType(d *cbor.Decoder) (NodeType, error) {
	val, err := d.ReadString()
	if err != nil {
		return "", err
	}
	return NodeType(val), nil
}

type Nodes map[string]Node

// MEncode serializes a Nodes using msgpack
func (o *Nodes) MEncode(encoder msgpack.Writer) error {
	encoder.WriteMapSize(uint32(len(*o)))
	for key_o, val_o := range *o {
		encoder.WriteString(key_o)
		val_o.MEncode(encoder)
	}

	return encoder.CheckError()
}

// MDecodeNodes deserializes a Nodes using msgpack
func MDecodeNodes(d *msgpack.Decoder) (Nodes, error) {
	isNil, err := d.IsNextNil()
	if err != nil || isNil {
		return make(map[string]Node, 0), err
	}
	size, err := d.ReadMapSize()
	if err != nil {
		return make(map[string]Node, 0), err
	}
	val := make(map[string]Node, size)
	for i := uint32(0); i < size; i++ {
		k, _ := d.ReadString()
		v, err := MDecodeNode(d)
		if err != nil {
			return val, err
		}
		val[k] = v
	}
	return val, nil
}

// CEncode serializes a Nodes using cbor
func (o *Nodes) CEncode(encoder cbor.Writer) error {
	encoder.WriteMapSize(uint32(len(*o)))
	for key_o, val_o := range *o {
		encoder.WriteString(key_o)
		val_o.CEncode(encoder)
	}

	return encoder.CheckError()
}

// CDecodeNodes deserializes a Nodes using cbor
func CDecodeNodes(d *cbor.Decoder) (Nodes, error) {
	isNil, err := d.IsNextNil()
	if err != nil || isNil {
		return make(map[string]Node, 0), err
	}
	size, indef, err := d.ReadMapSize()
	if err != nil && indef {
		err = cbor.NewReadError("indefinite maps not supported")
	}
	if err != nil {
		return make(map[string]Node, 0), err
	}
	val := make(map[string]Node, size)
	for i := uint32(0); i < size; i++ {
		k, _ := d.ReadString()
		v, err := CDecodeNode(d)
		if err != nil {
			return val, err
		}
		val[k] = v
	}
	return val, nil
}

type Payload struct {
	EventNodeKey string
	Nodeid       string
}

// MEncode serializes a Payload using msgpack
func (o *Payload) MEncode(encoder msgpack.Writer) error {
	encoder.WriteMapSize(2)
	encoder.WriteString("eventNodeKey")
	encoder.WriteString(o.EventNodeKey)
	encoder.WriteString("nodeid")
	encoder.WriteString(o.Nodeid)

	return encoder.CheckError()
}

// MDecodePayload deserializes a Payload using msgpack
func MDecodePayload(d *msgpack.Decoder) (Payload, error) {
	var val Payload
	isNil, err := d.IsNextNil()
	if err != nil || isNil {
		return val, err
	}
	size, err := d.ReadMapSize()
	if err != nil {
		return val, err
	}
	for i := uint32(0); i < size; i++ {
		field, err := d.ReadString()
		if err != nil {
			return val, err
		}
		switch field {
		case "eventNodeKey":
			val.EventNodeKey, err = d.ReadString()
		case "nodeid":
			val.Nodeid, err = d.ReadString()
		default:
			err = d.Skip()
		}
		if err != nil {
			return val, err
		}
	}
	return val, nil
}

// CEncode serializes a Payload using cbor
func (o *Payload) CEncode(encoder cbor.Writer) error {
	encoder.WriteMapSize(2)
	encoder.WriteString("eventNodeKey")
	encoder.WriteString(o.EventNodeKey)
	encoder.WriteString("nodeid")
	encoder.WriteString(o.Nodeid)

	return encoder.CheckError()
}

// CDecodePayload deserializes a Payload using cbor
func CDecodePayload(d *cbor.Decoder) (Payload, error) {
	var val Payload
	isNil, err := d.IsNextNil()
	if err != nil || isNil {
		return val, err
	}
	size, indef, err := d.ReadMapSize()
	if err != nil && indef {
		err = cbor.NewReadError("indefinite maps not supported")
	}
	if err != nil {
		return val, err
	}
	for i := uint32(0); i < size; i++ {
		field, err := d.ReadString()
		if err != nil {
			return val, err
		}
		switch field {
		case "eventNodeKey":
			val.EventNodeKey, err = d.ReadString()
		case "nodeid":
			val.Nodeid, err = d.ReadString()
		default:
			err = d.Skip()
		}
		if err != nil {
			return val, err
		}
	}
	return val, nil
}

type FsSubscriber interface {
	HandleFsMessage(ctx *actor.Context, arg FsMsg) (*HandlerFsResponse, error)
	HandleFsStatus(ctx *actor.Context) (*HandlerFsStatus, error)
}

// FsSubscriberHandler is called by an actor during `main` to generate a dispatch handler
// The output of this call should be passed into `actor.RegisterHandlers`
func FsSubscriberHandler(actor_ FsSubscriber) actor.Handler {
	return actor.NewHandler("FsSubscriber", &FsSubscriberReceiver{}, actor_)
}

// FsSubscriberContractId returns the capability contract id for this interface
func FsSubscriberContractId() string { return "jordanrash:terminal:fs" }

// FsSubscriberReceiver receives messages defined in the FsSubscriber service interface
type FsSubscriberReceiver struct{}

func (r *FsSubscriberReceiver) Dispatch(ctx *actor.Context, svc interface{}, message *actor.Message) (*actor.Message, error) {
	svc_, _ := svc.(FsSubscriber)
	switch message.Method {

	case "HandleFsMessage":
		{

			d := msgpack.NewDecoder(message.Arg)
			value, err_ := MDecodeFsMsg(&d)
			if err_ != nil {
				return nil, err_
			}

			resp, err := svc_.HandleFsMessage(ctx, value)
			if err != nil {
				return nil, err
			}

			var sizer msgpack.Sizer
			size_enc := &sizer
			resp.MEncode(size_enc)
			buf := make([]byte, sizer.Len())
			encoder := msgpack.NewEncoder(buf)
			enc := &encoder
			resp.MEncode(enc)
			return &actor.Message{Method: "FsSubscriber.HandleFsMessage", Arg: buf}, nil
		}
	case "HandleFsStatus":
		{
			resp, err := svc_.HandleFsStatus(ctx)
			if err != nil {
				return nil, err
			}

			var sizer msgpack.Sizer
			size_enc := &sizer
			resp.MEncode(size_enc)
			buf := make([]byte, sizer.Len())
			encoder := msgpack.NewEncoder(buf)
			enc := &encoder
			resp.MEncode(enc)
			return &actor.Message{Method: "FsSubscriber.HandleFsStatus", Arg: buf}, nil
		}
	default:
		return nil, actor.NewRpcError("MethodNotHandled", "FsSubscriber."+message.Method)
	}
}

// FsSubscriberSender sends messages to a FsSubscriber service
type FsSubscriberSender struct{ transport actor.Transport }

// NewActorSender constructs a client for actor-to-actor messaging
// using the recipient actor's public key
func NewActorFsSubscriberSender(actor_id string) *FsSubscriberSender {
	transport := actor.ToActor(actor_id)
	return &FsSubscriberSender{transport: transport}
}

func (s *FsSubscriberSender) HandleFsMessage(ctx *actor.Context, arg FsMsg) (*HandlerFsResponse, error) {

	var sizer msgpack.Sizer
	size_enc := &sizer
	arg.MEncode(size_enc)
	buf := make([]byte, sizer.Len())

	var encoder = msgpack.NewEncoder(buf)
	enc := &encoder
	arg.MEncode(enc)

	out_buf, _ := s.transport.Send(ctx, actor.Message{Method: "FsSubscriber.HandleFsMessage", Arg: buf})
	d := msgpack.NewDecoder(out_buf)
	resp, err_ := MDecodeHandlerFsResponse(&d)
	if err_ != nil {
		return nil, err_
	}
	return &resp, nil
}
func (s *FsSubscriberSender) HandleFsStatus(ctx *actor.Context) (*HandlerFsStatus, error) {
	buf := make([]byte, 0)
	out_buf, _ := s.transport.Send(ctx, actor.Message{Method: "FsSubscriber.HandleFsStatus", Arg: buf})
	d := msgpack.NewDecoder(out_buf)
	resp, err_ := MDecodeHandlerFsStatus(&d)
	if err_ != nil {
		return nil, err_
	}
	return &resp, nil
}

// This file is generated automatically using wasmcloud/weld-codegen 0.4.6
