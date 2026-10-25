import type { Metadata } from 'next'
import Link from 'next/link'
import Header from '../../components/Header'
import Footer from '../../components/Footer'
import BottomTabs from '../../components/BottomTabs'
import ContactBox from '../../components/ContactBox'

export const metadata: Metadata = {
  title: 'Contact, Mathua',
  description:
    'How to report a problem with a question, a bug in the site, or your account, and who reads it.',
}

const monoFont = "var(--font-jetbrains-mono), 'JetBrains Mono', monospace"
const bodyFont = "var(--font-space-grotesk), 'Space Grotesk', serif"

const textStyle: React.CSSProperties = {
  fontFamily: bodyFont,
  fontSize: '0.95rem',
  color: 'var(--text-secondary)',
  lineHeight: 1.85,
  marginBottom: '1rem',
}

const h2Style: React.CSSProperties = {
  fontFamily: bodyFont,
  fontWeight: 400,
  fontSize: '1.3rem',
  color: 'var(--text-primary)',
  marginTop: '2.5rem',
  marginBottom: '1rem',
  borderBottom: '0.5px solid var(--border)',
  paddingBottom: '0.5rem',
}

const linkStyle: React.CSSProperties = {
  color: 'var(--text-primary)',
  textDecoration: 'underline',
}

// The public route to a person.
//
// Three problems arrive here and only one of them belongs in the in-app report box:
// a wrong expected answer needs the captured question to be reproducible at all, a
// broken account needs a channel that is not public, and a code bug is best served
// by an issue tracker where a fix can be referenced. Saying which is which is the
// whole job of this page, because a learner who sends an account problem to a public
// tracker has published it.
//
// No email address: an address on a page is harvested. The box has no separate
// inbox either -- it writes into the report queue, the page the administrator opens
// whenever a question is reported, so a message cannot quietly pile up somewhere
// nobody looks. X and the issue tracker are the two-way channels.
export default function ContactPage() {
  return (
    <>
      <Header />
      <div className="max-w-[680px] mx-auto px-4 sm:px-6 py-12 sm:py-20 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-20 overflow-x-hidden min-w-0">
        <Link href="/" style={{ fontFamily: monoFont, fontSize: '13px', color: 'var(--text-muted)' }}>
          ← Back
        </Link>
        <div style={{ fontFamily: monoFont, fontSize: '12px', color: 'var(--text-muted)', marginTop: '1.5rem', marginBottom: '0.5rem' }}>
          WHO TO TELL
        </div>
        <h1
          style={{
            fontFamily: "var(--font-space-grotesk), 'Space Grotesk', serif",
            fontWeight: 400,
            fontSize: '1.9rem',
            color: 'var(--text-primary)',
            marginBottom: '2rem',
          }}
        >
          Contact
        </h1>

        <p style={textStyle}>
          Mathua is built and run by one person, so the same person reads all of this. Which
          channel to use depends on what went wrong: a question, the site, or your account.
        </p>

        <h2 style={h2Style}>Leave a message</h2>
        <p style={textStyle}>
          If you would rather write here than on another site, use the box. It reaches the
          same queue as a report about a question, and one person reads it. There is no
          automatic reply, so use X below if you need an answer.
        </p>
        <ContactBox />

        <h2 style={h2Style}>Discord</h2>
        <p style={textStyle}>
          The{' '}
          <a href="https://discord.gg/qXPeQ2jys" target="_blank" rel="noreferrer" style={linkStyle}>
            Mathua Discord server
          </a>{' '}
          is the fastest way to reach other learners and the administrator. Ask about a lesson,
          say what you are working on, or explain a problem that is easier to talk through than
          to type into a form. Keep anything private, an account problem, a password, a message
          meant for one person, to the box above, not the public server.
        </p>

        <h2 style={h2Style}>A question is wrong</h2>
        <p style={textStyle}>
          Use the report control on the question itself in Learn, a quiz, a review, a diagnostic
          or a concept reference. Choose <em>Report a problem</em>, then say what you expected.
          The report keeps the question, the answer you gave, the answer Mathua expected and how
          your answer was being read, so the problem can be replayed exactly instead of guessed
          at from a description. In History, open the concept and report from there, or use the
          message box and include the concept and question.
        </p>
        <p style={textStyle}>
          That is the right route for a wrong expected answer, a question that cannot be typed
          into the box, a figure that does not load, or a mistake in an explanation.
        </p>

        <h2 style={h2Style}>The site is broken, or your account is</h2>
        <p style={textStyle}>
          Message me on X:{' '}
          <a
            href="https://x.com/wisdomsglow"
            target="_blank"
            rel="noreferrer"
            style={linkStyle}
          >
            @wisdomsglow
          </a>
          . Account problems do not belong in a public issue tracker, because everything posted
          there is visible to everyone. If you cannot sign in, say which page you were on and
          what you saw.
        </p>
        <p style={textStyle}>
          Never send your password to anyone, including me. Nothing in Mathua will ask for it.
        </p>

        <h2 style={h2Style}>A bug in Mathua itself</h2>
        <p style={textStyle}>
          If a page does not load, a layout breaks, or a button does nothing, and the problem is
          safe to discuss in public, open an issue:{' '}
          <a
            href="https://github.com/chuma-beep/mathua/issues"
            target="_blank"
            rel="noreferrer"
            style={linkStyle}
          >
            github.com/chuma-beep/mathua/issues
          </a>
          . Include the page address, what you did, what you expected, what happened, and your
          browser and device if you know them.
        </p>

        <h2 style={h2Style}>What not to send</h2>
        <p style={textStyle}>
          Do not put an email address, a password or anything private in a public issue. Mathua
          will not ask for them.
        </p>

        <h2 style={h2Style}>What happens next</h2>
        <p style={textStyle}>
          Every report is read. Mathua is in beta and maintained by one person, so a reply can
          take a while. When a fix ships, the reports that led to it are closed with the reason
          recorded against each one, so the trail says what was wrong and what changed.
        </p>
      </div>
      <Footer />
      <BottomTabs />
    </>
  )
}
