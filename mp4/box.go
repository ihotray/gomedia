package mp4

import (
    "bytes"
    "encoding/binary"
)

const (
    BasicBoxLen = 8
    FullBoxLen  = 12
)

var (
    MOOV BasicBox = BasicBox{Type: [4]byte{'m', 'o', 'o', 'v'}}
    TRAK BasicBox = BasicBox{Type: [4]byte{'t', 'r', 'a', 'k'}}
    MDIA BasicBox = BasicBox{Type: [4]byte{'m', 'd', 'i', 'a'}}
    MINF BasicBox = BasicBox{Type: [4]byte{'m', 'i', 'n', 'f'}}
    NMHD FullBox  = FullBox{Box: NewBasicBox([4]byte{'n', 'm', 'h', 'd'}), Version: 0}
    STBL BasicBox = BasicBox{Type: [4]byte{'s', 't', 'b', 'l'}}
    MDAT BasicBox = BasicBox{Type: [4]byte{'m', 'd', 'a', 't'}}
    AVCC BasicBox = BasicBox{Type: [4]byte{'a', 'v', 'c', 'C'}}
    HVCC BasicBox = BasicBox{Type: [4]byte{'h', 'v', 'c', 'C'}}
    ESDS FullBox  = FullBox{Box: NewBasicBox([4]byte{'e', 's', 'd', 's'}), Version: 0}
    DINF BasicBox = BasicBox{Type: [4]byte{'d', 'i', 'n', 'f'}}
    EDTS BasicBox = BasicBox{Type: [4]byte{'e', 'd', 't', 's'}}
    MVEX BasicBox = BasicBox{Type: [4]byte{'m', 'v', 'e', 'x'}}
    MOOF BasicBox = BasicBox{Type: [4]byte{'m', 'o', 'o', 'f'}}
    TRAF BasicBox = BasicBox{Type: [4]byte{'t', 'r', 'a', 'f'}}
    MFRA BasicBox = BasicBox{Type: [4]byte{'m', 'f', 'r', 'a'}}
)

type BoxEncoder interface {
    Encode(buf []byte) (int, []byte)
}

type BoxDecoder interface {
    Decode(buf []byte) (int, error)
}

type BoxSize interface {
    Size() uint64
}

// aligned(8) class Box (unsigned int(32) boxtype, optional unsigned int(8)[16] extended_type) {
//     unsigned int(32) size;
//     unsigned int(32) type = boxtype;
//     if (size==1) {
//        unsigned int(64) largesize;
//     } else if (size==0) {
//        // box extends to end of file
//     }
//     if (boxtype==‘uuid’) {
//     unsigned int(8)[16] usertype = extended_type;
//  }
// }

type BasicBox struct {
    Size     uint64
    Type     [4]byte
    UserType [16]byte
}

func NewBasicBox(boxtype [4]byte) *BasicBox {
    return &BasicBox{
        Type: boxtype,
    }
}

func (box *BasicBox) Decode(rh Reader) (int, error) {
    buf := make([]byte, 8)
    if n, err := rh.ReadAtLeast(buf); err != nil {
        return n, err
    }
    nn := 0
    boxsize := binary.BigEndian.Uint32(buf)
    copy(box.Type[:], buf[4:8])
    nn = 8
    if boxsize == 1 {
        _ = buf[nn+8]
        box.Size = binary.BigEndian.Uint64(buf[nn:])
        nn += 8
    } else {
        box.Size = uint64(boxsize)
    }
    if bytes.Equal(box.Type[:], []byte("uuid")) {
        uuid := make([]byte, 16)
        if n, err := rh.ReadAtLeast(uuid); err != nil {
            return n + nn, err
        }
        copy(box.UserType[:], uuid[:])
        nn += 16
    }

    return nn, nil
}

func (box *BasicBox) Encode() (int, []byte) {
    nn := 8
    buf := make([]byte, box.Size)
    if box.Size > 0xFFFFFFFF {
        binary.BigEndian.PutUint32(buf, 1)
        copy(buf[4:], box.Type[:])
        nn += 8
        binary.BigEndian.PutUint32(buf[8:], uint32(box.Size))
    } else {
        binary.BigEndian.PutUint32(buf, uint32(box.Size))
        copy(buf[4:], box.Type[:])
    }
    if bytes.Equal(box.Type[:], []byte("uuid")) {
        copy(buf[nn:nn+16], box.UserType[:])
    }
    return nn, buf
}

// aligned(8) class FullBox(unsigned int(32) boxtype, unsigned int(8) v, bit(24) f) extends Box(boxtype) {
//     unsigned int(8) version = v;
//     bit(24) flags = f;
// }

type FullBox struct {
    Box     *BasicBox
    Version uint8
    Flags   [3]byte
}

func NewFullBox(boxtype [4]byte, version uint8) *FullBox {
    return &FullBox{
        Box:     NewBasicBox(boxtype),
        Version: version,
    }
}

func (box *FullBox) Size() uint64 {
    if box.Box.Size > 0 {
        return box.Box.Size
    } else {
        return 12
    }
}

func (box *FullBox) Decode(rh Reader) (int, error) {
    buf := make([]byte, 4)
    if n, err := rh.ReadAtLeast(buf); err != nil {
        return n, err
    }
    box.Version = buf[0]
    copy(box.Flags[:], buf[:])
    return 4, nil
}

func (box *FullBox) Encode() (int, []byte) {
    box.Box.Size = box.Size()
    offset, buf := box.Box.Encode()
    buf[offset] = box.Version
    copy(buf[offset+1:], box.Flags[:])
    return offset + 4, buf
}
