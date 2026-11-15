import type { MetadataRoute } from 'next'
import { INDEXABLE_ROUTES, SITE_URL } from '@/lib/site'

// Next 16 requires an explicit static marker for metadata routes under
// `output: 'export'`; without it the build fails while collecting page data.
export const dynamic = 'force-static'

export default function sitemap(): MetadataRoute.Sitemap {
  const lastModified = new Date()
  return INDEXABLE_ROUTES.map(route => ({
    url: `${SITE_URL}${route}`,
    lastModified,
    changeFrequency: route === '/' ? 'weekly' : 'monthly',
    priority: route === '/' ? 1 : route.startsWith('/docs') ? 0.5 : 0.7,
  }))
}
