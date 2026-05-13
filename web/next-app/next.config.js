const nextConfig = {
  output: process.env.NODE_ENV === 'development' ? undefined : 'export',
  async rewrites() {
    return [
      { source: '/api/:path*', destination: 'http://localhost:8080/api/:path*' },
    ]
  },
}

module.exports = nextConfig