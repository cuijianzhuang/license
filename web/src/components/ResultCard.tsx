import React, {useState} from 'react';
import styled from 'styled-components';
import {App, Button, Card, Divider, message, Space, Tooltip, Typography} from 'antd';
import {CheckOutlined, CopyOutlined, DownloadOutlined} from '@ant-design/icons';
import {useTranslation} from "react-i18next";
import {copyAndManageState} from '../utils/clipboard';

const { Title, Text } = Typography;

interface ResultCardProps {
  title: string;
  data: Record<string, string>;
  fileName?: string;
}

const StyledCard = styled(Card)`
  margin-bottom: 32px;
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
  border: 1px solid #e5e7eb;
  overflow: hidden;
  transition: box-shadow 0.3s;

  &:hover {
    box-shadow: 0 6px 16px rgba(0, 0, 0, 0.08);
  }
  
  .ant-card-head {
    background-color: #f9fafb;
    border-bottom: 1px solid #e5e7eb;
    min-height: 56px;
  }
  
  .ant-card-head-title {
    padding: 16px 0;
  }
`;

const LicenseContent = styled.div`
  margin-top: 12px;
  background-color: #f9fafb;
  padding: 16px;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
  white-space: pre-wrap;
  word-break: break-all;
  position: relative;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 14px;
  overflow-x: auto;
`;

const CardTitle = styled(Title)`
  font-size: 16px !important;
  margin-bottom: 0 !important;
  font-weight: 600 !important;
`;

const LabelText = styled(Text)`
  font-weight: 500;
  color: #4b5563;
`;

const ButtonContainer = styled.div`
  position: absolute;
  top: 8px;
  right: 8px;
  z-index: 2;
`;

const CopyButton = styled(Button)`
  opacity: 0.8;
  
  &:hover {
    opacity: 1;
  }
`;

const DownloadButton = styled(Button)`
  border-radius: 6px;
`;

const ResultCard: React.FC<ResultCardProps> = ({ title, data, fileName }) => {
  const { t } = useTranslation();
  const { notification } = App.useApp();
  const [copying, setCopying] = useState<Record<string, boolean>>({});

  const copyToClipboard = (key: string, text: string) => {
    copyAndManageState(
      key,
      text,
      copying,
      setCopying,
      notification,
      t('common.copied'),
      t('common.copyFail')
    );
  };

  const downloadAsFile = () => {
    if (!fileName) return;
    
    const dataStr = Object.entries(data)
      .map(([key, value]) => `${key}: ${value}`)
      .join('\n');
    
    const blob = new Blob([dataStr], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = fileName;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
    
    message.success('下载成功');
  };

  return (
    <StyledCard
      title={<CardTitle>{title}</CardTitle>}
      extra={
        fileName && (
          <Tooltip title="下载为文件">
            <DownloadButton 
              type="primary"
              ghost
              icon={<DownloadOutlined />} 
              onClick={downloadAsFile}
            >
              下载
            </DownloadButton>
          </Tooltip>
        )
      }
    >
      {Object.entries(data).map(([key, value], index) => (
        <div key={key}>
          {index > 0 && <Divider style={{ margin: '16px 0' }} />}
          <Space direction="vertical" style={{ width: '100%' }}>
            <LabelText>{key}:</LabelText>
            <LicenseContent>
              {value}
              <ButtonContainer>
                <CopyButton
                  size="small"
                  type="primary"
                  ghost
                  icon={copying[key] ? <CheckOutlined /> : <CopyOutlined />}
                  onClick={() => copyToClipboard(key, value)}
                />
              </ButtonContainer>
            </LicenseContent>
          </Space>
        </div>
      ))}
    </StyledCard>
  );
};

export default ResultCard; 