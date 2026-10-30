# Friends, shared graphs and friendly competition

*Status: Proposed. Date: 2026-10-29.*

## Summary

Mathua will add mutual friends, a read-only friend profile, a read-only friend graph and a weekly friends leaderboard. The feature is opt-in, requires verified email for the first rollout and is disabled by default for accounts that may be under 18.

This document records the schema and the `FriendStore` boundary. It does not implement the UI or the server routes.

## Why this shape

The request is social, so the cost is not the graph query. The cost is deciding who may see another learner's progress and what happens when that consent is withdrawn.

The existing leaderboard already has a privacy model: `named`, `anonymous` and `hidden`. It also has a stable week boundary in `internal/leaderboard/compute.go` `WeekStart`. The friends feature should reuse both, not invent a second meaning for "this week" and not create a second privacy vocabulary.

The existing `share_token` is not suitable as a friend identity. It is a secret parent/teacher report token. The account UUID is an internal key and must not appear in a public URL. Usernames are mutable: `SetUsername` exists in both stores. A profile URL keyed on a username would break on a rename. The feature therefore needs an immutable opaque profile ID separate from the username, the account UUID and the share token.

## Decisions

### Mutual requests

A friendship exists only after both accounts consent. There is no one-way follow in the first version.

### Add by username or friend code

A learner can send a request with an exact username or an eight-character friend code. There is no fuzzy search, no email lookup and no contact import.

The friend code alphabet is `23456789ABCDEFGHJKLMNPQRSTUVWXYZ`. It excludes `0`, `1`, `I` and `O`. Codes are stored uppercase and compared case-insensitively. A learner can rotate the code. Rotation invalidates the old code immediately.

Lookups are rate-limited by both the submitted username and the submitted code. The server does not reveal whether a failed lookup was an unknown account or a block.

### Real names to friends

Accepted friends see the display name and avatar even if the public leaderboard is anonymous. `friends_visibility = nobody` is the opt-out.

The add-friend form and the accept screen must state what accepting reveals before the learner confirms: display name, avatar, level, weekly XP, all-time XP, mastered count and the read-only graph. They must also state what it does not reveal: email, raw answers, attempts and review schedule.

`friends_visibility = nobody` is the opt-out, and it hides the friendship from both sides without deleting it.

### Verified email

Sending a friend request requires a verified email. Accepting does not add a second verification gate in the first rollout; the sender's verification is what limits spam. An unverified account can still accept a request from a verified friend and can still use the rest of Mathua.

### Minors

**Decision: assume minors are possible.** Mathua has no age gate today, so treating the population as adults would be an unsupported assumption.

Social features are default-off for accounts that have not attested that they are 18 or older. The friends UI must not be publicly enabled until a written minors policy is approved. This is a phase 3 gate, not a late phase 6 audit. Phase 2 backend work can merge dark behind the feature flag.

The first rollout therefore needs either an age declaration or an explicit 18+ attestation at the point a learner enables friends. An account that cannot attest stays default-off. This is a launch gate, not a code detail.

### Feature flag

The server reads `MATHUA_FRIENDS` with values `off` and `on`. The default is `off`.

When the flag is off, every `/api/friends` route returns 404, not 403. A disabled surface should not advertise that it exists. The frontend reads `friends_enabled` from the existing config payload and hides the Friends navigation item.

Phase 2 merges with the flag off. Phase 3 enables it for internal accounts. Public enablement waits for the minors policy.

## Identity model

Add two columns to `students`:

```sql
profile_id  TEXT NOT NULL DEFAULT ''
friend_code TEXT NOT NULL DEFAULT ''
```

Both get unique indexes. `profile_id` is generated once and never changes. `friend_code` is generated once and may rotate.

- `profile_id` is the stable key in friend URLs and API paths. It is `p_` plus 16 characters from the same unambiguous alphabet. It is not secret and it is not the account UUID.
- `friend_code` is the eight-character code a learner shares to be added. It is rotatable and is never used as a bookmark key.
- `username` remains the login identifier and the human-readable search key. It is mutable and is not a URL key.

Existing accounts need a backfill. A versioned Go backfill after schema migration selects rows with an empty `profile_id` or `friend_code`, generates values with `crypto/rand` and retries on a unique conflict. Account creation inserts both values in the same transaction as the student row.

## Schema

### `friendships`

