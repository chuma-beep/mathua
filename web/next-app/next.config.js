/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'export',
  images: { unoptimized: true },
  // Keep Turbopack rooted at the app. Without this it can walk above the
  // directory and read a package-lock.json outside the repository.
  turbopack: { root: __dirname },
  // The dev server is read from a phone on the LAN and over the tailnet. Next
  // blocks dev assets from an origin it was not told about, which leaves the
  // page shell rendering and its JavaScript never loading.
  allowedDevOrigins: ['192.168.1.16', '100.89.234.95', 'localhost', '127.0.0.1'],
}

module.exports = nextConfig
