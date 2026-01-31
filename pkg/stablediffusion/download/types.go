package download

import "fmt"

// String returns the string representation of Arch
func (a Arch) String() string {
	switch a {
	case AMD64:
		return "amd64"
	case ARM64:
		return "arm64"
	default:
		return "unknown"
	}
}

// String returns the string representation of OS
func (o OS) String() string {
	switch o {
	case Linux:
		return "linux"
	case Darwin:
		return "darwin"
	case Windows:
		return "windows"
	default:
		return "unknown"
	}
}

// MarshalText implements encoding.TextMarshaler
func (a Arch) MarshalText() ([]byte, error) {
	return []byte(a.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler
func (a *Arch) UnmarshalText(text []byte) error {
	parsed, err := ParseArch(string(text))
	if err != nil {
		return err
	}
	*a = parsed
	return nil
}

// MarshalText implements encoding.TextMarshaler
func (o OS) MarshalText() ([]byte, error) {
	return []byte(o.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler
func (o *OS) UnmarshalText(text []byte) error {
	parsed, err := ParseOS(string(text))
	if err != nil {
		return err
	}
	*o = parsed
	return nil
}

// GoString implements fmt.GoStringer
func (a Arch) GoString() string {
	return fmt.Sprintf("Arch(%s)", a.String())
}

// GoString implements fmt.GoStringer
func (o OS) GoString() string {
	return fmt.Sprintf("OS(%s)", o.String())
}
