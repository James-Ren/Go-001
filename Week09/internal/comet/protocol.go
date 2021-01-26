package comet

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
)

type Proto struct {
	Ver  int32
	Op   int32
	Seq  int32
	Body []byte
}

const (
	// MaxBodySize max proto body size
	MaxBodySize = int32(1 << 12)
)

const (
	// size
	_packSize      = 4
	_headerSize    = 2
	_verSize       = 2
	_opSize        = 4
	_seqSize       = 4
	_heartSize     = 4
	_rawHeaderSize = _packSize + _headerSize + _verSize + _opSize + _seqSize
	_maxPackSize   = MaxBodySize + int32(_rawHeaderSize)
)

var (
	// ErrProtoPackLen proto packet len error
	ErrProtoPackLen = errors.New("default server codec pack length error")
	// ErrProtoHeaderLen proto header len error
	ErrProtoHeaderLen = errors.New("default server codec header length error")
)

var (
	// ProtoReady proto ready
	ProtoReady = &Proto{Op: OpProtoReady}
	// ProtoFinish proto finish
	ProtoFinish = &Proto{Op: OpProtoFinish}
)

// ReadTCP read a proto from TCP reader.
func (p *Proto) ReadTCP(rr *bufio.Reader) (err error) {
	var (
		bodyLen   int
		headerLen int16
		packLen   int32
		ver       int16
	)

	err = binary.Read(rr, binary.BigEndian, &packLen)
	if err != nil {
		return err
	}
	err = binary.Read(rr, binary.BigEndian, &headerLen)
	if err != nil {
		return err
	}
	err = binary.Read(rr, binary.BigEndian, &ver)
	if err != nil {
		return err
	}
	p.Ver = int32(ver)
	err = binary.Read(rr, binary.BigEndian, &p.Op)
	if err != nil {
		return err
	}
	err = binary.Read(rr, binary.BigEndian, &p.Seq)
	if err != nil {
		return err
	}
	if packLen > _maxPackSize {
		return ErrProtoPackLen
	}
	if headerLen != _rawHeaderSize {
		return ErrProtoHeaderLen
	}
	if bodyLen = int(packLen - int32(headerLen)); bodyLen > 0 {
		p.Body = make([]byte, bodyLen)
		_, err = io.ReadFull(rr, p.Body)
	} else {
		p.Body = nil
	}
	return
}

func (p *Proto) WriteTCP(wr *bufio.Writer) (err error) {
	var (
		packLen int
	)
	packLen = _rawHeaderSize + _heartSize
	err = binary.Write(wr, binary.BigEndian, int32(packLen))
	if err != nil {
		return err
	}
	err = binary.Write(wr, binary.BigEndian, int16(_rawHeaderSize))
	if err != nil {
		return err
	}
	err = binary.Write(wr, binary.BigEndian, int16(p.Ver))
	if err != nil {
		return err
	}
	err = binary.Write(wr, binary.BigEndian, p.Op)
	if err != nil {
		return err
	}
	err = binary.Write(wr, binary.BigEndian, p.Op)
	if err != nil {
		return err
	}
	if p.Body != nil {
		_, err = wr.Write(p.Body)
		if err != nil {
			return
		}
	}
	err = wr.Flush()
	return
}
