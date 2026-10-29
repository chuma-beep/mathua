# FIRe activation and rollback

FIRe has three modes: `off`, `shadow` and `active`. Production starts in `shadow`.
Active mode applies banked implicit repetitions to eligible review schedules for
an explicit cohort. It does not grant XP and does not change mastery evidence.

## Activation requires two keys

Active mode needs both a process configuration and a durable activation row. A
client request cannot activate FIRe. A learner request cannot activate FIRe.

The process configuration comes from these environment variables:

- `MATHUA_FIRE_MODE=active`
- `MATHUA_FIRE_ACK=I-ACKNOWLEDGE-FIRE-APPLIES-SCHEDULE-CHANGES`
- `MATHUA_FIRE_COHORT=<stable cohort name>`
- `MATHUA_FIRE_ALLOWLIST=<comma-separated student IDs>`

The durable row comes from `mathua fire activate` or the owner-only admin
endpoint. Startup checks that the row is active and that its cohort equals the
environment cohort. A missing row, a shadow row or a cohort mismatch degrades
active to shadow and logs the reason.

The acknowledgement is an exact string. An empty cohort, an empty allowlist or a
wrong acknowledgement degrades to shadow.

## Operator commands

Run these commands with shell access to the deployment:

```text
mathua fire status
mathua fire report
mathua fire activate --ack I-ACKNOWLEDGE-FIRE-APPLIES-SCHEDULE-CHANGES --cohort pilot --allowlist <id>,<id> --actor <email>
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
{"ack":"I-ACKNOWLEDGE-FIRE-APPLIES-SCHEDULE-CHANGES","cohort":"pilot","allowlist":"<id>,<id>"}
```

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
queue for the selected cohort. Keep the cohort small until the report and the
conflict rate are stable.