```sql
CREATE TABLE IF NOT EXISTS friendships (
    student_a    TEXT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    student_b    TEXT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    status       TEXT NOT NULL CHECK (status IN ('pending', 'accepted')),
    requested_by TEXT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    created_at   INTEGER NOT NULL,
    responded_at INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (student_a, student_b),
    CHECK (student_a < student_b),
    CHECK (requested_by = student_a OR requested_by = student_b)
);

CREATE INDEX IF NOT EXISTS idx_friendships_a_status
    ON friendships(student_a, status);
CREATE INDEX IF NOT EXISTS idx_friendships_b_status
    ON friendships(student_b, status);
```

`student_a` is always the lexicographically smaller account ID. That canonical pair makes a duplicate request in either direction impossible.

The table holds both pending requests and accepted friendships. A pending row records the direction in `requested_by`. An accepted row keeps the original requester for history.

### `friend_blocks`

```sql
CREATE TABLE IF NOT EXISTS friend_blocks (
    blocker_id TEXT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    blocked_id TEXT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    created_at INTEGER NOT NULL,
    PRIMARY KEY (blocker_id, blocked_id),
    CHECK (blocker_id <> blocked_id)
);

CREATE INDEX IF NOT EXISTS idx_friend_blocks_blocked
    ON friend_blocks(blocked_id);
```

A block is directional in storage and symmetric in effect. Every friend read checks both directions.

### Settings

The existing `students.settings` JSON blob gains one key:

```json
{
  "friends_visibility": "friends"
}
```

Allowed values are `friends` and `nobody`. Unknown values fall back to `friends`, matching the tolerant behaviour of `leaderboard_privacy`.

`friends_visibility = nobody` means:

- No new requests can be sent or accepted.
- Existing friendship rows remain.
- The friendship is hidden from both lists.
- Profile and graph reads return the same not-found response as a missing account.
- Switching back to `friends` restores the same friendships.

The row is not deleted because the privacy choice must be reversible. The hidden state is applied at read time.

## `FriendStore` interface

The new surface is a separate interface. Adding methods to `Repository` would force every existing test double to grow methods it does not use.

```go
type FriendshipStatus string

const (
    FriendshipPending  FriendshipStatus = "pending"
    FriendshipAccepted FriendshipStatus = "accepted"
)

type Friendship struct {
    StudentA    string
    StudentB    string
    Status      FriendshipStatus
    RequestedBy string
    CreatedAt   time.Time
    RespondedAt time.Time
}

type FriendRequest struct {
    ProfileID string
    Name      string
    Username  string
    Direction string // "incoming" or "outgoing"
    CreatedAt time.Time
}

type FriendSummary struct {
    ProfileID      string
    Name           string
    Username       string
    AvatarURL      string
    AvatarSettings string
    XPTotal        int
    WeeklyXP       int
    Mastered       int
    FriendsSince   time.Time
}

type FriendProfile struct {
    FriendSummary
    Visibility string
}

type FriendGraphNode struct {
    ConceptID  string
    Status     string
    MasteredAt *time.Time
}

type FriendLeaderboardRow struct {
    ProfileID string
    Name      string
    Username  string
    AvatarURL string
    WeeklyXP  int
    XPTotal   int
    Mastered  int
}

type FriendStore interface {
    ResolveFriendTarget(query string) (profileID string, err error)
    SendFriendRequest(studentID, targetProfileID string, now time.Time) (Friendship, error)
    AcceptFriendRequest(studentID, otherProfileID string, now time.Time) (Friendship, error)
    DeclineFriendRequest(studentID, otherProfileID string, now time.Time) error
    CancelFriendRequest(studentID, otherProfileID string, now time.Time) error
    RemoveFriend(studentID, otherProfileID string, now time.Time) error
    BlockUser(studentID, otherProfileID string, now time.Time) error
    UnblockUser(studentID, otherProfileID string) error
    ListFriends(studentID string) ([]FriendSummary, error)
    ListFriendRequests(studentID string) (incoming, outgoing []FriendRequest, err error)
    FriendProfile(studentID, otherProfileID string) (*FriendProfile, error)
    FriendGraph(studentID, otherProfileID string) ([]FriendGraphNode, error)
    FriendsLeaderboard(studentID, view string, week time.Time) ([]FriendLeaderboardRow, error)
    RotateFriendCode(studentID string, now time.Time) (string, error)
}
```

The storage layer returns raw avatar fields. The server applies the same avatar resolver as the leaderboard. The storage layer does not import `internal/levels`; the server computes the level from `Mastered`.

Errors:

