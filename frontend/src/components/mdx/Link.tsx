'use client';

import { FC } from 'react';

const EXTERNAL_HREF_REGEX = /https?:\/\//;

const A: FC<{ href: string } & React.HTMLAttributes<HTMLAnchorElement>> = ({ href, ...props }) => {
  const isOutbound = EXTERNAL_HREF_REGEX.test(href);
  const isOfficial = String(href).includes('lobechat') || String(href).includes('lobehub');
  return (
    <a
      data-wml-openURL={href}
      rel={isOutbound && !isOfficial ? 'nofollow' : undefined}
      target={isOutbound ? '_blank' : undefined}
      {...props}
    />
  );
};

export default A;
