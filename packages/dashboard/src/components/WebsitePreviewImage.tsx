import * as React from "react";
import Image from "next/image";
import useSWR from "swr";

import { PreviewRepository } from "@/lib/repositories/preview";
import PageLoader from "@/components/PageLoader";

export interface WebsitePreviewImageProps {
  websiteId: string;
}

function WebsitePreviewImage({ websiteId }: WebsitePreviewImageProps) {
  const { data, isLoading, error } = useSWR(
    ["/preview/:websiteId", websiteId],
    ([_, websiteId]) => PreviewRepository.get(websiteId),
    {
      revalidateOnFocus: false,
      revalidateOnReconnect: false,
    },
  );

  return (
    <>
      {isLoading && <PageLoader size={100} />}

      {/* @todo: load fallback image */}
      {error && <PageLoader size={100} />}

      {data?.data?.url && (
        <Image
          src={data.data.url}
          alt="dupman"
          priority={true}
          width="350"
          height="200"
          style={{
            objectFit: "contain",
          }}
        />
      )}
    </>
  );
}

export default WebsitePreviewImage;
