/** @type {import('next').NextConfig} */
const nextConfig = {
  publicRuntimeConfig: {
    DUPMAN_API: process.env.DUPMAN_API,
    PREVIEW_API: process.env.PREVIEW_API,
  },
};

module.exports = nextConfig;
