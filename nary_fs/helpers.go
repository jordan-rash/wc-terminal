package main

import msgpack "github.com/wasmcloud/tinygo-msgpack"

func DecodeFS(arg []byte) (*NaryFs, error) {
	d := msgpack.NewDecoder(arg)
	fs, err := MDecodeNaryFs(&d)
	if err != nil {
		return nil, err
	}
	return &fs, nil
}

func EncodeFS(arg NaryFs) []byte {
	var sizer msgpack.Sizer
	size_enc := &sizer
	arg.MEncode(size_enc)
	buf := make([]byte, sizer.Len())
	encoder := msgpack.NewEncoder(buf)
	enc := &encoder
	arg.MEncode(enc)
	return buf
}

func DecodeFSMsg(arg []byte) (*FsMsg, error) {
	d := msgpack.NewDecoder(arg)
	fs, err := MDecodeFsMsg(&d)
	if err != nil {
		return nil, err
	}
	return &fs, nil
}

func EncodeFSMsg(arg NaryFs) []byte {
	var sizer msgpack.Sizer
	size_enc := &sizer
	arg.MEncode(size_enc)
	buf := make([]byte, sizer.Len())
	encoder := msgpack.NewEncoder(buf)
	enc := &encoder
	arg.MEncode(enc)
	return buf
}
