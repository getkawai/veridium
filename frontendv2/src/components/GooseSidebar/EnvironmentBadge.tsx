import React from 'react';
import { Tooltip, TooltipContent, TooltipTrigger } from '../ui/Tooltip';

const EnvironmentBadge: React.FC = () => {
  const isDevelopment = import.meta.env.DEV;

  // Don't show badge in production
  if (!isDevelopment) {
    return null;
  }

  const tooltipText = 'Alpha';
  const bgColor = 'bg-purple-600';

  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <div
          className={`${bgColor} w-3 h-3 rounded-full cursor-default`}
          data-testid="environment-badge"
          aria-label={tooltipText}
        />
      </TooltipTrigger>
      <TooltipContent side="right">{tooltipText}</TooltipContent>
    </Tooltip>
  );
};

export default EnvironmentBadge;
