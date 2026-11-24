import type { Metadata } from 'next'
import Link from 'next/link'
import Header from '../../components/Header'
import Footer from '../../components/Footer'
import BottomTabs from '../../components/BottomTabs'

export const metadata: Metadata = {
  title: 'Privacy notice, Mathua',
  description:
    'What the assistant sends to an AI provider, what Mathua keeps and what it never sends.',
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

// The notice the assistant dock links to before a visitor's first message.
//
// Every claim here is one the code can be held to: the prompt carries the message and the
// last turns only (internal/ai/redact.go), the handlers write no message to storage
// (internal/server/ai.go logs failures, not text), the quota row holds a counter and a
// date, and both providers are declared zero-retention in the registry. When one of those
// changes, this page changes with it or it becomes the lie that a notice must never be.
export default function PrivacyPage() {
  return (
    <>
      <Header />
      <div className="max-w-[680px] mx-auto px-4 sm:px-6 py-12 sm:py-20 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-20 overflow-x-hidden min-w-0">
        <Link href="/" style={{ fontFamily: monoFont, fontSize: '13px', color: 'var(--text-muted)' }}>
          ← Back
        </Link>
        <div style={{ fontFamily: monoFont, fontSize: '12px', color: 'var(--text-muted)', marginTop: '1.5rem', marginBottom: '0.5rem' }}>
          WHAT LEAVES MATHUA
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
          Privacy notice
        </h1>

        <p style={textStyle}>
          One feature sends what you type to another company: the assistant. This page says
          what it sends, who receives it and what Mathua keeps. Everything else in Mathua runs
          on Mathua&rsquo;s own server.
        </p>

        <h2 style={h2Style}>The assistant</h2>
        <p style={textStyle}>
          Your message is sent with up to eight earlier turns of the same conversation, so that
          a follow-up question can be answered. Nothing else is sent: not your name, not your
          email address, not your student id, not your answers, not your progress and not the
          expected answer to an exercise.
        </p>

        <h2 style={h2Style}>Who receives it</h2>
        <p style={textStyle}>
          Two providers answer, in this order: Groq, and Cloudflare Workers AI if Groq cannot
          answer. Both are configured for zero retention, which means they process the request
          and do not store it.
        </p>

        <h2 style={h2Style}>What Mathua keeps</h2>
        <p style={textStyle}>
          The assistant&rsquo;s messages are not written to Mathua&rsquo;s database. The server
          keeps a daily count of requests &mdash; tied to your account when you are signed in,
          and to your IP address when you are not &mdash; so that the daily limit can be
          enforced. That counter holds a number and a date, not your text. If a request fails,
          the server logs the failure so it can be fixed; the log records the error, not your
          message.
        </p>

        <h2 style={h2Style}>Accounts</h2>
        <p style={textStyle}>
          If you register, Mathua stores the email address you signed up with and the practice
          history tied to your account. The display name is whatever you type, and you can
          leave it blank. Passwords are stored as bcrypt hashes, never as text.
        </p>

        <h2 style={h2Style} id="account-deletion">Deleting your account</h2>
        <p style={textStyle}>
          Settings has a Delete account control. It asks you to type
          &ldquo;delete my account&rdquo; and, for a password account, the current
          password. It then removes the account and everything owned by it:
          progress, attempts, sessions, quiz history, connected sign-ins,
          settings, plans and avatars. Nothing is kept. Guest progress lives in
          the browser and goes when you clear site data.
        </p>
        <p style={textStyle}>
          If you cannot sign in, open an issue in the repository and name the
          account to remove. That reaches the same queue as a report.
        </p>

        <h2 style={h2Style}>The contact box</h2>
        <p style={textStyle}>
          The box on the contact page writes your message into the same table as a report
          about a question. The text you write is stored, and so is your account id if you
          are signed in, or the guest id your browser already carries if you are not. It is
          read in the moderation queue and kept so that the conversation about a fix has a
          record. Do not put a password in it, and do not put anyone else&rsquo;s details in it.
        </p>

        <h2 style={h2Style}>Age</h2>
        <p style={textStyle}>
          The assistant is for people aged 13 or older. If you are younger than that, do not
          use it. The rest of Mathua works without it.
        </p>

        <h2 style={h2Style}>Changes and questions</h2>
        <p style={textStyle}>
          If what the assistant sends changes, the notice in the panel asks you again rather
          than assuming your earlier answer still applies. Questions about this notice can go
          to the repository&rsquo;s issue tracker:{' '}
          <a
            href="https://github.com/chuma-beep/mathua/issues"
            style={{ color: 'var(--text-primary)', textDecoration: 'underline' }}
          >
            github.com/chuma-beep/mathua
          </a>
          .
        </p>
      </div>
      <Footer />
      <BottomTabs />
    </>
  )
}
