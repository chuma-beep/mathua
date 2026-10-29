# FIRe activation and rollback

FIRe has three modes: `off`, `shadow` and `active`. Production starts in `shadow`.
Active mode applies banked implicit repetitions to eligible review schedules for
a configured cohort and scope. It does not grant XP and does not change mastery
evidence.

## Activation requires two keys

Active mode needs both a process configuration and a durable activation row. A
client request cannot activate FIRe. A learner request cannot activate FIRe.

The process configuration comes from these environment variables:

- `MATHUA_FIRE_MODE=active`
- `MATHUA_FIRE_ACK="I acknowledge that FIRe may apply schedule changes"`
- `MATHUA_FIRE_COHORT=<stable cohort name>`
- `MATHUA_FIRE_ALLOWLIST=<comma-separated student IDs>` or `all`

Set `MATHUA_FIRE_ALLOWLIST=all` to apply active FIRe to every learner, including
accounts created after activation. `all` must be the only token. Mixing `all`
with student IDs is ambiguous and degrades to shadow. The `*` token has the same
meaning.

The durable row comes from `mathua fire activate` or the owner-only admin
endpoint. Startup checks that the row is active and that its cohort equals the
environment cohort. A missing row, a shadow row or a cohort mismatch degrades
active to shadow and logs the reason.

The acknowledgement is an exact string. An empty cohort, an empty scope or a
wrong acknowledgement degrades to shadow.

## Operator commands

Run these commands with shell access to the deployment:

```text
mathua fire status
mathua fire report
mathua fire activate --ack "I acknowledge that FIRe may apply schedule changes" --cohort pilot --allowlist <id>,<id> --actor <email>
mathua fire activate --ack "I acknowledge that FIRe may apply schedule changes" --cohort all-users --allowlist all --actor <email>
mathua fire deactivate --actor <email>
```

`mathua fire report` is read-only. It reads banked ordinals and calls
`PreviewCredit` only. It does not call `ApplyCredit` or `ProjectCredit`. Use a
copied database or a read replica for a production-shaped report. Do not run a
mutating replay on production.

The report includes the build revision, a checksum, each target and epoch, the
ordinal, the eligibility reason, the actual and proposed due dates, the due-date
delta, stability and difficulty before and after, conflicts and the number of
reviews deferred.

## Admin endpoints

Owner and admin accounts can read these endpoints with `content.read`:

- `GET /api/admin/fire`
- `GET /api/admin/fire/report`

Only an owner can call these endpoints with `content.manage`:

- `POST /api/admin/fire/activate`
- `POST /api/admin/fire/deactivate`

The activation request body is:

```json
{"ack":"I acknowledge that FIRe may apply schedule changes","cohort":"pilot","allowlist":"<id>,<id>"}
```

Use `"allowlist":"all"` for every learner. The environment must carry the same
scope, and the process must restart.

The durable row does not change the running process. Set the environment values
and restart the process to apply the mode. This is deliberate: a row written by
mistake cannot move schedules without a second operator action.

## What active mode writes

An application transaction verifies the mode cohort, the credit epoch, the
banked ordinal, the target state, the target due date and the source attempt. It
then updates the target due date and FSRS-6 schedule, consumes exactly one
ordinal and writes one application row. The application row has a unique key.
A replay returns the stored result. A failure before commit leaves the ordinal
available for a safe retry.

An implicit review is a scheduling event. It is not evidence that the learner
answered correctly. It does not change XP, attempts, streak, mastery status,
weakness, placement state or the source concept.

## Scheduler read paths

Four read paths call `SatisfyImplicitReviews` for a learner in the configured scope:

- `NextQuestion`, when a study or practice question is served.
- `NextReviewQuestion`, when a spaced-review question is served.
- `ReviewsDue`, when the profile or dashboard reads the due count.
- `conceptSnapshots`, when the recommender builds its view.

Each call applies at most one ordinal per target and at most 32 ordinals in
total. A repeated read cannot apply the same ordinal twice: the ordinal moves
from `banked` to `projected`, the application key is unique and a retry reads
the stored row. A storage failure is logged and the ordinary learning response
continues. The next read can retry the same banked ordinal.

The 32-application bound is per call, not per lifetime. A large bank drains over
successive reads. The prefilter skips future and ineligible targets before a
transaction is opened, so a backlog of future targets cannot starve a due one.

## Explicit review precedence

A later explicit review is real learner evidence and takes precedence over an
earlier implicit application. If an explicit review commits first, the
application's due-baseline check sees the moved date and returns `stale_baseline`
without writing. If an application commits first, the explicit review runs from
the new schedule and can shorten or extend it. In both cases the application row
remains the audit record of what FIRe changed.

A retry cannot apply one ordinal twice, and an application cannot overwrite a
newer explicit review because the progress update is conditional on the exact
due date it read.

## Reset, deletion and activation races

On PostgreSQL every credit writer takes a per-student transaction advisory lock
before touching a credit table. `ResetProgress` and `DeleteAccount` take the same
lock. On SQLite the write transaction serialises them. The result is that a
reset either runs before the application, so the application finds no target and
writes nothing, or after it, so the reset removes the application with the rest
of the learner's credit. A reset cannot leave an application row behind.

The row update increments the activation revision and returns it with
`RETURNING`. Two concurrent operator changes cannot both write the same revision
number.

## Metrics

`GET /api/admin/fire` returns process counters:

- `allocations`: new credit allocations.
- `applications`: ordinals applied to schedules.
- `conflicts`: ordinals rejected by an eligibility or baseline check.
- `replays`: an application key that already existed.
- `retries`: outbox effects reprocessed after a retryable failure.
- `failures`: allocation or application errors.
- `due_shifts`: applications that moved a due date.
- `latency_micros` and `latency_count`: total application latency.

Counters are process-local. A restart resets them. Durable application rows are
the audit record.

## Alerts and rollback thresholds

Use these thresholds during a pilot:

- Stop the pilot if `failures` rises by more than five in ten minutes.
- Stop the pilot if `conflicts` is more than ten percent of `applications` plus
  `conflicts` over one hour.
- Stop the pilot if `due_shifts` is zero while `applications` is positive, or if
  `applications` is zero while eligible banked ordinals exist.
- Stop the pilot if any learner reports a due date earlier than its previous
  value, or if an implicit review changes XP, mastery, attempts or placement.
- Stop the pilot if `latency_micros / latency_count` exceeds 250 ms for ten
  minutes.

The log line `fire: applied ordinal` records each due shift. The log lines
`fire: conflict on ordinal` and `fire: deferred ordinal` record the eligibility
reason. Alert on conflicts and failures.

## Disable and rollback

Set `MATHUA_FIRE_MODE=shadow` and restart, or set the durable row to shadow with
`mathua fire deactivate`. The durable row stops new applications immediately,
even before the process restarts. Restarting with shadow or off also downgrades
the running process so it stops trying. Previously applied rows remain auditable
and are not reversed.

FIRe has no automatic reversal. A reversal must compare the application row, the
current progress row and any later learner activity. Do not overwrite a due date
that a later explicit review moved. Reversal needs a separate reconciliation
plan and an explicit operator decision.

## Limits

Active FIRe does not prove that an implicit repetition improves recall. Shadow
projection stores counterfactual due dates. Active mode changes the real review
queue for the selected cohort. An explicit allowlist is the narrowest scope.
`all` applies to every current and future learner, so use it only after the
report and the conflict rate are stable and the rollback path is ready.
