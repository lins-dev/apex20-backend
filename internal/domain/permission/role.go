package permission

import apexv1 "github.com/apex20/contracts/proto/apex20/v1"

// Role é um alias do tipo gerado via Protobuf.
// O apex20-contracts é a fonte única de verdade para roles em todos os serviços.
type Role = apexv1.Role

const (
	RoleGM      = apexv1.Role_ROLE_GM
	RolePlayer  = apexv1.Role_ROLE_PLAYER
	RoleTrusted = apexv1.Role_ROLE_TRUSTED
)
