/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'export',
  images: { unoptimized: true },
  // Keep Turbopack inside the app directory. Without this it walks up to
  // /home/wisdom and reads a package-lock.json outside the repository.
  turbopack: { root: __dirname },
}

module.exports = nextConfig
