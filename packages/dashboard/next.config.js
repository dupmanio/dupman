/** @type {import('next').NextConfig} */
const nextConfig = {
  publicRuntimeConfig: {
    DUPMAN_API: process.env.DUPMAN_API,
  },
};

module.exports = nextConfig;
