/** @type {import('next').NextConfig} */
const nextConfig = {
  publicRuntimeConfig: {
    DUPMAN_API: process.env.DUPMAN_API,
  },
  experimental: {
    instrumentationHook: true,
  },
};

module.exports = nextConfig;
