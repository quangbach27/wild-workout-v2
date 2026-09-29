import type { NextConfig } from 'next';

const nextConfig: NextConfig = {
  reactCompiler: true,
  // Emit a minimal self-contained server for Docker deployment
  output: 'standalone',
  poweredByHeader: false,
};

export default nextConfig;
