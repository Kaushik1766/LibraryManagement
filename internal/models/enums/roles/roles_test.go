package roles

import "testing"

func TestUserRoles_String(t *testing.T) {
	tests := []struct {
		name string
		role UserRoles
		want string
	}{
		{
			name: "Staff role",
			role: Staff,
			want: "Staff",
		},
		{
			name: "Customer role",
			role: Customer,
			want: "Customer",
		},
		{
			name: "Invalid role",
			role: UserRoles(999),
			want: "Invalid Role",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.role.String(); got != tt.want {
				t.Errorf("UserRoles.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
