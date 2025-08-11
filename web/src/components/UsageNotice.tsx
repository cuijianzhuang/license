import React from 'react';
import {Alert} from 'antd';
import styled from 'styled-components';

interface UsageNoticeProps {
  message: string;
  description: string;
  type?: 'info' | 'success' | 'warning' | 'error';
  showIcon?: boolean;
  style?: React.CSSProperties;
}

// 限制提示框的最大宽度，避免占满整个屏幕
const StyledAlert = styled(Alert)`
  max-width: 800px;
  margin-bottom: 24px;
  transition: max-width 0.3s;
  
  @media (max-width: 1000px) {
    max-width: 100%;
  }
`;

const UsageNotice: React.FC<UsageNoticeProps> = ({ 
  message, 
  description, 
  type = 'info', 
  showIcon = true, 
  style 
}) => {
  return (
    <StyledAlert
      message={message}
      description={description}
      type={type}
      showIcon={showIcon}
      style={style}
    />
  );
};

export default UsageNotice; 