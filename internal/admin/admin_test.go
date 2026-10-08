package admin

import (
	"testing"

	"github.com/chuma-beep/mathua/internal/storage"
)

// The permission matrix, as executable specification. These are the properties the server
// relies on; if the model changes, this is where the change is felt first.

func TestPermittedMatrix(t *testing.T) {
	for _, tc := range []struct {
		role storage.Role
		perm Permission
		want bool
	}{
		// Owner holds everything.
		{storage.RoleOwner, Access, true},
		{storage.RoleOwner, ReportsModerate, true},
		{storage.RoleOwner, UsersRead, true},
		{storage.RoleOwner, AdminsManage, true},
		{storage.RoleOwner, RolesGrantOwner, true},
		{storage.RoleOwner, AuditRead, true},
		{storage.RoleOwner, ContentManage, true},

		// Admin runs the place but cannot crown an owner, and cannot touch the audit trail's
		// shape — it can read it, which is the point.
		{storage.RoleAdmin, ReportsModerate, true},
		{storage.RoleAdmin, AdminsManage, true},
		{storage.RoleAdmin, AuditRead, true},
		{storage.RoleAdmin, RolesGrantOwner, false},

		// Moderator works reports and inspects content. No account management, no audit.
		{storage.RoleModerator, Access, true},
		{storage.RoleModerator, ReportsRead, true},
		{storage.RoleModerator, ReportsModerate, true},
		{storage.RoleModerator, ContentRead, true},
		{storage.RoleModerator, UsersRead, false},
		{storage.RoleModerator, AdminsRead, false},
		{storage.RoleModerator, AdminsManage, false},
		{storage.RoleModerator, AuditRead, false},

		// A learner holds nothing, including not the floor.
		{storage.RoleStudent, Access, false},
		{storage.RoleStudent, ReportsRead, false},
	} {
		if got := Permitted(tc.role, tc.perm); got != tc.want {
			t.Errorf("Permitted(%s, %s) = %v, want %v", tc.role, tc.perm, got, tc.want)
		}
	}
}

func TestCanAssign(t *testing.T) {
	for _, tc := range []struct {
		name          string
		actor         storage.Role
		targetCurrent storage.Role
		next          storage.Role
		want          bool
	}{
		{"owner grants owner", storage.RoleOwner, storage.RoleStudent, storage.RoleOwner, true},
		{"owner grants admin", storage.RoleOwner, storage.RoleStudent, storage.RoleAdmin, true},
		{"owner demotes peer owner", storage.RoleOwner, storage.RoleOwner, storage.RoleStudent, true},
		{"owner demotes admin", storage.RoleOwner, storage.RoleAdmin, storage.RoleStudent, true},

		{"admin grants moderator", storage.RoleAdmin, storage.RoleStudent, storage.RoleModerator, true},
		{"admin grants admin (at own rank)", storage.RoleAdmin, storage.RoleStudent, storage.RoleAdmin, true},
		{"admin cannot grant owner", storage.RoleAdmin, storage.RoleStudent, storage.RoleOwner, false},
		{"admin cannot act on an owner", storage.RoleAdmin, storage.RoleOwner, storage.RoleStudent, false},
		{"admin cannot act on a peer admin", storage.RoleAdmin, storage.RoleAdmin, storage.RoleStudent, false},

		{"moderator has no manage permission", storage.RoleModerator, storage.RoleStudent, storage.RoleModerator, false},
		{"student has no manage permission", storage.RoleStudent, storage.RoleStudent, storage.RoleModerator, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := CanAssign(tc.actor, tc.targetCurrent, tc.next); got != tc.want {
				t.Errorf("CanAssign(%s, %s, %s) = %v, want %v",
					tc.actor, tc.targetCurrent, tc.next, got, tc.want)
			}
		})
	}
}

func TestCouldEscalate(t *testing.T) {
	if !CouldEscalate(storage.RoleStudent, storage.RoleModerator) {
		t.Error("student → moderator is an escalation")
	}
	if CouldEscalate(storage.RoleAdmin, storage.RoleAdmin) {
		t.Error("admin → admin is not a change, let alone an escalation")
	}
	if CouldEscalate(storage.RoleOwner, storage.RoleStudent) {
		t.Error("owner → student is a demotion")
	}
}

// The picker must offer exactly what the check allows: a role the client is shown but the server
// refuses is a broken control, and a role the server allows but the client hides is a feature
// nobody can reach.
func TestAssignableRolesMatchCanAssign(t *testing.T) {
	for _, actor := range []storage.Role{storage.RoleOwner, storage.RoleAdmin, storage.RoleModerator, storage.RoleStudent} {
		got := map[storage.Role]bool{}
		for _, r := range AssignableRoles(actor) {
			got[r] = true
		}
		for _, candidate := range []storage.Role{storage.RoleStudent, storage.RoleModerator, storage.RoleAdmin, storage.RoleOwner} {
			want := CanAssign(actor, storage.RoleStudent, candidate)
			if got[candidate] != want {
				t.Errorf("AssignableRoles(%s) lists %s = %v, but CanAssign says %v",
					actor, candidate, got[candidate], want)
			}
		}
	}
}

func TestInvitationTokenIsHashedAndUnique(t *testing.T) {
	raw1, hash1, err := NewInvitationToken()
	if err != nil {
		t.Fatalf("token: %v", err)
	}
	raw2, hash2, err := NewInvitationToken()
	if err != nil {
		t.Fatalf("token: %v", err)
	}
	if raw1 == raw2 || hash1 == hash2 {
		t.Fatal("two invitations produced the same token")
	}
	if hash1 == raw1 {
		t.Fatal("the stored hash equals the raw token")
	}
	if HashInvitationToken(raw1) != hash1 {
		t.Fatal("the hash is not the documented function of the raw token")
	}
}
