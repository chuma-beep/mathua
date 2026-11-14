# Launch copy

Status: draft, 2026-11-14.

Copy for the launch week. Keep the claims to what the product does today. Do not promise email or push reminders until their providers are configured.

## Boilerplate

**One line:** Mathua is an open-source adaptive math learning engine with a checked prerequisite graph, generated problems and spaced review.

**One paragraph:** Mathua teaches mathematics as a graph of 3,015 atomic concepts across 24 domains. New material stays locked until its prerequisites are mastered, mastery is judged on evidence from the last six attempts, and FSRS-6 brings mastered material back before it decays. It is written in Go with a Next.js frontend, runs on SQLite or Postgres and is MIT licensed.

## Show HN

**Title**

Show HN: Mathua – open-source adaptive math learning engine

**Body**

Hi HN,

Mathua is an open-source adaptive math learning engine. It covers 3,015 concepts from arithmetic through real analysis and optimization, and it refuses to show new material until the prerequisites are actually mastered.

The parts I think are interesting:

- The prerequisite graph is checked for cycles and orphans at build time. Eligibility is enforced, not suggested.
- Mastery is evidence-based: the last six attempts are scored on accuracy, difficulty, variety and time. A single miss no longer resets a 30-answer streak.
- Every problem is generated on demand by a Go generator, with answer-recomputing tests. There is no static question bank to memorise.
- Review uses FSRS-6, and decayed reviews bypass the prerequisite gate because retention is not acquisition.

It is written in Go with a Next.js frontend, runs on SQLite or Postgres and is MIT licensed. The hosted version is at mathua.com and self-hosting is one Docker Compose command.

I would value feedback on the graph and the mastery model, especially from anyone who teaches.

## Reddit, learner communities

**Title**

I built an open-source math app that will not let you move on until you have actually mastered the prerequisite

**Body**

Most math apps let you click Next. Mathua does not. It treats mathematics as a graph of 3,015 concepts, and new material stays locked until the prerequisite is mastered.

A few things that make it different:

- A diagnostic places you at your frontier instead of starting everyone at arithmetic.
- Mastery is based on evidence from your last six attempts, not one lucky streak.
- Reviews come back before you forget, using FSRS-6.
- Every problem is generated, so you cannot memorise a question bank.

It is open source and free to use at mathua.com. No account is needed to try it.

I am looking for feedback from people actually studying: where does the path feel wrong, and which explanations are unclear?

## Reddit, open-source communities

**Title**

Mathua: an open-source adaptive math engine in Go with a checked prerequisite DAG and generated problems

**Body**

I have been building Mathua, an MIT-licensed adaptive math learning engine. It is a single Go binary serving a Next.js static export, with SQLite for local use and Postgres for multi-user deployments.

The architecture choices that might interest this community:

- The concept graph is validated at build time for cycles, orphans and duplicate IDs.
- Generators are registered per concept with a contract test for three section-aligned question families, answer recomputation and a diversity floor.
- Mastery and review are separate: evidence decides acquisition, FSRS-6 decides retention.
- An edge guard requires a shared secret on API requests, so the Fly origin cannot be called directly around the Cloudflare Worker.
- Everything is embeddable and self-hostable, with Docker Compose in the repository.

Repository: https://github.com/chuma-beep/mathua

I would appreciate a review of the storage and scheduler design and any advice on making the contributor path easier.

## X thread

**Post 1**

Most math apps let you click Next. Mathua does not.

It is an open-source adaptive math engine with a checked prerequisite graph. It covers 3,015 concepts across 24 domains, and new material stays locked until the prerequisite is mastered.

mathua.com

**Post 2**

Mastery is not a 30-answer streak.

The engine reads your last six attempts and scores accuracy, difficulty, variety and time. One miss no longer erases the run.

**Post 3**

Every problem is generated, not pulled from a bank. The same concept gives you a different question every time, and the tests recompute every answer.

**Post 4**

Review uses FSRS-6. Mastered concepts come back before they decay, and reviews bypass the prerequisite gate because retention is not acquisition.

**Post 5**

It is Go, Next.js, SQLite or Postgres, MIT licensed and self-hostable.

Hosted at mathua.com. Repository: github.com/chuma-beep/mathua

## LinkedIn

I have spent the last months building Mathua, an open-source adaptive mathematics learning engine.

It covers 3,015 concepts across 24 domains, from arithmetic through real analysis and optimization. Instead of a video library with a quiz button, it enforces a checked prerequisite graph, judges mastery on evidence from the last six attempts and schedules review with FSRS-6.

The part I am most interested in is the measurement: what should mastery cost, and how do you keep a learner moving without letting them skip the foundation?

Mathua is MIT licensed, written in Go with a Next.js frontend and self-hostable with SQLite or Postgres. You can try it at mathua.com.

I would welcome feedback from teachers, tutors and maintainers.

## Discord

Mathua now has 3,015 concepts across 24 domains, a new Optimization domain with 200 concepts and in-app reminders for reviews, study days, milestones and inactivity.

The hosted site is at mathua.com. If you find a question that is wrong, use the report button. If you want to help, the repository has good first issues.

## Replies to common questions

**How is this different from Khan Academy or Math Academy?**

The graph and the evidence model. Mathua will not show new material until the prerequisite is mastered, and mastery is judged on a window of attempts rather than a single streak. It is also open source and self-hostable.

**Is this AI?**

No. Problems are generated by deterministic Go functions, and grading is rule-based, with SymPy for symbolic equivalence. There is an optional assistant, but the learning loop does not depend on a model.

**Where do the questions come from?**

Every question is generated on demand from a parameterised Go function. The generator contract requires three section-aligned families per concept, and the tests recompute each answer from the facts the generator published.

**Is it free?**

There is no paywall on the hosted site today. The code is MIT licensed, so self-hosting is free.

**Can I self-host it?**

Yes. The repository has a Docker Compose file, a single Go binary and a SQLite default. Postgres is supported for multi-user deployments.

**Why another math app?**

Because the gap is not content. It is the path: what to do next, when to review and when a learner is actually ready to move on.
