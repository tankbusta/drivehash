// +build windows

package drivehash

import (
	"golang.org/x/sys/windows"
)

func checkIfAdmin() bool {
	var sid *windows.SID

	// See https://docs.microsoft.com/en-us/windows/desktop/api/securitybaseapi/nf-securitybaseapi-checktokenmembership
	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		2,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
		&sid)
	if err != nil {
		return false
	}
	defer windows.FreeSid(sid)
	token := windows.Token(0)

	member, err := token.IsMember(sid)
	if err != nil {
		return false
	}

	return token.IsElevated() && member
}
