package currency

import "fmt"

// по go стилю структуры данны в основном файле пакета
type Currency struct {
	Id       int
	Code     string
	FullName string
	Sign     string
}

// методы валидации на структурах, без отдельного валидатора
func (c *Currency) Validate() error {
	if len(c.Code) != 3 {
		return fmt.Errorf("currency code must be exactly 3 characters")
	}

	if c.Code == "" {
		return fmt.Errorf("currency code is required")
	}

	for _, char := range c.Code {
		if char < 'A' || (char > 'Z' && char < 'a') || char > 'z' {
			return fmt.Errorf("currency code must contain only letters")
		}
	}

	if c.FullName == "" {
		return fmt.Errorf("currency full name is required")
	}

	if len(c.FullName) < 2 || len(c.FullName) > 50 {
		return fmt.Errorf("currency full name must be between 2 and 50 characters")
	}

	if c.Sign == "" {
		return fmt.Errorf("currency sign is required")
	}

	if len(c.Sign) > 10 {
		return fmt.Errorf("currency sign must be no more than 10 characters")
	}

	return nil
}
