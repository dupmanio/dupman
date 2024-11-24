/** @type {import('next').NextConfig} */
const nextConfig = {
  publicRuntimeConfig: {
    DUPMAN_API: process.env.DUPMAN_API,
  },
  experimental: {
    instrumentationHook: true,
  },
  typescript: {
    ignoreBuildErrors: true,
  },
};

module.exports = nextConfig;
