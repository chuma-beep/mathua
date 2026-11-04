// One place for the public identity links.
//
// The footer, the contact page and the Creator's Note all point at the same set
// of places. Keeping the URLs here means a handle change is one edit rather than
// three, and a test can assert that every surface publishes the same addresses.
export const SUPPORT_URL = 'https://selar.com/showlove/wisddom'

export interface CreatorLink {
  label: string
  href: string
}

export const CREATOR_LINKS: CreatorLink[] = [
  { label: 'Website', href: 'https://wisdomac.com' },
  { label: 'X', href: 'https://x.com/wisdomsglow' },
  { label: 'LinkedIn', href: 'https://linkedin.com/in/wisdom-anwaegbu' },
  { label: 'GitHub profile', href: 'https://github.com/chuma-beep' },
]
