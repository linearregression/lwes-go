package lwes

import (
    "bytes"
    "encoding/binary"
    "net"
)

// http://golang.org/doc/articles/json_and_go.html
type eventAttrs map[string]interface{}

type Event struct {
    // TODO should this be a normal struct?
    name string
    attributes eventAttrs
}

// NewEvent returns an initialized Event
func NewEvent() *Event {
    return &Event{attributes: make(eventAttrs)}
}

// Name returns the name or class of an event. This is separate from an attribute
func (e *Event) Name() string {
    return e.name
}

// Iterator interface
func (e *Event) Iterator() eventAttrs {
    return e.attributes
}

// Get an attribute
func (e *Event) Get(s string) interface{} {
    return e.attributes[s]
}

func (event *Event) fromBytes(buf []byte) {
    p := bytes.NewBuffer(buf)

    var nameSize byte
    binary.Read(p, binary.BigEndian, &nameSize)

    event.name = string(p.Next(int(nameSize)))

    var attrSize uint16
    binary.Read(p, binary.BigEndian, &attrSize)

    // temporary types
    var tmpUint16 uint16
    var tmpInt16  int16
    var tmpUint32 uint32
    var tmpInt32  int32
    var tmpUint64 uint64
    var tmpInt64  int64

    for i:=0; i < int(attrSize); i++ {
        var attrNameSize byte
        var attrName string
        var attrType byte

        binary.Read(p, binary.BigEndian, &attrNameSize)
        // TODO should we camelCase attrName?
        attrName = string(p.Next(int(attrNameSize)))

        binary.Read(p, binary.BigEndian, &attrType)

        switch int(attrType) {
        case 1: // LWES_U_INT_16_TOKEN
            binary.Read(p, binary.BigEndian, &tmpUint16)
            event.attributes[attrName] = tmpUint16
        case 2: // LWES_INT_16_TOKEN
            binary.Read(p, binary.BigEndian, &tmpInt16)
            event.attributes[attrName] = tmpInt16
        case 3: // LWES_U_INT_32_TOKEN
            binary.Read(p, binary.BigEndian, &tmpUint32)
            event.attributes[attrName] = tmpUint32
        case 4: // LWES_INT_32_TOKEN
            binary.Read(p, binary.BigEndian, &tmpInt32)
            event.attributes[attrName] = tmpInt32
        case 5: // LWES_STRING_TOKEN
            binary.Read(p, binary.BigEndian, &tmpUint16)
            event.attributes[attrName] = string(p.Next(int(tmpUint16)))
        case 6: // LWES_IP_ADDR_TOKEN
            tmpIp := p.Next(4)
            // not sure if this is completely accurate
            event.attributes[attrName] = net.IPv4(tmpIp[3], tmpIp[2], tmpIp[1], tmpIp[0])
        case 7: // LWES_INT_64_TOKEN
            binary.Read(p, binary.BigEndian, &tmpInt64)
            event.attributes[attrName] = tmpInt64
        case 8: // LWES_U_INT_64_TOKEN
            binary.Read(p, binary.BigEndian, &tmpUint64)
            event.attributes[attrName] = tmpUint64
        case 9: // LWES_BOOLEAN_TOKEN
            event.attributes[attrName] = 1 == p.Next(1)[0]
        }
    }
}