package enums

// AddressType is the chain family key of an agent's registry addresses map.
type AddressType string

const (
	AddressTypeEvm AddressType = "evm"
)

func (t AddressType) String() string {
	return string(t)
}