```go
var (
    ErrFriendNotFound      = errors.New("friend target not found")
    ErrFriendBlocked       = errors.New("friend target not found")
    ErrFriendNotVerified   = errors.New("verified email required")
    ErrFriendSelf          = errors.New("cannot add yourself")
    ErrFriendVisibilityOff = errors.New("friends are disabled for this account")
)
```

`ResolveFriendTarget` and every read path return `ErrFriendNotFound` for an unknown account and `ErrFriendBlocked` for a block. The HTTP layer maps both to the same 404 body.

## Atomic block

Blocking is one transaction:

1. Delete the canonical `friendships` row for the pair. This removes a pending request and an accepted friendship.
2. Insert the `friend_blocks` row with `ON CONFLICT DO NOTHING`.
3. Commit.

The endpoints still check `friend_blocks` before every profile, graph and request read. The transaction is the source of truth; the endpoint check is defence in depth.

`UnblockUser` deletes the block. It does not recreate the friendship.

## Simultaneous requests

Two users may send to each other at the same moment. The result must be one accepted friendship, not a unique constraint error.

`SendFriendRequest` uses the canonical pair and a transactional upsert:

```sql
INSERT INTO friendships
    (student_a, student_b, status, requested_by, created_at, responded_at)
VALUES (?, ?, 'pending', ?, ?, 0)
ON CONFLICT (student_a, student_b) DO UPDATE SET
    status = CASE
        WHEN friendships.status = 'pending'
         AND friendships.requested_by <> excluded.requested_by
        THEN 'accepted'
        ELSE friendships.status
    END,
    responded_at = CASE
        WHEN friendships.status = 'pending'
         AND friendships.requested_by <> excluded.requested_by
        THEN excluded.created_at
        ELSE friendships.responded_at
    END;
```

SQLite uses `BEGIN IMMEDIATE`; PostgreSQL locks the canonical row through the insert conflict. The storage test starts both directions concurrently and requires exactly one row with `status = accepted`.

If one request arrives first and the second arrives after the first transaction commits, the second statement sees the pending row and promotes it to accepted. If both arrive before either commits, one insert wins and the other takes the conflict path. Neither case returns a constraint error.

## Visibility and access rules

`ListFriends` and `FriendsLeaderboard` join both students' settings and exclude a friendship when either side has `friends_visibility = nobody`. The row stays in `friendships`.

`FriendProfile` and `FriendGraph` check, in order:

1. The caller is authenticated.
2. The target profile ID exists.
3. The caller and target have an accepted friendship.
4. Neither side blocks the other.
5. Neither side has `friends_visibility = nobody`.

A failure at any of those steps returns the same 404. Verification is checked when a request is sent, not when a friendship is read.

`FriendGraph` returns only `concept_id`, stored `status` and `mastered_at`. It does not return attempts, answers, weakness scores, review due dates or SM-2 fields.

## Friends leaderboard

The friends board uses the same week boundary as the public board:

```go
week := leaderboard.WeekStart(now)
```

It uses the same `weekly` and `all_time` view names. It does not add a second timezone or a separate Monday calculation.

The query filters to accepted friends of the caller. It excludes either side with `friends_visibility = nobody`. It excludes any learner with `leaderboard_privacy = hidden`. Anonymous learners remain visible to friends under their display name because the friendship is mutual; the public board still shows the alias.

The response adds `Level` from `levels.Compute(Mastered)` in the server, not in storage.

## API surface

All routes are authenticated and behind `MATHUA_FRIENDS`.

```text
GET    /api/friends
POST   /api/friends/requests                 { "query": "username or code" }
POST   /api/friends/requests/{profile_id}/accept
POST   /api/friends/requests/{profile_id}/decline
POST   /api/friends/requests/{profile_id}/cancel
DELETE /api/friends/{profile_id}
POST   /api/friends/{profile_id}/block
DELETE /api/friends/{profile_id}/block
GET    /api/friends/{profile_id}/profile
GET    /api/friends/{profile_id}/graph
GET    /api/friends/leaderboard?view=weekly
POST   /api/friends/code/rotate
```

`GET /api/friends` returns:

```json
{
  "visibility": "friends",
  "friends": [],
  "incoming": [],
  "outgoing": [],
  "friend_code": "7KQ2M9XA"
}
```

The profile and graph handlers resolve `profile_id` to an account ID internally. The account UUID never appears in a response body.

## Profile URLs

Friend profile pages use a query parameter because the frontend is a static export:

```text
/friends/view?id=<profile_id>
```

