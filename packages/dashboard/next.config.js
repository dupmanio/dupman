/** @type {import('next').NextConfig} */
const nextConfig = {
  publicRuntimeConfig: {
    DUPMAN_LANDING: process.env.DUPMAN_LANDING,
    DUPMAN_API: process.env.DUPMAN_API,
  },
};

module.exports = nextConfig;
