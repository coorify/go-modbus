package modbus

import (
	"fmt"
)

type rtuEncoder struct {
	Crc16 CRC16

	MaxSize int
	MinSize int
	ErrSize int
}

func NewRTUEncoder() ModbusEncoder {
	return &rtuEncoder{
		MaxSize: 256,
		MinSize: 4,
		ErrSize: 5,
	}
}

func (e *rtuEncoder) Endode(pdu *ModbusPDU) (stream ModbusStream, err error) {
	stream = NewStream(nil)
	stream.WriteByte(pdu.Slave)
	stream.WriteByte(pdu.OpCode)
	stream.Write(pdu.Payload)

	checksum := e.Crc16.Reset().PushBytes(stream.Bytes()...).Value()
	stream.WriteByte(byte(checksum >> 0))
	stream.WriteByte(byte(checksum >> 8))
	return
}

func (e *rtuEncoder) Decode(stream ModbusStream) (pdu *ModbusPDU, err error) {
	pdu = &ModbusPDU{}
	e.Crc16.Reset()

	pdu.Slave, err = stream.ReadByte()
	if err != nil {
		return
	}
	e.Crc16.PushBytes(pdu.Slave)

	pdu.OpCode, err = stream.ReadByte()
	if err != nil {
		return
	}
	e.Crc16.PushBytes(pdu.OpCode)

	if (pdu.OpCode & 0x80) == 0x80 {
		merr := &ModbusError{
			OpCode: pdu.OpCode & 0x7F,
		}

		merr.ErCode, err = stream.ReadByte()
		if err != nil {
			return
		}
		err = merr
		e.Crc16.PushBytes(merr.ErCode)
	} else {
		switch pdu.OpCode {
		case OPCODE_READ_COILS,
			OPCODE_READ_DISCRETE_INPUTS,
			OPCODE_READ_HOLDING_REGISTERS,
			OPCODE_READ_INPUT_REGISTERS:

			var pll byte
			pll, err = stream.ReadByte()
			if err != nil {
				return
			}
			e.Crc16.PushBytes(pll)

			pdu.Payload = make([]byte, pll)
			_, err = stream.Read(pdu.Payload)
			if err != nil {
				return
			}
			e.Crc16.PushBytes(pdu.Payload...)
		case OPCODE_WRITE_COIL,
			OPCODE_WRITE_COILS,
			OPCODE_WRITE_REGISTER,
			OPCODE_WRITE_REGISTERS:
			pdu.Payload = make([]byte, 4)
			_, err = stream.Read(pdu.Payload)
			if err != nil {
				return
			}
			e.Crc16.PushBytes(pdu.Payload...)
		default:
			err = fmt.Errorf("modbus: unknown opcode %v", pdu.OpCode)
			return
		}
	}

	checkval := byte(0)
	checksum := uint16(0)
	checkval, err = stream.ReadByte()
	if err != nil {
		return
	}
	checksum |= uint16(checkval) << 0

	checkval, err = stream.ReadByte()
	if err != nil {
		return
	}
	checksum |= uint16(checkval) << 8

	if checksum != e.Crc16.Value() {
		err = fmt.Errorf("modbus: response crc '%v' does not match expected '%v'", checksum, e.Crc16.Value())
		return
	}

	return pdu, nil
}
