import React, { useState, useEffect, useRef } from 'react';
import { Typography, Alert, Card, Button, Tooltip } from 'antd';
import styled from 'styled-components';
import { useTranslation } from 'react-i18next';
import { 
  CopyOutlined, 
  CheckOutlined, 
  ReloadOutlined 
} from '@ant-design/icons';
import PageHeader from '../components/PageHeader';
import { jrebel } from '../api';

const { Title, Paragraph } = Typography;

const FormCard = styled(Card)`
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  margin-bottom: 32px;
  border: 1px solid #e5e7eb;
  
  .ant-card-head {
    border-bottom: 1px solid #e5e7eb;
  }
`;

const StepItem = styled.div`
  margin-bottom: 16px;
  display: flex;
  align-items: flex-start;
`;

const StepNumber = styled.span`
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 24px;
  height: 24px;
  background-color: #1890ff;
  color: #fff;
  border-radius: 50%;
  margin-right: 12px;
  font-size: 14px;
  flex-shrink: 0;
`;

const StepContent = styled.div`
  flex: 1;
`;

const ServerAddressContainer = styled.div`
  position: relative;
  background-color: #f9fafb;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 16px;
  margin-bottom: 16px;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  overflow-wrap: break-word;
  word-break: break-all;
`;

const CopyButton = styled(Button)`
  position: absolute;
  top: 8px;
  right: 8px;
  opacity: 0.8;
  z-index: 2;
  
  &:hover {
    opacity: 1;
  }
`;

const ReloadButton = styled(Button)`
  position: absolute;
  top: 8px;
  right: 42px;
  opacity: 0.8;
  z-index: 2;
  
  &:hover {
    opacity: 1;
  }
`;

const JRebel: React.FC = () => {
  const { t } = useTranslation();
  const [serverAddress, setServerAddress] = useState<string>('');
  const [jrebelAddress, setJrebelAddress] = useState<string>('');
  const [guid, setGuid] = useState<string>('');
  const [copying, setCopying] = useState<{[key: string]: boolean}>({});
  const [regenerating, setRegenerating] = useState(false);
  
  // 生成GUID和设置地址
  const generateAndSetAddresses = () => {
    const protocol = window.location.protocol;
    const host = window.location.host;
    const baseAddress = `${protocol}//${host}`;
    setServerAddress(baseAddress);
    
    // 生成一个随机的GUID
    const newGuid = jrebel.generateGuid();
    setGuid(newGuid);
    
    // 设置JRebel完整地址
    setJrebelAddress(`${baseAddress}/${newGuid}`);
  };
  
  // 获取浏览器主机地址和生成GUID
  useEffect(() => {
    generateAndSetAddresses();
  }, []);

  // 复制到剪贴板
  const copyToClipboard = (key: string, text: string) => {
    navigator.clipboard.writeText(text).then(() => {
      setCopying({ ...copying, [key]: true });
      
      setTimeout(() => {
        setCopying({ ...copying, [key]: false });
      }, 2000);
    });
  };
  
  // 重新生成GUID
  const handleRegenerateGuid = () => {
    setRegenerating(true);
    
    // 生成新的GUID并更新地址
    generateAndSetAddresses();
    
    // 显示重新生成动画
    setTimeout(() => {
      setRegenerating(false);
    }, 500);
  };

  const breadcrumbs = [
    {
      path: '/',
      breadcrumbName: t('nav.home'),
    },
    {
      path: '',
      breadcrumbName: t('nav.jrebel'),
    },
  ];

  return (
    <div>
      <PageHeader
        title={t('jrebel.title')}
        subTitle={t('jrebel.subTitle')}
        breadcrumbs={breadcrumbs}
      />

      <Paragraph>
        {t('jrebel.description')}
      </Paragraph>

      <Alert
        message={t('jrebel.usageNotice')}
        description={t('jrebel.authorizationDescription')}
        type="info"
        showIcon
        style={{ marginBottom: 24 }}
      />

      <FormCard title={t('jrebel.serverConfig')}>
        <Paragraph style={{ marginBottom: 8 }}>
          {t('jrebel.baseServerAddress')}:
        </Paragraph>
        <ServerAddressContainer>
          {jrebelAddress}
          <CopyButton
            size="small"
            type="primary"
            ghost
            icon={copying['jrebelAddress'] ? <CheckOutlined /> : <CopyOutlined />}
            onClick={() => copyToClipboard('jrebelAddress', jrebelAddress)}
          />
          <ReloadButton
            size="small"
            type="primary"
            ghost
            icon={<ReloadOutlined spin={regenerating} />}
            onClick={handleRegenerateGuid}
            title={t('jrebel.regenerateGuid')}
          />
        </ServerAddressContainer>
      </FormCard>

      <FormCard title={t('jrebel.usageSteps')}>
        <StepItem>
          <StepNumber>1</StepNumber>
          <StepContent>{t('jrebel.step1')}</StepContent>
        </StepItem>
        <StepItem>
          <StepNumber>2</StepNumber>
          <StepContent>{t('jrebel.step2')}</StepContent>
        </StepItem>
        <StepItem>
          <StepNumber>3</StepNumber>
          <StepContent>{t('jrebel.step3')}</StepContent>
        </StepItem>
        <StepItem>
          <StepNumber>4</StepNumber>
          <StepContent>{t('jrebel.step4')}</StepContent>
        </StepItem>
        <StepItem>
          <StepNumber>5</StepNumber>
          <StepContent>{t('jrebel.step5')}</StepContent>
        </StepItem>
        <StepItem>
          <StepNumber>6</StepNumber>
          <StepContent>{t('jrebel.step6')}</StepContent>
        </StepItem>
      </FormCard>
    </div>
  );
};

export default JRebel; 