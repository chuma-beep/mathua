/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'export',
  images: { unoptimized: true },
  // Keep Turbopack rooted at the app. Without this it can walk above the
  // directory and read a package-lock.json outside the repository.
  turbopack: { root: __dirname },
}

module.exports = nextConfig
