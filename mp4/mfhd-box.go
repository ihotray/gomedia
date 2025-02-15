package mp4

import "encoding/binary"

// aligned(8) class MovieFragmentHeaderBox extends FullBox(‘mfhd’, 0, 0){
// 	unsigned int(32) sequence_number;
// }

type MovieFragmentHeaderBox struct {
    Box            *FullBox
    SequenceNumber uint32
}

func NewMovieFragmentHeaderBox(sequence uint32) *MovieFragmentHeaderBox {
    return &MovieFragmentHeaderBox{
        Box:            NewFullBox([4]byte{'m', 'f', 'h', 'd'}, 0),
        SequenceNumber: sequence,
    }
}

func (mfhd *MovieFragmentHeaderBox) Size() uint64 {
    return mfhd.Box.Size() + 4
}

func (mfhd *MovieFragmentHeaderBox) Decode(rh Reader) (offset int, err error) {
    if offset, err = mfhd.Box.Decode(rh); err != nil {
        return
    }
    buf := make([]byte, 4)
    if _, err = rh.ReadAtLeast(buf); err != nil {
        return 0, err
    }
    mfhd.SequenceNumber = binary.BigEndian.Uint32(buf)
    return offset + 4, nil
}

func (mfhd *MovieFragmentHeaderBox) Encode() (int, []byte) {
    mfhd.Box.Box.Size = mfhd.Size()
    offset, boxdata := mfhd.Box.Encode()
    binary.BigEndian.PutUint32(boxdata[offset:], mfhd.SequenceNumber)
    return offset + 4, boxdata
}

func makeMfhdBox(frament uint32) (boxdata []byte) {
    mfhd := NewMovieFragmentHeaderBox(frament)
    _, boxdata = mfhd.Encode()
    return
}
