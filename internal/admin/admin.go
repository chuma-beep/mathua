// Package admin holds the Admin V2 authorization model: the roles and the
// permissions they grant, as data rather than as scattered checks.
//
// It is a package of its own for one reason: every route's authorization is a
// lookup here, so "who may do this" is answerable in one file and testable
// without a server. A permission matrix spread across handlers is a matrix
// nobody can review, and the failure it produces is a route that forgot a check
// rather than a route that got the check wrong.
package admin

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/chuma-beep/mathua/internal/storage"
)

// Permission names one thing an administrator may do. They are coarse on
// purpose — a permission per button is an enterprise IAM model, and Mathua is a
// small open-source project. Each one corresponds to an area of the admin
// surface, not to an endpoint.
type Permission string

const (
	// Access is the floor: may reach the admin surface at all. Every staff role has it; every
	// staff route requires it before it checks anything finer.
	Access Permission = "admin.access"

	// ReportsRead and ReportsModerate are the moderation queue. Moderators hold both: the role
	// exists to work reports, and a report you may read but not resolve is not a workflow.
	ReportsRead     Permission = "reports.read"
	ReportsModerate Permission = "reports.moderate"

	// UsersRead is read-only visibility into accounts.
	UsersRead Permission = "users.read"

	// AdminsRead is who currently holds a staff role; AdminsManage is changing it. They are
	// separate so that "can see the team" and "can rewrite the team" are different grants.
	AdminsRead   Permission = "admins.read"
	AdminsManage Permission = "admins.manage"

	// RolesGrantOwner is the one grant that is not implied by AdminsManage. Handing out the role
	// that can hand out every other role is the escalation that matters, so it has its own
	// permission that only an owner holds.
	RolesGrantOwner Permission = "roles.grant_owner"

	// AuditRead is the trail.
	AuditRead Permission = "audit.read"

	// ContentRead is inspection; ContentManage is changing content state. Inspection is what a
	// moderator needs to judge a report; management is admin work.
	ContentRead   Permission = "content.read"
	ContentManage Permission = "content.manage"
)

// permissionsByRole is the whole model. A role absent from the map holds nothing, which is the
// safe default: adding a role without listing it grants it nothing rather than everything.
var permissionsByRole = map[storage.Role]map[Permission]bool{
	storage.RoleOwner: {
		Access: true, ReportsRead: true, ReportsModerate: true,
		UsersRead: true, AdminsRead: true, AdminsManage: true, RolesGrantOwner: true,
		AuditRead: true, ContentRead: true, ContentManage: true,
	},
	storage.RoleAdmin: {
		Access: true, ReportsRead: true, ReportsModerate: true,
		UsersRead: true, AdminsRead: true, AdminsManage: true,
		AuditRead: true, ContentRead: true, ContentManage: true,
	},
	storage.RoleModerator: {
		Access: true, ReportsRead: true, ReportsModerate: true, ContentRead: true,
	},
}

// Permitted reports whether a role holds a permission.
func Permitted(role storage.Role, p Permission) bool {
	return permissionsByRole[role][p]
}

// Permissions lists everything a role may do, in the declaration order above. The client uses
// it to render only the sections a caller can reach — as a convenience, never as the check.
func Permissions(role storage.Role) []Permission {
	order := []Permission{
		Access, ReportsRead, ReportsModerate, UsersRead,
		AdminsRead, AdminsManage, RolesGrantOwner, AuditRead, ContentRead, ContentManage,
	}
	out := make([]Permission, 0, len(order))
	for _, p := range order {
		if Permitted(role, p) {
			out = append(out, p)
		}
	}
	return out
}

// CouldEscalate reports whether changing a target from `current` to `next` would raise the
// target's privilege. It is the check behind "you cannot promote yourself": nothing else about
// a self-change is dangerous, but a self-promotion is.
func CouldEscalate(current, next storage.Role) bool {
	return next.AtLeast(current) && next != current
}

// CanAssign reports whether an actor may set an arbitrary account's role to `next`.
//
// Three rules, and each exists because of a specific escalation:
//
//  1. You cannot hand out a role above your own. Otherwise an admin mints an owner and is an
//     owner by proxy.
//  2. You cannot act on someone at least as privileged as you — unless you are an owner, who is
//     the only role that manages peers. Otherwise an admin demotes an owner.
//  3. Handing out `owner` itself needs RolesGrantOwner, which only an owner holds. This is what
//     separates "an admin can manage moderators and learners" from "an admin can crown a king".
//
// It does not decide self-changes or the last-owner invariant; the caller does those, because
// they need the actor's and target's ids and the row counts.
func CanAssign(actor, targetCurrent, next storage.Role) bool {
	if !Permitted(actor, AdminsManage) {
		return false
	}
	// Rule 1: you cannot grant a role strictly above your own. An admin who could mint an owner
	// would be an owner by proxy.
	if next != actor && next.AtLeast(actor) {
		return false
	}
	// Rule 3: handing out `owner` needs its own permission, so an admin managing moderators and
	// learners is not thereby able to crown a peer.
	if next == storage.RoleOwner && !Permitted(actor, RolesGrantOwner) {
		return false
	}
	// Rule 2: you cannot act on a peer or a superior, except as an owner. This is what stops an
	// admin demoting an owner.
	if actor != storage.RoleOwner && targetCurrent.AtLeast(actor) {
		return false
	}
	return true
}

// InvitationTTL is how long a staff invitation stays usable. Long enough to survive a weekend
// and a slow reply, short enough that a forgotten invitation is not a standing door.
const InvitationTTL = 7 * 24 * time.Hour

// NewInvitationToken returns a fresh raw token and the hash to store. The raw value is returned
// exactly once, to the caller that will put it in the response; only the hash is persisted.
//
// The construction mirrors the password-reset tokens elsewhere in the project on purpose. A
// second, subtly different scheme for the same job is how one of them ends up weaker than the
// other, and nobody notices which.
func NewInvitationToken() (raw, hash string, err error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", fmt.Errorf("generate invitation token: %w", err)
	}
	raw = hex.EncodeToString(b)
	return raw, HashInvitationToken(raw), nil
}

// HashInvitationToken is the one-way map from a raw token to its stored form. SHA-256 is
// correct here and not for passwords: the input is 256 bits of uniform randomness, so there is
// nothing to brute-force and no need for a slow KDF.
func HashInvitationToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

// AssignableRoles lists the roles an actor may assign to some account, for a role picker. Like
// Permissions it is a convenience for the client and never the enforcement.
func AssignableRoles(actor storage.Role) []storage.Role {
	out := []storage.Role{}
	for _, r := range []storage.Role{storage.RoleModerator, storage.RoleAdmin, storage.RoleOwner, storage.RoleStudent} {
		// Reuse the real rule rather than restating it, so the picker cannot drift from the check.
		if CanAssign(actor, storage.RoleStudent, r) {
			out = append(out, r)
		}
	}
	return out
}
