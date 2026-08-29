# WASC Accreditation & UC A-G Approval — Research for Mathua

> Research task (2026-08-29). No code changes. Primary sources only: acswasc.org, wscuc.org, ucop.edu, admission.universityofcalifornia.edu, cde.ca.gov / leginfo (statute), ftc.gov, eCFR, W3C. Mathua facts cited to repo files (`internal/diagnostic/report.go:11`, `data/courses.json:1`, `improve.md:148`).
>
> Note: `https://www.ucop.edu/ag-guide/` and `https://hs-articulation.ucal.edu/` no longer resolve; the canonical UC pages are the UC Admissions site (`admission.universityofcalifornia.edu`) and the A-G Policy Resource Guide (`hs-articulation.ucop.edu/guide`), used below. `cde.ca.gov` was unreachable from this environment (timeouts); the governing statute is cited directly from California Legislative Information instead.

---

## Section 1 — Which accreditor applies, and the WASC process

### 1.1 Two different commissions, two different missions

**WASC ACS (Accrediting Commission for Schools)** accredits pre-K-12 and adult schools — and, critically for Mathua, **Supplementary Education Programs (SEPs)**, which WASC defines as "a non-degree, non-diploma granting organization that offers programs or instruction or education services designed for elementary, secondary age or adult students seeking educational development in one or more areas" ([https://acswasc.org/supplementary-education-programs/](https://acswasc.org/supplementary-education-programs/)). WASC ACS also explicitly accredits **online schools and programs** and has aligned its criteria with the iNACOL standards "that are also used by the University of California in the evaluation of online a-g courses" ([https://acswasc.org/about/frequently-asked-questions/](https://acswasc.org/about/frequently-asked-questions/)).

**WSCUC (WASC Senior College and University Commission)** accredits senior (bachelor's and above) degree-granting colleges and universities — "Postsecondary institutions may apply if they offer **degree programs** and have a significant focus on higher education" ([https://www.wscuc.org/resources/becoming-accredited/](https://www.wscuc.org/resources/becoming-accredited/); the 2023 Handbook states WSCUC was formed "for the purposes of accrediting senior (bachelor's degrees and above) colleges and universities" ([https://www.wscuc.org/handbook2023/](https://www.wscuc.org/handbook2023/)).

**Answer: WASC ACS applies to Mathua**, and the SEP track is the natural fit for a course platform that does not itself grant diplomas or degrees. WSCUC only becomes relevant if Mathua ever becomes a degree-granting postsecondary institution. UC itself hard-wires this distinction into its a-g policy: an online course publisher that employs its own instructors "must earn accreditation as a supplementary education center/program from the Western Association of Schools and Colleges (WASC)" — i.e., WASC ACS, not WSCUC ([https://hs-articulation.ucop.edu/guide/register-your-institution/online-course-publishers/](https://hs-articulation.ucop.edu/guide/register-your-institution/online-course-publishers/)). WASC ACS also confirms it collaborates "with the University of California regarding the a-g course requirements" ([https://acswasc.org/about/acs-wasc-overview/](https://acswasc.org/about/acs-wasc-overview/)).

### 1.2 WASC ACS process, statuses, timeline, cost

**Steps** (synthesized from the official process pages):

1. **Review conditions of eligibility.** SEPs review the "SEP Conditions of Eligibility" and submit the SEP application ([https://acswasc.org/beginning-the-accreditation-process/](https://acswasc.org/beginning-the-accreditation-process/)). Eligibility requires at minimum a history of operation: "Initial visits to newly established schools can not occur prior to the second semester of operation" (same page).
2. **Application** with a nonrefundable $160 application fee (California schools) ([https://acswasc.org/beginning-the-accreditation-process/](https://acswasc.org/beginning-the-accreditation-process/)); SEP-specific: "Application Fee: $160, Nonrefundable Initial Visit Fee: $850" plus annual membership fees ([https://acswasc.org/supplementary-education-programs/](https://acswasc.org/supplementary-education-programs/); fee schedules at [https://acswasc.org/about/acs-wasc-fees/](https://acswasc.org/about/acs-wasc-fees/)).
3. **Initial visit** by a two-member team for one or two days; the team produces a report with recommendations to the Commission ([https://acswasc.org/initial-visit-process-and-procedures/](https://acswasc.org/initial-visit-process-and-procedures/)).
4. **Commission action** — one of three statuses:
   - **Initial Accreditation** (meets criteria for full accreditation), or
   - **Candidacy** (progressing toward accreditation; "Schools are normally expected to apply for full accreditation by the third year of candidacy"; a written progress report is due the spring after the initial visit), or
   - **Accreditation Status Withheld** (may reapply after deficiencies are remedied).
   Both initial accreditation and candidacy are granted "for a period of time not to exceed three years" ([https://acswasc.org/initial-accreditation-status-options/](https://acswasc.org/initial-accreditation-status-options/); [https://acswasc.org/beginning-the-accreditation-process/](https://acswasc.org/beginning-the-accreditation-process/)).
5. **First full self-study by the end of the third year**, followed by the **self-study visit** (peer visiting committee validates the self-study against WASC criteria), then a **follow-up cycle** of annual progress reports and action-plan assessment ([https://acswasc.org/about/acs-wasc-overview/](https://acswasc.org/about/acs-wasc-overview/)). Full accreditation is a six-year status with a review; status runs July 1 – June 30 ([https://acswasc.org/about/frequently-asked-questions/](https://acswasc.org/about/frequently-asked-questions/)).

**Rough timeline:** application → initial visit (a semester after opening at the earliest) → candidacy/initial accreditation (≤3 years) → first full self-study by year 3 → six-year accreditation cycle with annual reports. **Rough cost:** $160 application + $850 initial visit (SEP rates) + annual membership fees (2026-27 schedule not yet published as of fetch date; prior schedules via [https://acswasc.org/about/acs-wasc-fees/](https://acswasc.org/about/acs-wasc-fees/)).

**Evaluation criteria** — nine standards in three categories (revised 2025–26): Category A (Vision/Leadership/Faculty/Resources/Continuous Improvement/Accountability & Compliance, incl. A6 "compliance with all applicable laws and regulations"), Category B (Curriculum, Teaching & Learning, Assessment & Data Analysis), Category C (School Culture, Systems of Support, Student Success & Community Partnerships) ([https://acswasc.org/about/acs-wasc-standards/](https://acswasc.org/about/acs-wasc-standards/)). The self-study is judged against "schoolwide learner outcomes" and student achievement data — relevant for Mathua: per-concept mastery data and diagnostic reports are exactly the kind of evidence the process asks for ([https://acswasc.org/about/frequently-asked-questions/](https://acswasc.org/about/frequently-asked-questions/)).

**Non-negotiable structural prerequisites:** a minimum of 6 full-time students to be eligible for affiliation; qualified staff (private-school teachers "encouraged to be highly qualified" with subject expertise; credentialed staff for public schools); a plan for permanent retention of transcripts/records if the organization closes (records retention is called "crucial at the high school level") ([https://acswasc.org/about/frequently-asked-questions/](https://acswasc.org/about/frequently-asked-questions/)).

### 1.3 WSCUC (higher ed) — only if degree-granting

For completeness: WSCUC's path is **Inquiry → Alignment & Financial Review (eligibility decision within 60 business days) → Application for Accreditation**, then a multi-year candidacy/initial-accreditation self-study ([https://www.wscuc.org/resources/becoming-accredited/](https://www.wscuc.org/resources/becoming-accredited/)). Candidacy is limited to five years; initial accreditation is granted for 6, 8, or 10 years; reaffirmation "typically... approximately two years" ([https://www.wscuc.org/handbook2023/](https://www.wscuc.org/handbook2023/)). WSCUC's Standards focus on institutional mission/integrity, student outcomes, resources, and quality assurance ([https://www.wscuc.org/handbook2023/](https://www.wscuc.org/handbook2023/)). Accreditation by WSCUC is also the prerequisite for federal student-aid eligibility for institutions ([https://www.wscuc.org/about/about-accreditation/](https://www.wscuc.org/about/about-accreditation/)). **Not applicable to a non-degree course platform; noted only to rule it out.**

---

## Section 2 — UC a-g: Category C (Mathematics) and the course approval pathway

### 2.1 The 15-course subject requirement and Category C

UC first-year admission requires **15 yearlong high school courses with a letter grade of C or better**, at least 11 completed before the last year of high school: a) History 2 years, b) English 4, c) **Mathematics 3 (4 recommended)**, d) Science 2, e) Language other than English 2, f) Visual & performing arts 1, g) College-preparatory elective 1 ([https://admission.universityofcalifornia.edu/admission-requirements/freshman-requirements/](https://admission.universityofcalifornia.edu/admission-requirements/freshman-requirements/)).

**Category C (Mathematics)** — UC-approved high school courses must include "the topics covered in elementary and advanced algebra and two- and three-dimensional geometry; a fourth year of math is strongly recommended. A geometry course or an integrated math course with a sufficient amount of geometry content must be completed." Math taken in 7th/8th grade counts if the high school accepts it as equivalent; and "only the math and language other than English requirements may be met with coursework completed in 7th and/or 8th grade" ([https://admission.universityofcalifornia.edu/admission-requirements/freshman-requirements/subject-requirement-a-g.html](https://admission.universityofcalifornia.edu/admission-requirements/freshman-requirements/subject-requirement-a-g.html); [https://admission.universityofcalifornia.edu/counselors/preparing-freshman-students/freshman-requirements.html](https://admission.universityofcalifornia.edu/counselors/preparing-freshman-students/freshman-requirements.html)).

**Official Course Criteria & Guidance for Category C** ([https://hs-articulation.ucop.edu/guide/a-g-subject-requirements/c-mathematics/](https://hs-articulation.ucop.edu/guide/a-g-subject-requirements/c-mathematics/)):

- Courses must be consistent with the **Common Core Standards for Mathematical Practice** (high school), per the content guidelines effective 2025-26.
- Courses "recognize the hierarchical nature of mathematics"; traditional (Algebra 1, Geometry, Algebra 2) or integrated (Math I/II/III) sequences are acceptable.
- Advanced courses for 12th grade with a math prerequisite (pre-calculus, calculus, linear algebra, trigonometry) may satisfy the third required year.
- Courses "based largely on repetition of material from a prerequisite or prior course (e.g., pre-college review) will not be approved."
- A trig-only course fulfills only half a year; multi-term courses beyond one year still satisfy only one year of the requirement.
- Five **Skills Guidelines** (apply knowledge to novel problems, justify solutions, use multiple representations, abstraction/generalization, mathematical modeling) and three **Core Competencies** (recognize when math applies; goals of clarity/brevity/universality/objectivity; fluency with formulas and computation).
- Honors criteria: 3 years of prior C-math prerequisite, primarily Common Core (+) advanced standards, plus the general honors requirements.

### 2.2 The approval pathway: CMP registration, submission, deadlines, review

1. **Institution registration.** "Any public or private California high school may establish an A-G course list if they are accredited... and a diploma-granting institution." Non-California high schools are ineligible "with the exception of online schools that serve California students." Online schools must complete a self-assessment against the National Standards for Quality Online Programs. **Accreditation is mandatory**: BOARS policy (December 2002) requires all schools to be accredited to establish/maintain an A-G list, and "As of July 2020, UC highly recommends accreditation by... WASC for California high schools"; WASC candidacy also qualifies ([https://hs-articulation.ucop.edu/guide/register-your-institution/schools/](https://hs-articulation.ucop.edu/guide/register-your-institution/schools/); confirmed by WASC: "in November 2002 the University of California passed a policy that requires all schools that want UC approval of courses to meet the a-g requirements to be accredited by WASC or at least be an initial WASC candidate" ([https://acswasc.org/about/frequently-asked-questions/](https://acswasc.org/about/frequently-asked-questions/))).
2. **Non-classroom-based / independent-study criteria** (directly relevant to a self-paced platform used at home): students must spend "at least one hour per week per A-G course engaged in interactive instruction and/or academic tutoring/advising"; regular (at least weekly) access to a **subject expert teacher** (defined as: 3+ years teaching the subject at HS/postsecondary level, **or** a teaching credential plus a bachelor's/advanced degree in the subject, **or** Highly Qualified Teacher certification); prompt responses; **all courses must require a final exam or significant final project**; "Major assessments (i.e., unit tests, final exams) shall be **proctored by a qualified professional** (e.g., a school teacher, administrator, or counselor)"; student work evaluated by an impartial professional involved in the student's learning ([https://hs-articulation.ucop.edu/guide/register-your-institution/schools/](https://hs-articulation.ucop.edu/guide/register-your-institution/schools/)).
3. **Course submission via the A-G Course Management Portal (CMP)** ([https://hs-articulation.ucop.edu/agcmp](https://hs-articulation.ucop.edu/agcmp)). Submissions follow the **annual calendar**: Primary Phase **February 1 – June 30** (so lists are current before the August 1 UC application opening); Course Management Month **July** (auto-approval-only types); Supplementary Phase **August 1 – 31** ([https://hs-articulation.ucop.edu/guide/update-your-a-g-list/a-g-course-submission-deadlines/](https://hs-articulation.ucop.edu/guide/update-your-a-g-list/a-g-course-submission-deadlines/)).
4. **Course description requirements** ("Writing A-G courses"): course overview, **unit summary + one sample assignment for every unit** (5–7 sentences each), showing alignment with the subject-area criteria ([https://hs-articulation.ucop.edu/guide/update-your-a-g-list/writing-a-g-courses/](https://hs-articulation.ucop.edu/guide/update-your-a-g-list/writing-a-g-courses/)). Honors-level courses additionally need a **comprehensive final examination or substantive culminating project** (same page).
5. **Review timeline**: "New courses are generally reviewed within **two to four weeks** of submission. However, depending on the time of year and volume of submissions, it may take longer." Non-approved courses may be resubmitted after revisions ([https://hs-articulation.ucop.edu/guide/update-your-a-g-list/submitting-courses/](https://hs-articulation.ucop.edu/guide/update-your-a-g-list/submitting-courses/)).
6. **Renewal/maintenance**: course lists are re-validated annually via the CMP (annual update checklist); course revisions and activations of archived courses are submitted in the same annual windows; course-learning-environment edits can be made year-round ([https://hs-articulation.ucop.edu/guide/update-your-a-g-list/a-g-course-submission-deadlines/](https://hs-articulation.ucop.edu/guide/update-your-a-g-list/a-g-course-submission-deadlines/); [https://hs-articulation.ucop.edu/guide/update-your-a-g-list/annual-a-g-update-checklist/](https://hs-articulation.ucop.edu/guide/update-your-a-g-list/annual-a-g-update-checklist/)). There is no separate "renewal fee" — maintenance runs through the annual update cycle. (For schools, see also the "Course revisions" page: [https://hs-articulation.ucop.edu/guide/update-your-a-g-list/course-revisions/](https://hs-articulation.ucop.edu/guide/update-your-a-g-list/course-revisions/).)

### 2.3 The online course publisher path (the one that fits an online platform)

Since 2019–20 UC no longer maintains A-G reference lists for online course publishers; instead publishers register in UC's **online course publisher directory** ([https://hs-articulation.ucop.edu/agcourselist/publisher-directory](https://hs-articulation.ucop.edu/agcourselist/publisher-directory)) and **high schools** certify + self-report publisher courses onto their own A-G lists ([https://hs-articulation.ucop.edu/guide/update-your-a-g-list/online-courses/](https://hs-articulation.ucop.edu/guide/update-your-a-g-list/online-courses/); FAQ: [https://hs-articulation.ucop.edu/guide/news-resources/updated-online-course-publisher-policy-faqs](https://hs-articulation.ucop.edu/guide/news-resources/updated-online-course-publisher-policy-faqs)).

**Publisher registration criteria** — must meet at least #1–#4 (plus #5 only if instructors are offered):

1. Author their own curriculum;
2. Rely on the student's **home high school** to report grades/credits on the official transcript;
3. Documentation that public high schools in at least **two California school districts** (or two private high schools) have committed to using the publisher's courses for A-G purposes;
4. Maintain **at least two years of enrollment data** on California high school students "for UC audit purposes";
5. *If the publisher employs instructors:* **WASC accreditation as a supplementary education center/program** ([https://hs-articulation.ucop.edu/guide/register-your-institution/online-course-publishers/](https://hs-articulation.ucop.edu/guide/register-your-institution/online-course-publishers/); same list on [https://hs-articulation.ucop.edu/guide/update-your-a-g-list/online-courses/](https://hs-articulation.ucop.edu/guide/update-your-a-g-list/online-courses/)).

**Annual verification:** UC randomly reviews self-reported publisher courses May 1–15 (notification) with submissions via CMP May 15–June 30 (one resubmission, due June 30); courses that fail alignment are removed from school lists; **two strikes removes the publisher from the directory** until courses earn Quality Matters certification. QM-certified courses are exempt from the annual verification and are spotlighted in the directory ([https://hs-articulation.ucop.edu/guide/update-your-a-g-list/online-courses/](https://hs-articulation.ucop.edu/guide/update-your-a-g-list/online-courses/)).

**UC's 10 online learning guidelines** (institution-level requirements when online courses appear on an A-G list): equal access; development by content experts + ed-tech experts; alignment with A-G criteria; "opportunities for **substantial interactions between students and the teacher, and between students and other students**"; frequent varied assessment with prompt feedback; **qualified teachers who are content experts**, with adequate professional development and reasonable caseloads; adequate technology infrastructure; advising; local access to qualified professionals/paraprofessionals; and — directly relevant to self-paced platforms — "Institutions must have a process in place to ensure that the person submitting material for assessment is actually the student enrolled in the online course" ([https://hs-articulation.ucop.edu/guide/update-your-a-g-list/online-courses/](https://hs-articulation.ucop.edu/guide/update-your-a-g-list/online-courses/)).

UC also points schools and publishers to the **Quality Matters K-12 Publisher Rubric** and the **National Standards for Quality Online Courses** (NSQOL) for assessment integrity and course quality ([https://hs-articulation.ucop.edu/guide/update-your-a-g-list/online-courses/](https://hs-articulation.ucop.edu/guide/update-your-a-g-list/online-courses/)).

**Seat time vs mastery:** UC does not prescribe minimum seat time for a-g courses (and the revised Category C guidelines explicitly reject only "repetition of material"); the operative requirements are content coverage, rigor, the final-exam requirement, proctored major assessments, and teacher interaction (independent-study criteria, [https://hs-articulation.ucop.edu/guide/register-your-institution/schools/](https://hs-articulation.ucop.edu/guide/register-your-institution/schools/)). A mastery-based course can therefore qualify **if** it demonstrates standards coverage and includes supervised, proctored, instructor-involved assessment — the exact three things a pure self-paced app lacks by default.

### 2.4 Interaction with California graduation requirements

The state minimum for a diploma (Education Code §51225.3): 3 years English, 2 years math, 2 years science, 3 years social studies (incl. U.S. history, world history, ½-year civics, ½-year economics), 1 year VPA/world language/CTE, 2 years PE; a one-semester ethnic studies course commencing with the class of 2029-30; and a separate one-semester personal finance course commencing with the class of 2030-31 ([https://leginfo.legislature.ca.gov/faces/codes_displaySection.xhtml?lawCode=EDC&sectionNum=51225.3](https://leginfo.legislature.ca.gov/faces/codes_displaySection.xhtml?lawCode=EDC&sectionNum=51225.3); CDE's summary page is [https://www.cde.ca.gov/ci/gs/hs/hsgrgen.asp](https://www.cde.ca.gov/ci/gs/hs/hsgrgen.asp) — unreachable from this environment; the statute is the authoritative source).

How WASC and a-g interact with this: (1) WASC accreditation is the *quality stamp* the UC requires schools to hold; it does **not** itself approve courses — a-g approval is a separate UC process ([https://hs-articulation.ucop.edu/guide/register-your-institution/schools/](https://hs-articulation.ucop.edu/guide/register-your-institution/schools/)). (2) a-g (UC/CSU admission) is *stricter* than the state diploma minimum (e.g., 3 vs 2 math years, geometry requirement), and the statute itself draws the distinction when districts substitute CTE courses for the VPA/world-language requirement (EC §51225.3(a)(1)(E)(iv)) ([https://leginfo.legislature.ca.gov/faces/codes_displaySection.xhtml?lawCode=EDC&sectionNum=51225.3](https://leginfo.legislature.ca.gov/faces/codes_displaySection.xhtml?lawCode=EDC&sectionNum=51225.3)). (3) WASC ACS and CDE have a formal joint-review relationship (the "Focus on Learning" joint process) for California public schools ([https://acswasc.org/about/acs-wasc-overview/](https://acswasc.org/about/acs-wasc-overview/)).

---

## Section 3 — Mathua gap map

Legend: **SUPPORTED** = evidence exists today; **PARTIAL** = exists but lacks a required form/field/process; **MISSING** = does not exist. "Evidence needed" column names the artifact an auditor would request.

### 3.1 WASC ACS (SEP / online-program track)

| # | Requirement (source) | Mathua status | Evidence needed |
|---|---|---|---|
| W1 | Be a defined organization with mission & schoolwide learner outcomes (Standards A1, A4) ([acswasc.org/about/acs-wasc-standards](https://acswasc.org/about/acs-wasc-standards/)) | **PARTIAL** — mission exists implicitly ("mastery-gated, locally-first", `README.md:14`); no published SLOs or self-study artifact | Written mission + measurable learner-outcome statements; continuous-improvement plan with data |
| W2 | Self-study evidencing student achievement data (B3: "varied, comprehensive, and valid assessments"; data analysis) ([standards](https://acswasc.org/about/acs-wasc-standards/)) | **SUPPORTED** — per-concept mastery states, SM-2 histories, `DiagnosticReport{MasteryLevels, CompletionEstimates, GapsByDomain}` (`internal/diagnostic/report.go:11`), attempt logs (`internal/storage/store.go`) | Anonymized aggregate reports (completion, mastery, retention by concept/course) |
| W3 | Student enrollment ≥ 6 full-time students ([FAQ](https://acswasc.org/about/frequently-asked-questions/)) | **MISSING** — no enrollment/roster concept; guests are ephemeral (`lib/auth.ts:4` `ensureGuestId`) | Persistent student records, rosters, enrollment dates |
| W4 | Qualified faculty/staff; "highly qualified" subject teachers (Standard A3; [FAQ](https://acswasc.org/about/frequently-asked-questions/)) | **MISSING** — no teacher/instructor concept anywhere in the schema (`system-design.md:68` student table; no staff table) | Instructor-of-record field + credential/expertise documentation per course |
| W5 | Oversight of curriculum + delivery by qualified teacher; grades/transcripts produced by teacher, not algorithm alone ([FAQ](https://acswasc.org/about/frequently-asked-questions/); independent-study criteria in §2.2) | **MISSING** — completion is algorithm-determined (streak+time thresholds, `internal/mastery/machine.go:22`) | Teacher sign-off workflow on mastery awards/transcripts |
| W6 | Records retention plan; transcripts housed for many years ([FAQ](https://acswasc.org/about/frequently-asked-questions/)) | **PARTIAL** — CSV transcript export exists (`internal/server/server.go:1491` `mathua-transcript.csv`, 19 courses from `data/courses.json:1`); no retention policy, no archive guarantee | Written records-retention policy + exportable archival format per student |
| W7 | Accountability/compliance with applicable laws (Standard A6) ([standards](https://acswasc.org/about/acs-wasc-standards/)) | **MISSING** — no COPPA/FERPA policies, no privacy officer, no complaint process | Privacy policy, FERPA/COPPA compliance docs (§4.6, §4.7) |
| W8 | "Respect for all students" / equitable access, accommodations (Standards C1, C2; adaptive pacing as "systems of support") ([standards](https://acswasc.org/about/acs-wasc-standards/)) | **PARTIAL** — per-concept adaptive difficulty and remediation exist (`internal/engine/engine.go`); no accommodations/time-multiplier setting shipped (planned: `improve.md:139`) | Accessibility statement (WCAG), accommodation policy, ADA-accessible UI |
| W9 | Annual reports + six-year cycle (follow-up process) ([overview](https://acswasc.org/about/acs-wasc-overview/)) | **MISSING** — no reporting process (product-appropriate) | Designated accreditation liaison + annual data pack |

### 3.2 UC a-g — general pathway prerequisites

| # | Requirement (source) | Mathua status | Evidence needed |
|---|---|---|---|
| G1 | Accredited (WASC) as a school **or** SEP (BOARS 2002 policy; publisher criterion #5) ([ucop schools](https://hs-articulation.ucop.edu/guide/register-your-institution/schools/); [ucop publishers](https://hs-articulation.ucop.edu/guide/register-your-institution/online-course-publishers/)) | **MISSING** — no accreditation | Completed WASC SEP application + initial visit (§1.2) |
| G2 | Publisher criteria #1–4: own curriculum; home-school reports grades; 2 CA districts committed; 2 years of CA enrollment data for audit ([ucop publishers](https://hs-articulation.ucop.edu/guide/register-your-institution/online-course-publishers/)) | **#1 SUPPORTED** (curriculum authored in-repo: `data/lessons/`, generators); **#2 PARTIAL** (transcript CSV exists; no school-partner reporting flow); **#3 MISSING** (no school partnerships); **#4 MISSING** (no CA-enrollment tracking) | Publisher registration docs; partner letters/purchase agreements; enrollment analytics retained ≥2 years |
| G3 | CMP registration + annual submission (Feb 1–Jun 30) ([deadlines](https://hs-articulation.ucop.edu/guide/update-your-a-g-list/a-g-course-submission-deadlines/)) | **MISSING** — no CMP account, no course-description documents | CMP account + per-course submissions |
| G4 | Course descriptions: overview + unit summary + sample assignment per unit ([writing](https://hs-articulation.ucop.edu/guide/update-your-a-g-list/writing-a-g-courses/)) | **MISSING** — lessons are concept-keyed markdown (`internal/lessons/lesson.go:28`), not unit/course-structured with assignments | Per-course syllabus docs mapping units → A-G criteria |
| G5 | All courses require final exam or significant final project; proctoring of major assessments ([independent study](https://hs-articulation.ucop.edu/guide/register-your-institution/schools/)) | **MISSING** — no capstone/final exam; no proctoring (quizzes are timed but unattended) | Final-exam component + proctoring policy (§4.3) |
| G6 | Online guidelines: teacher-student and student-student interaction; qualified teachers; identity verification (guideline #10) ([online courses](https://hs-articulation.ucop.edu/guide/update-your-a-g-list/online-courses/)) | **MISSING** — self-paced, no instructor, no peer interaction, no identity check | Instructor-of-record + interaction cadence + identity/anti-cheat process |
| G7 | QM certification or annual verification readiness ([online courses](https://hs-articulation.ucop.edu/guide/update-your-a-g-list/online-courses/)) | **MISSING** | Either QM K-12 Publisher Rubric certification or readiness to submit course descriptions May 15–Jun 30 |

### 3.3 UC a-g — Category C (Mathematics) content criteria

| # | Criterion (source: [C-mathematics](https://hs-articulation.ucop.edu/guide/a-g-subject-requirements/c-mathematics/)) | Mathua status | Evidence needed |
|---|---|---|---|
| M1 | Covers elementary algebra, two- and three-dimensional geometry, advanced algebra; 3 years / 4 recommended | **SUPPORTED** — 19 courses from 4th grade through university (`data/courses.json:1`) including geometry/geo content; Algebra 1/2, Geometry, Precalc, Calc, LinAlg | Course→CCSM standards coverage matrix |
| M2 | Common Core-aligned; hierarchical prerequisites; no repetition-only courses | **SUPPORTED** — DAG with hard prerequisite gating (`internal/concepts/dag.go`, `Available` filter) | DAG→CCSM alignment doc; demonstrate non-repetition |
| M3 | Advanced 12th-grade courses w/ math prereq satisfy 3rd year | **SUPPORTED** — calc/linAlg/discrete/stat courses exist | Prereq listing per advanced course |
| M4 | Trig-only ≠ full year; multi-term courses = one year only | **PARTIAL** — course granularity is year-shaped by design, but no per-course credit mapping documented | Credit-value statements per course |
| M5 | Skills guidelines (novel problems, justification, multiple representations, modeling) | **PARTIAL** — generators produce novel problems (`internal/generator/registry.go`); **no justification/open-response grading**; lessons are worked examples not proof practice | Assessment samples showing problem-solving depth |
| M6 | Honors designation: comprehensive final exam/project + 3-yr prereq + (+) standards ([writing](https://hs-articulation.ucop.edu/guide/update-your-a-g-list/writing-a-g-courses/)) | **PARTIAL** — honors-level content exists (calc, abstract algebra) but no honors designation docs | Honors application w/ final-exam design |

### 3.4 CA state law

| # | Requirement (source) | Mathua status | Evidence needed |
|---|---|---|---|
| S1 | Diploma minimums (EC §51225.3): 2 math years etc. ([leginfo](https://leginfo.legislature.ca.gov/faces/codes_displaySection.xhtml?lawCode=EDC&sectionNum=51225.3)) | **SUPPORTED** — course catalog far exceeds state minimums | — (no change needed; affects partner schools) |
| S2 | WASC/CDE alignment for public schools ([overview](https://acswasc.org/about/acs-wasc-overview/)) | Not applicable until a school-partner model exists | — |

---

## Section 4 — Recommended roadmap (ordered, small steps)

Each step is an artifact or a small product change; none require the full accreditation process to begin.

1. **Establish the legal/org identity and privacy base.** Incorporate an entity with a published mission, then publish: privacy policy, terms, and complaint process. **COPPA** applies because Mathua serves children under 13 (16 CFR Part 312 — operators of services "directed to children under 13" must provide notice and obtain verifiable parental consent; [https://www.ftc.gov/legal-library/browse/rules/childrens-online-privacy-protection-rule-coppa](https://www.ftc.gov/legal-library/browse/rules/childrens-online-privacy-protection-rule-coppa)). **FERPA** applies once a school relationship exists (34 CFR 99.1 — applies to agencies/institutions receiving ED funds; the definition of "education records" and PII at 99.3, 45-day inspection at 99.10, annual notification at 99.7 are the core obligations; [https://www.ecfr.gov/current/title-34/subtitle-A/part-99](https://www.ecfr.gov/current/title-34/subtitle-A/part-99)). **WCAG 2.2 AA** conformance for the web app (the W3C standard: [https://www.w3.org/WAI/standards-guidelines/wcag/](https://www.w3.org/WAI/standards-guidelines/wcag/)).
2. **Add an instructor-of-record / teacher model** (schema + UI). This one field unblocks W4, W5, G6, and WASC Standard A3. Also add a teacher-sign-off step on transcript generation (grades produced by a qualified teacher, per WASC FAQ for independent-study programs).
3. **Course metadata + syllabus pack.** Add per-course (not per-concept) metadata: credit value, grade level, CCSM alignment, prerequisites — then produce the "unit summary + sample assignment" documents the CMP requires ([https://hs-articulation.ucop.edu/guide/update-your-a-g-list/writing-a-g-courses/](https://hs-articulation.ucop.edu/guide/update-your-a-g-list/writing-a-g-courses/)). Map the DAG onto Common Core (the DAG's hierarchy is a natural fit for M2).
4. **Assessment & proctoring policy.** Add a final-exam/capstone component per course (G5) and a proctoring policy (live proctor or proctored assessment window) for major assessments — required by the independent-study criteria ([https://hs-articulation.ucop.edu/guide/register-your-institution/schools/](https://hs-articulation.ucop.edu/guide/register-your-institution/schools/)) and UC online guideline #10 (identity verification) ([https://hs-articulation.ucop.edu/guide/update-your-a-g-list/online-courses/](https://hs-articulation.ucop.edu/guide/update-your-a-g-list/online-courses/)).
5. **Seat-time ↔ mastery equivalence documentation.** UC does not require seat time, but auditors expect proof that mastery thresholds substitute for it: write a white paper mapping streak+time thresholds, SM-2 spacing, and per-concept mastery to "hours of engaged learning" and to the Category C skills guidelines (M5), using `internal/mastery/machine.go` + `DiagnosticReport` data.
6. **Records retention & audit export.** Formalize retention (FERPA-driven, 45-day access; WASC wants a long-term transcript plan) and expose a per-student archival export (transcript + attempt history) beyond the current CSV ([https://www.ecfr.gov/current/title-34/subtitle-A/part-99](https://www.ecfr.gov/current/title-34/subtitle-A/part-99); [https://acswasc.org/about/frequently-asked-questions/](https://acswasc.org/about/frequently-asked-questions/)).
7. **Enrollment tracking.** Add persistent student records with enrollment dates and CA-school affiliation (needed for publisher criterion #4: two years of CA enrollment data) ([https://hs-articulation.ucop.edu/guide/register-your-institution/online-course-publishers/](https://hs-articulation.ucop.edu/guide/register-your-institution/online-course-publishers/)).
8. **Quality review against QM/NSQ.** Self-audit courses against the Quality Matters K-12 Publisher Rubric and National Standards for Quality Online Courses — QM certification is the shortcut that exempts courses from UC's annual verification and is the remediation path after a directory removal ([https://hs-articulation.ucop.edu/guide/update-your-a-g-list/online-courses/](https://hs-articulation.ucop.edu/guide/update-your-a-g-list/online-courses/)).
9. **WASC SEP application (parallel, after steps 1–2 and 6–7).** $160 application + $850 initial visit, ≥6 enrolled students, initial visit no earlier than second semester of operation; target candidacy then full self-study by year 3 ([https://acswasc.org/beginning-the-accreditation-process/](https://acswasc.org/beginning-the-accreditation-process/); [https://acswasc.org/supplementary-education-programs/](https://acswasc.org/supplementary-education-programs/)).
10. **UC publisher registration + partner schools.** Secure two committed CA school districts (criterion #3), then register in the CMP and have partners certify + self-report courses during Feb 1–Jun 30 ([https://hs-articulation.ucop.edu/guide/update-your-a-g-list/online-courses/](https://hs-articulation.ucop.edu/guide/update-your-a-g-list/online-courses/)).

**Key sequencing logic:** (1)–(7) are prerequisite artifacts for both WASC and UC; (9) WASC SEP accreditation is a hard precondition for the publisher-with-instructors model (UC criterion #5) and for schools' own a-g list eligibility; (10) requires (9) if Mathua staffs instructors, and requires (3)–(5) before courses survive UC's review. The cheapest wins are steps 1–3 (all in-repo work, no fees).

---

## Sources (all fetched 2026-08-29)

**WASC ACS**
- SEP track, definition, fees: https://acswasc.org/supplementary-education-programs/
- Beginning the process (fees, initial visit, statuses): https://acswasc.org/beginning-the-accreditation-process/
- Initial visit procedures: https://acswasc.org/initial-visit-process-and-procedures/
- Initial status options: https://acswasc.org/initial-accreditation-status-options/
- Nine standards: https://acswasc.org/about/acs-wasc-standards/
- Overview, cycle, UC collaboration: https://acswasc.org/about/acs-wasc-overview/
- Fees: https://acswasc.org/about/acs-wasc-fees/
- FAQ (online schools/iNACOL, UC 2002 policy, 6 students, records): https://acswasc.org/about/frequently-asked-questions/

**WSCUC**
- About accreditation: https://www.wscuc.org/about/about-accreditation/
- Becoming accredited: https://www.wscuc.org/resources/becoming-accredited/
- 2023 Handbook (Standards/CFRs, candidacy, reaffirmation): https://www.wscuc.org/handbook2023/
- Dues and Fees Schedule (PDF): https://wascsenior.box.com/s/yinvbvulfn2vtjg60tb0fqkjs5jdrkef

**UC a-g**
- UC admissions — first-year requirements: https://admission.universityofcalifornia.edu/admission-requirements/freshman-requirements/
- Subject requirement (A-G), Category C: https://admission.universityofcalifornia.edu/admission-requirements/freshman-requirements/subject-requirement-a-g.html
- Counselors — first-year requirements: https://admission.universityofcalifornia.edu/counselors/preparing-freshman-students/freshman-requirements.html
- A-G Policy Resource Guide (index): https://hs-articulation.ucop.edu/guide/
- Category C course criteria: https://hs-articulation.ucop.edu/guide/a-g-subject-requirements/c-mathematics/
- Schools registration (accreditation policy, independent-study criteria): https://hs-articulation.ucop.edu/guide/register-your-institution/schools/
- Online course publishers (registration criteria): https://hs-articulation.ucop.edu/guide/register-your-institution/online-course-publishers/
- Online courses policy (10 guidelines, verification): https://hs-articulation.ucop.edu/guide/update-your-a-g-list/online-courses/
- Online publisher policy FAQ: https://hs-articulation.ucop.edu/guide/news-resources/updated-online-course-publisher-policy-faqs
- Submission deadlines: https://hs-articulation.ucop.edu/guide/update-your-a-g-list/a-g-course-submission-deadlines/
- Submitting courses (2–4 week review): https://hs-articulation.ucop.edu/guide/update-your-a-g-list/submitting-courses/
- Writing A-G courses: https://hs-articulation.ucop.edu/guide/update-your-a-g-list/writing-a-g-courses/
- Annual update checklist: https://hs-articulation.ucop.edu/guide/update-your-a-g-list/annual-a-g-update-checklist/
- Course Management Portal (CMP): https://hs-articulation.ucop.edu/agcmp
- A-G Course List: https://hs-articulation.ucop.edu/agcourselist
- Publisher directory: https://hs-articulation.ucop.edu/agcourselist/publisher-directory
- Quick Reference Guide to UC Admissions (principal certification, p.17): https://admission.universityofcalifornia.edu/counselors/_files/documents/quick-reference-guide-to-uc-admissions.pdf
- Webinar "Making sense of online coursework": https://admission.universityofcalifornia.edu/counselors/_files/documents/counselor-webinars/hs-webinar-online-coursework-3.24.26.pdf

**California**
- Education Code §51225.3 (graduation requirements): https://leginfo.legislature.ca.gov/faces/codes_displaySection.xhtml?lawCode=EDC&sectionNum=51225.3
- CDE high school graduation requirements (unreachable from this environment): https://www.cde.ca.gov/ci/gs/hs/hsgrgen.asp

**Privacy / accessibility**
- COPPA rule (16 CFR Part 312): https://www.ftc.gov/legal-library/browse/rules/childrens-online-privacy-protection-rule-coppa
- FERPA regulations (34 CFR Part 99): https://www.ecfr.gov/current/title-34/subtitle-A/part-99
- WCAG 2 overview (W3C): https://www.w3.org/WAI/standards-guidelines/wcag/

**Referenced by UC (third-party quality standards)**
- Quality Matters K-12 Publisher Rubric: https://www.qualitymatters.org/sites/default/files/PDFs/StandardsfromtheK-12PublisherRubric.pdf
- National Standards for Quality Online Courses: https://www.nsqol.org/

**Mathua (repo, not web sources)**
- `README.md:14,26,66-89` — mission, mastery gating, SM-2; `data/courses.json:1` — 19 courses; `internal/diagnostic/report.go:11-33` — `DiagnosticReport{PlacementCourseID, MasteryLevels, CompletionEstimates, GapsByDomain}`; `internal/server/server.go:1448-1491` — `/api/transcript?format=csv`; `internal/mastery/machine.go:22` — threshold gating; `internal/lessons/lesson.go:28` — lesson loading; `improve.md:139,148` — planned accommodations, accreditation-track design.
