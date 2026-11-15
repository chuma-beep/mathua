import type { MetadataRoute } from 'next'
import { SITE_URL } from '@/lib/site'

// Next 16 requires an explicit static marker for metadata routes under
// `output: 'export'`; without it the build fails while collecting page data.
export const dynamic = 'force-static'

export default function robots(): MetadataRoute.Robots {
  return {
    rules: [
      {
        userAgent: '*',
        allow: '/',
        disallow: ['/api/', '/admin/'],
      },
    ],
    sitemap: `${SITE_URL}/sitemap.xml`,
  }
}
