/** @type {import('next').NextConfig} */
const nextConfig = {
  images: {
    remotePatterns: [
      {
        protocol: "https",
        hostname: "fakeimg.pl",
      },
    ],
  },
  publicRuntimeConfig: {
    DUPMAN_API: process.env.DUPMAN_API,
  },
};

module.exports = nextConfig;
