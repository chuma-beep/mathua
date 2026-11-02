import type { Metadata } from 'next'
import Link from 'next/link'
import Header from '../../components/Header'
import Footer from '../../components/Footer'
import BottomTabs from '../../components/BottomTabs'
import { SUPPORT_URL } from '../../lib/creatorLinks'

export const metadata: Metadata = {
  title: "Creator's Note, Mathua",
  description: 'Why Mathua exists, from the person building it.',
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

const quoteStyle: React.CSSProperties = {
  fontFamily: bodyFont,
  fontSize: '1.05rem',
  fontStyle: 'italic',
  color: 'var(--text-primary)',
  lineHeight: 1.7,
  borderLeft: '2px solid var(--accent-blue)',
  paddingLeft: '1rem',
  margin: '1.5rem 0',
}

const citeStyle: React.CSSProperties = {
  display: 'block',
  fontFamily: monoFont,
  fontSize: '12px',
  fontStyle: 'normal',
  color: 'var(--text-muted)',
  marginTop: '0.5rem',
}

const ipaStyle: React.CSSProperties = {
  fontFamily: monoFont,
  fontSize: '13px',
  color: 'var(--text-muted)',
  marginBottom: '0.75rem',
}

// The note is the creator's own words, kept as written: sentence shape, capitalisation
// and spelling are part of the voice, so nothing here is edited for house style. The
// structure only marks what the author marked, headings, a quotation and one bolded
// phrase, because those carry meaning in the original.
type Block =
  | { text: string }
  | { parts: Array<string | { strong: string }> }
  | { quote: string; cite: string }
  | { ipa: string }

type Section = { heading: string; blocks: Block[] }

const SECTIONS: Section[] = [
  {
    heading: 'Hello',
    blocks: [
      {
        text: "it's 4:23am, and a '3rd worlder' who doesnt really have any deep experience in maths is writing about a platform he is creating in hopes that it will help people such as myself learn maths.",
      },
    ],
  },
  {
    heading: 'Pitstop',
    blocks: [
      {
        text: "Why did I start this? I started this because I was really desperate to learn maths (gain automaticity) and I couldnt find any better platform other than mathacademy and it was pricey for me. 50 dollars might not be much for most countries in the western world but it is about 66,000 in my countries currency where the minimum wage per month is 77,000... dont mind me just talking economics here. :|",
      },
      {
        text: "This project really is a pitstop. my aim is that with enough support from kind people we'll be able to create something as good as mathacademy, and it will be a project that improves overtime as the community supports it.",
      },
      {
        text: 'Okay so i have been trying to read research papers why? I am transitioning into AI Robotics Engineering. So to know or gain knowledge I have met a lot of maths I have either forgotten, do not understand or have never seen in my life.',
      },
      {
        text: "my honest thought was this 'okay something has to be done about this' I went to look for platforms to fill up the knowledge gap places like khanacademy,educare, openstax are cool but in those places you just encounter lessons and move along nothing else, I understand that things like active recall, spaced repetition and interleaving are meant to be your duty but what if the person is 'lazy' but still wants to try? or has never even heard of those terms. Thanks to people like Justin sung and Justin Skycak I just learnt about them myself while researching on how to be a better learner.",
      },
      {
        text: 'I was progressing along the lines until I stumbled on mathacademy, it really is an excellent learning platform. Most of the things I tried to do here are things that are on mathacademy I am trying to create something just as good and over time it can probably even be better. From the first beta version I also added some aspects i think they are missing.',
      },
    ],
  },
  {
    heading: 'an ODE for the learner',
    blocks: [
      {
        quote: 'as long as i am alive i have infinite chances',
        cite: 'Monkey D. Luffy',
      },
      {
        text: 'The truth is that math is brutally hierarchical, I mean take out time to ask the real geniuses around you, assuming you have one of course, they always give the same advice you have to start from the basics.',
      },
      {
        text: 'This makes a question come to mind where exactly? arithmetic? counting? prealgebra? or algebra? anyways that problem is one of the problems I am trying to solve. One thing I want the learner to note is that you get to go much quicker than you would expect if you just show up everyday and practice the right way I mean The brain is a sucker for repetition.',
      },
      {
        text: 'One analogy is this, imagine you want to write proper english or be creative with english like writing a novel. how far can you go with out never learning consonants or vowels, what an action word is? or what an adjective is? how far can one go in creativity? For the sky scrapper to stand tall the foundations have to be sturdy right?.',
      },
    ],
  },
  {
    heading: 'Automaticity and grunt work',
    blocks: [
      { ipa: '/ˌɔːtəməˈtɪsɪti/' },
      {
        text: 'Automaticity is the ability to perform tasks, behaviors, or mental processes effortlessly and without conscious thought as a result of extensive practice, repetition and learning.',
      },
      {
        text: "Can we do maths with this ability on display? what is required in order to do maths like this? is it just Genius? I acknowledge that there are intellectual differences in individuals but dont count your self short. All great individuals have a period of doing low-level grunt work in a whatever manner some structured or in a disciplined and passionate way, because of intellectual differences some people might make it look easy but the grunt work was done regardless.",
      },
      {
        parts: [
          "Luffy didn't immediately go striaght to gear 5 remember? instead he went through gear 1 2 3 4 then gear 5, and gear 4 even had ",
          { strong: 'Bounceman/Boundman, Tankman and Snakeman' },
          '.',
        ],
      },
      {
        text: 'how far we can go with delibrate practice and focused effort?',
      },
    ],
  },
]

export default function NotePage() {
  return (
    <>
      <Header />
      <div className="max-w-[680px] mx-auto px-4 sm:px-6 py-12 sm:py-20 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-20 overflow-x-hidden min-w-0">
        <Link href="/" style={{ fontFamily: monoFont, fontSize: '13px', color: 'var(--text-muted)' }}>
          ← Back
        </Link>
        <div style={{ fontFamily: monoFont, fontSize: '12px', color: 'var(--text-muted)', marginTop: '1.5rem', marginBottom: '0.5rem' }}>
          WHY MATHUA EXISTS
        </div>
        <h1 style={{
          fontFamily: "var(--font-space-grotesk), 'Space Grotesk', serif",
          fontWeight: 400,
          fontSize: '1.9rem',
          color: 'var(--text-primary)',
          marginBottom: '2rem',
          borderBottom: '0.5px solid var(--border)',
          paddingBottom: '0.5rem',
        }}>
          Creator&apos;s Note
        </h1>

        {SECTIONS.map(section => (
          <section key={section.heading}>
            <h2 style={h2Style}>{section.heading}</h2>
            {section.blocks.map((block, i) => {
              if ('quote' in block) {
                return (
                  <blockquote key={i} style={quoteStyle}>
                    {`\u201c${block.quote}\u201d`}
                    <cite style={citeStyle}>{block.cite}</cite>
                  </blockquote>
                )
              }
              if ('ipa' in block) {
                return <p key={i} style={ipaStyle}>{block.ipa}</p>
              }
              if ('parts' in block) {
                return (
                  <p key={i} style={textStyle}>
                    {block.parts.map((part, j) =>
                      typeof part === 'string' ? part : <strong key={j}>{part.strong}</strong>,
                    )}
                  </p>
                )
              }
              return <p key={i} style={textStyle}>{block.text}</p>
            })}
          </section>
        ))}

        <section>
          <h2 style={h2Style}>If you want to help</h2>
          <p style={textStyle}>
            Mathua is free and open source. If it has been useful and you want to support the work,
            you can{' '}
            <a
              href={SUPPORT_URL}
              target="_blank"
              rel="noreferrer"
              style={{ color: 'var(--text-primary)', textDecoration: 'underline' }}
            >
              tip the project on Myhappr
            </a>
            . That support helps keep Mathua free, pays for hosting and diagrams, and funds ongoing
            improvements to the curriculum and platform, including future app work such as a mobile
            app. There is no paywall, and nothing here changes if you cannot contribute.
          </p>
        </section>

        <div className="flex gap-3 justify-center mt-10 max-sm:flex-col max-sm:items-center">
          <Link href="/onboard" style={{ fontFamily: monoFont, fontSize: '13px', color: 'var(--bg)', background: 'var(--accent-blue)', padding: '10px 22px', textDecoration: 'none' }}>
            Start diagnostic test →
          </Link>
          <Link href="/learn" style={{ fontFamily: monoFont, fontSize: '13px', color: 'var(--text-secondary)', border: '0.5px solid var(--border-strong)', padding: '10px 22px', textDecoration: 'none' }}>
            Start learning →
          </Link>
        </div>
      </div>
      <Footer />
      <BottomTabs />
    </>
  )
}
