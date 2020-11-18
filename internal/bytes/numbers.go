package bytes

import "fmt"

func GetDouble(value int) string {
	a := value & 0xFF
	b := value >> 8 & 0xFF
	c := value >> 16 & 0xFF
	d := value >> 24 & 0xFF
	return fmt.Sprintf("%02X %02X %02X %02X", a, b, c, d)
}

func GetByte(value int) string {
	a := value & 0xFF
	return fmt.Sprintf("%02X", a)
}

func GetVariableByte(value int) string {
	if value >= 255 {
		return fmt.Sprintf("FF %s", GetWord(value))
	}

	return GetByte(value)
}

func GetWord(value int) string {
	a := value & 0xFF
	b := value >> 8 & 0xFF
	return fmt.Sprintf("%02X %02X", a, b)
}

func GetCallbackResultByte(value int) string {
	a := value & 0xFF
	return fmt.Sprintf("%02X 80", a)
}