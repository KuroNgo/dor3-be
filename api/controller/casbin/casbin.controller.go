package casbin_controller

type RoleData struct {
	Role   string   `json:"role"`
	API    []string `json:"api"`
	Method []string `json:"method"`
}

type UserRole struct {
	UserID []string `json:"user_id"`
	Role   string   `json:"role"`
}