A dynamic `/friends/[username]` route cannot be statically generated for arbitrary accounts. The query parameter also survives a username change because `profile_id` is immutable.

The friend code is not a URL key. It rotates, so a bookmark made from a code would break on rotation.

## Migrations

SQLite and PostgreSQL receive the same schema.

1. Add `profile_id` and `friend_code` to `students` with a temporary empty default.
2. Backfill both values for existing accounts in a versioned Go migration.
3. Add the unique indexes only after the backfill, because multiple empty defaults would collide.
4. Create `friendships` and `friend_blocks`.
5. Insert the columns and tables into both `schema` and `pgSchema`.
6. Add the tables to `DeleteAccount` in both stores, child first.
7. Do not delete friendships in `ResetProgress`. A progress reset is not an account deletion.
8. Do not touch `share_token`.

Account creation paths must write both identifiers in the same transaction as the student row. The creation paths are `CreateStudent`, `ClaimGuestStudent`, `CreateUser`, `AdoptGuestAsUser`, `CreateGoogleUser` and `CreateOAuthUser`.

## Storage tests

| Test | Requirement |
| --- | --- |
| Canonical pair | A request from A to B and B to A produces one row |
| Reciprocal race | Two concurrent requests resolve to one accepted row and no error |
| Accept | A pending incoming request becomes accepted for both sides |
| Decline | A declined row disappears and a new request can be sent later |
| Cancel | The requester can cancel and the other side sees no request |
| Unfriend | Removing an accepted row removes it from both lists |
| Block cleanup | Blocking deletes an accepted friendship or pending request in the same transaction |
| Block defence | A blocked profile, graph and request read all return not-found |
| Not-found equivalence | Unknown code, unknown username and blocked-by-target return the same error |
| Visibility nobody | The friendship is hidden from both lists and both profile reads, then restored when visibility returns to friends |
| Friend code rotation | The old code fails immediately, the new code works and the alphabet has no ambiguous characters |
| Code lookup rate limit | Both code and username lookups are limited independently |
| Account deletion | Friendships and blocks cascade; the other account remains |
| Reset progress | Friendships remain after a progress reset |
| Leaderboard week | Friends board uses `leaderboard.WeekStart` and matches the public board's weekly numbers |
| Privacy | Hidden learners are excluded from the friends board |
| Graph payload | Graph rows contain no attempts, answers, weakness or due dates |
| SQLite and PostgreSQL | Every test above runs on both stores |

## Server tests

- Unauthenticated access returns 401.
- Unverified accounts cannot send a request.
- Self-add returns 400.
- Unknown and blocked targets return the same 404.
- A non-friend cannot read a profile or graph.
- `friends_visibility = nobody` revokes access immediately.
- The feature flag returns 404 when off.
- The leaderboard week matches the public board.
- A concurrent reciprocal request returns one accepted friendship.

## Rollout

1. **RFC** — this document.
2. **Schema and storage** — migrations, `FriendStore`, storage tests, flag off.
3. **Server routes** — request, accept, decline, cancel, block, profile and graph, flag off.
4. **Friends page** — add form, requests, list and sidebar badge.
5. **Friend profile and graph** — read-only surfaces.
6. **Friends leaderboard** — weekly and all-time views.
7. **Safety review and enablement** — minors policy, rate limits, block audit and public flag.

Phase 2 can merge dark. Phase 3 cannot be publicly enabled until the minors policy is written. Phase 6 cannot complete without phase 3.

## Risks

- **Harassment and spam.** Mitigations are verified email, exact-match lookup, rate limits, block and no chat.
- **Privacy leakage through the graph.** Mitigations are status-only payloads, a friend check on every read and no raw attempts.
- **Broken bookmarks after a username change.** Mitigated by the immutable `profile_id`.
- **Friend code sharing.** Rotation invalidates the old code immediately. A code is an add mechanism, not a login credential.
- **Block state leaking through errors.** Unknown and blocked targets share one response.
- **Minors.** The social surface is default-off for accounts that cannot attest to 18+. Public rollout waits for the written policy.

## Open implementation questions

1. Should the friends page show the leaderboard immediately, or only after the learner has at least one friend?
2. Should a pending request expire after 30 days? The first version can leave it until declined or cancelled.
3. Should the friend code be visible on the profile page, or only on the Friends page?
4. Should the friend graph include `DECAYING` as a computed state, or only the stored `MASTERED` status?
5. Should the first rollout allow friends to see the all-time XP, or only weekly XP?

These are UI details. None changes the schema or the `FriendStore` boundary.
