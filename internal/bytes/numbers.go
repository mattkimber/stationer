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

func GetWord(value int) string {
	a := value & 0xFF
	b := value >> 8 & 0xFF
	return fmt.Sprintf("%02X %02X", a, b)
}
