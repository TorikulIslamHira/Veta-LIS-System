import React from 'react';

interface StatusBadgeProps {
  status: 'connected' | 'disconnected' | 'error' | 'pending' | 'normal' | 'abnormal' | 'critical';
  label?: string;
}

const StatusBadge: React.FC<StatusBadgeProps> = ({ status, label }) => {
  const statusConfig = {
    connected: { bg: 'bg-green-100', text: 'text-green-800', label: 'Connected', icon: '✓' },
    disconnected: { bg: 'bg-gray-100', text: 'text-gray-800', label: 'Disconnected', icon: '○' },
    error: { bg: 'bg-red-100', text: 'text-red-800', label: 'Error', icon: '✗' },
    pending: { bg: 'bg-yellow-100', text: 'text-yellow-800', label: 'Pending', icon: '⋯' },
    normal: { bg: 'bg-green-100', text: 'text-green-800', label: 'Normal', icon: '✓' },
    abnormal: { bg: 'bg-yellow-100', text: 'text-yellow-800', label: 'Abnormal', icon: '!' },
    critical: { bg: 'bg-red-100', text: 'text-red-800', label: 'Critical', icon: '!!' },
  };

  const config = statusConfig[status];
  
  return (
    <span className={`inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-medium ${config.bg} ${config.text}`}>
      <span>{config.icon}</span>
      <span>{label || config.label}</span>
    </span>
  );
};

export default StatusBadge;
