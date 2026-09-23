// Display helpers for concept identity lines (concept/page.tsx).
//
// Subdomains are namespaced `domain.subtopic` throughout data/concepts/*,
// so rendering `{domain}.{subdomain}` doubles the domain
// ("abstract_algebra.abstract_algebra.rings"). Show the leaf instead.

/** Leaf segment of a namespaced subdomain ("a.rings" -> "rings"). */
export function subdomainLeaf(subdomain: string): string {
  const leaf = subdomain.split('.').pop() ?? ''
  return leaf || subdomain
}

/** "domain.leaf · id" scope line for the concept header. */
export function formatConceptPath(domain: string, subdomain: string, id: string): string {
  const leaf = subdomainLeaf(subdomain)
  const scope = leaf && leaf !== domain ? `${domain}.${leaf}` : domain
  return `${scope} · ${id}`
}
