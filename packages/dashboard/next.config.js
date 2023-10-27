/** @type {import('next').NextConfig} */
const nextConfig = {
  publicRuntimeConfig: {
    DUPMAN_API: process.env.DUPMAN_API,
    PREVIEW_API: process.env.PREVIEW_API,
    NOTIFY_URL: process.env.NOTIFY_URL,
  },
};

module.exports = nextConfig;
