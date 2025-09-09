package roles

type UserRoles int

const (
	Staff UserRoles = iota
	Customer
)

func (role UserRoles) String() string {
	switch role {
	case Staff:
		return "Staff"
	case Customer:
		return "Customer"
	default:
		return "Invalid Role"
	}
}
