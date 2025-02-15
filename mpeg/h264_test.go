package mpeg

import (
    "reflect"
    "testing"
)

var sps1 []byte = []byte{0x64, 0x00, 0x28, 0xAC, 0x2C, 0xA4, 0x01, 0xE0, 0x08, 0x9F, 0x97, 0xFF, 0x00, 0x01, 0x00, 0x01, 0x52, 0x02, 0x02, 0x02, 0x80, 0x00,
    0x01, 0xF4, 0x80, 0x00, 0x75, 0x30, 0x70, 0x10, 0x00, 0x16, 0xE3, 0x60, 0x00, 0x08, 0x95, 0x45, 0xF8, 0xC7, 0x07, 0x68, 0x58, 0xB4, 0x48,
}

var sps2 []byte = []byte{0x00, 0x00, 0x00, 0x01, 0x67, 0x64, 0x00, 0x28, 0xAC, 0x2C, 0xA4, 0x01, 0xE0, 0x08, 0x9F, 0x97, 0xFF, 0x00, 0x01, 0x00, 0x01, 0x52, 0x02, 0x02, 0x02, 0x80, 0x00,
    0x01, 0xF4, 0x80, 0x00, 0x75, 0x30, 0x70, 0x10, 0x00, 0x16, 0xE3, 0x60, 0x00, 0x08, 0x95, 0x45, 0xF8, 0xC7, 0x07, 0x68, 0x58, 0xB4, 0x48,
}

func TestSPS_Decode(t *testing.T) {
    type args struct {
        bs *BitStream
    }
    tests := []struct {
        name string
        sps  *SPS
        args args
    }{
        {
            name: "sps1",
            sps:  new(SPS),
            args: args{
                bs: NewBitStream(sps1),
            }},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.sps.Decode(tt.args.bs)
            t.Log(tt.sps)
        })
    }
}

func TestGetH264Resolution(t *testing.T) {
    type args struct {
        sps []byte
    }
    tests := []struct {
        name       string
        args       args
        wantWidth  uint32
        wantHeight uint32
    }{
        {
            name: "Resolution1",
            args: args{
                sps: sps2,
            }, wantWidth: 1920, wantHeight: 1080,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            gotWidth, gotHeight := GetH264Resolution(tt.args.sps)
            t.Logf("%d,%d", gotWidth, gotHeight)
            if gotWidth != tt.wantWidth {
                t.Errorf("GetH264Resolution() gotWidth = %v, want %v", gotWidth, tt.wantWidth)
            }
            if gotHeight != tt.wantHeight {
                t.Errorf("GetH264Resolution() gotHeight = %v, want %v", gotHeight, tt.wantHeight)
            }
        })
    }
}

var spss1 [][]byte = [][]byte{{0x00, 0x00, 0x00, 0x01, 0x67, 0x64, 0x00, 0x0A, 0xAC, 0x72, 0x84, 0x44,
    0x26, 0x84, 0x00, 0x00, 0x03, 0x00, 0x04, 0x00, 0x00, 0x03, 0x00, 0xCA, 0x3C, 0x48, 0x96, 0x11, 0x80}}
var ppss1 [][]byte = [][]byte{{0x00, 0x00, 0x00, 0x01, 0x68, 0xE8, 0x43, 0x8F, 0x13, 0x21, 0x30}}
var want1 []byte = []byte{0x01, 0x64, 0x00, 0x0A, 0xFF, 0xE1, 0x00, 0x19, 0x67, 0x64, 0x00, 0x0A, 0xAC, 0x72, 0x84, 0x44,
    0x26, 0x84, 0x00, 0x00, 0x03, 0x00, 0x04, 0x00, 0x00, 0x03, 0x00, 0xCA, 0x3C, 0x48, 0x96, 0x11,
    0x80, 0x01, 0x00, 0x07, 0x68, 0xE8, 0x43, 0x8F, 0x13, 0x21, 0x30}

func TestCreateH264AVCCExtradata(t *testing.T) {
    type args struct {
        spss [][]byte
        ppss [][]byte
    }
    tests := []struct {
        name string
        args args
        want []byte
    }{
        {name: "extrandata1", args: args{
            spss: spss1,
            ppss: ppss1,
        }, want: want1},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := CreateH264AVCCExtradata(tt.args.spss, tt.args.ppss); !reflect.DeepEqual(got, tt.want) {
                t.Errorf("CreateH264AVCCExtradata() = %v, want %v", got, tt.want)
            }
        })
    }
}
