package modbus

import (
	"fmt"
	"io"
)

const (
	OPCODE_READ_COILS             = 0x01
	OPCODE_READ_DISCRETE_INPUTS   = 0x02
	OPCODE_WRITE_COIL             = 0x05
	OPCODE_WRITE_COILS            = 0x0F
	OPCODE_READ_HOLDING_REGISTERS = 0x03
	OPCODE_READ_INPUT_REGISTERS   = 0x04
	OPCODE_WRITE_REGISTER         = 0x06
	OPCODE_WRITE_REGISTERS        = 0x10
)

const (
	ERCODE_ILLEGAL_FUNCTION         = 0x01
	ERCODE_ILLEGAL_DATA_ADDRESS     = 0x02
	ERCODE_ILLEGAL_DATA_VALUE       = 0x03
	ERCODE_SERVER_DEVICE_FAILURE    = 0x04
	ERCODE_ACKNOWLEDGE              = 0x05
	ERCODE_SERVER_DEVICE_BUSY       = 0x06
	ERCODE_MEMORY_PARITY_ERROR      = 0x08
	ERCODE_GATEWAY_PATH_UNAVAILABLE = 0x0A
	ERCODE_GATEWAY_DEVICE_FAILED    = 0x0B
)

type ModbusError struct {
	OpCode byte
	ErCode byte
}

func (e *ModbusError) Error() string {
	var name string
	switch e.ErCode {
	case ERCODE_ILLEGAL_FUNCTION:
		name = "illegal function"
	case ERCODE_ILLEGAL_DATA_ADDRESS:
		name = "illegal data address"
	case ERCODE_ILLEGAL_DATA_VALUE:
		name = "illegal data value"
	case ERCODE_SERVER_DEVICE_FAILURE:
		name = "server device failure"
	case ERCODE_ACKNOWLEDGE:
		name = "acknowledge"
	case ERCODE_SERVER_DEVICE_BUSY:
		name = "server device busy"
	case ERCODE_MEMORY_PARITY_ERROR:
		name = "memory parity error"
	case ERCODE_GATEWAY_PATH_UNAVAILABLE:
		name = "gateway path unavailable"
	case ERCODE_GATEWAY_DEVICE_FAILED:
		name = "gateway target device failed to respond"
	default:
		name = "unknown"
	}
	return fmt.Sprintf("modbus: exception '%v' (%s), function '%v'", e.ErCode, name, e.OpCode)
}

type ModbusStream interface {
	io.WriterTo
	io.ReaderFrom

	io.Reader
	io.Writer
	io.Seeker
	io.ByteReader
	io.ByteWriter

	Len() int64
	Bytes() []byte
}

type ModbusPDU struct {
	OpCode  byte
	Slave   byte
	Payload []byte
}

type ModbusEncoder interface {
	Endode(pdu *ModbusPDU) (ModbusStream, error)
	Decode(stream ModbusStream) (*ModbusPDU, error)
}

type ModbusClient interface {
	ReadCoils(address, quantity uint16, slave byte) ([]bool, error)
	ReadDiscreteInputs(address, quantity uint16, slave byte) ([]bool, error)
	WriteCoil(address uint16, value bool, slave byte) (bool, error)
	WriteCoils(address uint16, values []bool, slave byte) ([]uint16, error)

	ReadInputRegisters(address, quantity uint16, slave byte) ([]uint16, error)
	ReadHoldingRegisters(address, quantity uint16, slave byte) ([]uint16, error)
	WriteRegister(address, value uint16, slave byte) ([]uint16, error)
	WriteRegisters(address uint16, values []uint16, slave byte) ([]uint16, error)
}
