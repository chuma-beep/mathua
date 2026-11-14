# Launch plan

Status: draft, 2026-11-14.

The product is ready to show. The catalogue holds 3,015 concepts across 24 domains and 29 courses. Every concept has a lesson and three knowledge points. The hosted site at mathua.com runs on a Cloudflare Worker with a guarded Fly API origin.

The goal of this plan is not a single spike of traffic. It is a repeatable path from "someone sees Mathua" to "someone learns with Mathua and tells another person".

## Goal

After 90 days:

- 1,000 GitHub stars.
- 500 activated learners, where activation means 10 graded answers in the first week.
- 25 contributors with at least one merged change.
- A launch post that still sends referrals in month three.

Record the starting numbers on launch day rather than copying them here, because stars and signups move.

## Who it is for

- **Students preparing for exams.** They need a path, not a video library. The prerequisite graph answers "what should I do next" and the evidence model stops them moving on too early.
- **Self-learners returning to mathematics.** They need to find their frontier quickly. The diagnostic places them and the review queue keeps old material alive.
- **Teachers and tutors.** They need a read-only view of a learner's progress. Share links and transcripts cover that without exposing answers or review dates.
- **Maintainers and educators who want to self-host.** They need a real deployment story. One Go binary, SQLite or Postgres, Docker Compose and an MIT licence provide it.

## Positioning

Mathua is an open-source adaptive math learning engine with a checked prerequisite graph, generated problems and spaced review. It is not a content library with a quiz button.

Three proof points carry the message:

1. **The graph is enforced.** New material stays locked until its prerequisites are mastered. The graph is checked for cycles and orphans at build time.
2. **Mastery is evidence, not a streak.** The engine reads the last six attempts and scores accuracy, difficulty, variety and time.
3. **The catalogue is broad and deep.** 3,015 concepts, 24 domains and 9,045 knowledge points, from arithmetic through real analysis and optimization.

## Phase 1: conversion

The first job is to make the existing traffic count.

- README opens with a working demo link and a star call to action.
- The landing hero has a direct demo path for people who do not want an account.
- The signup to first graded answer path is under two minutes on a phone.
- Cloudflare Web Analytics is on for the landing page and the app.
- The Docker Compose path starts from a clean checkout and serves a lesson.
- The login monitor runs after every deploy and every 30 minutes.

Exit condition: a new visitor can reach a graded answer without reading documentation, on desktop and on a phone.

## Phase 2: launch week

One channel per day. Do not paste the same text everywhere.

- **Day 1, community.** Post in the Mathua Discord and ask existing learners what to fix before a wider launch. Fix the top two issues.
- **Day 2, Show HN.** Post the Show HN entry in the morning, then stay in the thread for six hours. Answer questions about the graph, the evidence model and the deployment.
- **Day 3, learner communities.** Post the learner-facing version to r/learnmath and a mathematics education community. Lead with the problem, not the repository.
- **Day 4, open-source communities.** Post the technical version to r/opensource and a Go community. Lead with the architecture and the self-hosting story.
- **Day 5, social.** Publish the short thread on X and the single post on LinkedIn. Use the same hero image and one concrete example.
- **Day 6, recap.** Publish what launched, what broke and what changed because of the feedback. Link the issues that came from the launch.

Rules for the week:

- Link to mathua.com and the repository, never to a signup wall.
- Answer every technical question with a concrete file or measurement.
- Do not argue with criticism. Turn it into an issue or a fix.
- Do not ask for stars in the first sentence. Earn the click with the demo.

## Phase 3: technical content

One substantial post every two weeks. Each post should teach something true about the system, not announce a feature.

Candidate posts:

1. **Why a checked prerequisite graph beats a recommendation feed.** The graph, the eligibility gate and the review exception.
2. **What mastery should cost.** The evidence model, the six-attempt window and the measurements behind the 0.6 threshold.
3. **Generating 3,015 concepts worth of problems.** The generator contract, the diversity floor and the tests that keep answers honest.
4. **One binary, two databases, one origin.** SQLite, Postgres, the Worker proxy and the edge guard.
5. **What eight hours of login downtime taught us.** The build-time environment leak, the asset routing bug and the browser-level monitor that now catches both.

Each post ends with a repository link and one specific invitation: try the demo, open an issue or self-host it.

## Phase 4: ecosystem

- Label 20 issues as good first issues with a file path and an acceptance test.
- Publish a contributor guide for a new concept, a new generator and a new lesson.
- Package three course paths for classrooms with a teacher share link.
- Add an export format that a school can import.
- Run a monthly "fix a question" session from the report queue.

## Metrics

| Metric | Target | Source |
|---|---|---|
| GitHub stars | 1,000 in 90 days | GitHub API |
| Activated learners | 500 in 90 days | signup and attempt tables |
| Week-one retention | 25 percent | attempts by student and date |
| Contributors | 25 merged | GitHub pull requests |
| Referral traffic | top three channels reported | Cloudflare Web Analytics |
| Question reports | under 5 open at any time | admin report queue |

Review the numbers weekly. Change one thing at a time.

## Risks

- **One Fly machine.** A launch spike can saturate the API while the Worker keeps serving static pages. Watch response time and scale the machine before the Show HN post.
- **No offsite backup yet.** Enable R2 and Litestream before the launch week, not after.
- **Email and push are not configured.** Do not market reminders by email or push. In-app reminders work today.
- **Content quality.** A launch brings more reports. Keep the report queue below five open items and publish the fixes.
- **Maintainer bandwidth.** A launch week is a support week. Block the calendar and keep the scope to fixes, not new features.

## Checklist

- [ ] README demo link and star line
- [ ] Landing demo path
- [ ] Cloudflare Web Analytics
- [ ] Docker Compose smoke test
- [ ] Login monitor green for 24 hours
- [ ] R2 backup and restore drill
- [ ] Workers Builds green on main
- [ ] Five screenshots and one short screen recording
- [ ] Show HN draft reviewed by two people
- [ ] Discord announcement drafted
- [ ] Support rota for launch week
- [ ] Rollback note: remove the Worker routes, Fly still serves the origin
