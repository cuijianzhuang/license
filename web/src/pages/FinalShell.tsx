import React, {useState} from 'react';
import {App, Button, Card, Col, Form, Input, Row, Typography} from 'antd';
import styled from 'styled-components';
import {useTranslation} from 'react-i18next';
import {CheckOutlined, CopyOutlined, InfoCircleOutlined} from '@ant-design/icons';
import PageHeader from '../components/PageHeader';
import UsageNotice from '../components/UsageNotice';
import {finalshell} from '../api';
import {copyAndManageState} from '../utils/clipboard';
import {FinalShellLicense} from '../api/finalshell';

const { Paragraph } = Typography;

const FormWrapper = styled.div`
  width: 100%;
  margin-bottom: 32px;
`;

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

const AuthorizationCodeContainer = styled.div`
  position: relative;
  background-color: #f9fafb;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 16px;
  padding-right: 50px;
  margin-bottom: 0;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  overflow-wrap: break-word;
  word-break: break-all;
  
  @media (max-width: 768px) {
    padding: 12px;
    padding-right: 40px;
    font-size: 12px;
  }
`;

const ButtonContainer = styled.div`
  position: absolute;
  top: 8px;
  right: 8px;
  z-index: 10;
  
  @media (max-width: 768px) {
    top: 4px;
    right: 4px;
  }
`;

const CopyButton = styled(Button)`
  opacity: 0.8;
  
  &:hover {
    opacity: 1;
  }
  
  @media (max-width: 768px) {
    padding: 0 8px;
    height: 24px;
    font-size: 12px;
  }
`;

const CodeContainer = styled.div`
  position: relative;
  margin-bottom: 16px;
  overflow: hidden;
  isolation: isolate;
`;

const CodeLabel = styled.div`
  font-weight: 600;
  margin-bottom: 4px;
  color: #4b5563;
`;

const HostRulesContainer = styled.div`
  background-color: #f9fafb;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 16px;
  padding-right: 50px;
  margin-bottom: 0;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 14px;
  line-height: 1.6;
  white-space: pre-line;
  position: relative;
  
  @media (max-width: 768px) {
    padding: 12px;
    padding-right: 40px;
    font-size: 12px;
  }
`;

const HostPathNotice = styled.div`
  background-color: #e6f7ff;
  border: 1px solid #91d5ff;
  border-radius: 6px;
  padding: 15px;
  margin-bottom: 16px;
  display: flex;
  align-items: flex-start;
  width: 100%;
  
  .info-icon {
    color: #1890ff;
    margin-right: 12px;
    margin-top: 1px;
    flex-shrink: 0;
    font-size: 14px;
  }
  
  .info-content {
    flex: 1;
    
    .info-title {
      font-weight: 500;
      color: rgba(0, 0, 0, 0.85);
      margin-bottom: 8px;
      font-size: 14px;
      line-height: 1.5715;
    }
    
    .path-item {
      font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
      color: rgba(0, 0, 0, 0.65);
      margin-bottom: 4px;
      font-size: 13px;
      line-height: 1.5715;
      
      &:last-child {
        margin-bottom: 0;
      }
      
      .platform {
        font-weight: 600;
        color: rgba(0, 0, 0, 0.85);
      }
    }
  }
  
  @media (max-width: 768px) {
    padding: 12px;
    
    .info-icon {
      margin-right: 8px;
    }
    
    .info-content {
      .info-title {
        font-size: 13px;
      }
      
      .path-item {
        font-size: 12px;
      }
    }
  }
`;

const FinalShell: React.FC = () => {
  const { t } = useTranslation();
  const { notification } = App.useApp();
  const [loading, setLoading] = useState(false);
  const [licenseData, setLicenseData] = useState<FinalShellLicense | null>(null);
  const [copying, setCopying] = useState<{[key: string]: boolean}>({});
  const [form] = Form.useForm();

  // Map of license keys to their translation keys
  const versionMap = {
    'advancedBelow396': 'finalshell.versions.advancedBelow396',
    'proBelow396': 'finalshell.versions.proBelow396',
    'advancedAbove396': 'finalshell.versions.advancedAbove396',
    'proAbove396': 'finalshell.versions.proAbove396',
    'advancedAbove45': 'finalshell.versions.advancedAbove45',
    'proAbove45': 'finalshell.versions.proAbove45',
    'advancedAbove46': 'finalshell.versions.advancedAbove46',
    'proAbove46': 'finalshell.versions.proAbove46',
  };

  // Host rules for blocking network verification
  const hostRules = `127.0.0.1    www.youtusoft.com
127.0.0.1    youtusoft.com
127.0.0.1    hostbuf.com
127.0.0.1    www.hostbuf.com
127.0.0.1    dkys.org
127.0.0.1    tcpspeed.com
127.0.0.1    www.wn1998.com
127.0.0.1    wn1998.com
127.0.0.1    pwlt.wn1998.com
127.0.0.1    backup.www.hostbuf.com`;

  const handleGenerateLicense = async (values: { machineCode: string }) => {
    setLoading(true);
    try {
      const data = await finalshell.generateLicense(values.machineCode);
      setLicenseData(data);
    } catch (error) {
      console.error('Generate license failed:', error);
    } finally {
      setLoading(false);
    }
  };

  // Copy to clipboard
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

  const breadcrumbs = [
    {
      path: '/',
      breadcrumbName: t('nav.home'),
    },
    {
      path: '',
      breadcrumbName: t('nav.finalshell'),
    },
  ];

  return (
    <div>
      <PageHeader
        title={t('finalshell.title')}
        subTitle={t('finalshell.subTitle')}
        breadcrumbs={breadcrumbs}
      />

      <Paragraph>
        {t('finalshell.description')}
      </Paragraph>

      <UsageNotice
        message={t('finalshell.usageNotice')}
        description={t('finalshell.warningDescription')}
      />

      <FormWrapper>
        <Form form={form} onFinish={handleGenerateLicense} layout="vertical">
          <Row gutter={16}>
            <Col xs={24} sm={12}>
              <Form.Item
                name="machineCode"
                label={t('finalshell.machineCode')}
                rules={[{ required: true, message: t('finalshell.machineCodeRequired') }]}
              >
                <Input placeholder={t('finalshell.enterMachineCode')} />
              </Form.Item>
            </Col>
          </Row>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>
              {t('finalshell.generateButton')}
            </Button>
          </Form.Item>
        </Form>
      </FormWrapper>

      {licenseData && (
        <FormCard title={t('finalshell.registrationSuccess')}>
          <Row gutter={16}>
            {Object.entries(licenseData).map(([key, code], index) => {
              const versionLabelKey = versionMap[key as keyof typeof versionMap] || '';
              const copyKey = `code-${key}`;

              return (
                <Col key={key} xs={24} sm={12} md={12} style={{marginBottom: 16}}>
                  <CodeLabel>{t(versionLabelKey)}</CodeLabel>
                  <CodeContainer>
                    <AuthorizationCodeContainer>
                      {code}
                      <ButtonContainer>
                        <CopyButton
                          size="small"
                          type="primary"
                          ghost
                          icon={copying[copyKey] ? <CheckOutlined /> : <CopyOutlined />}
                          onClick={() => copyToClipboard(copyKey, code)}
                        />
                      </ButtonContainer>
                    </AuthorizationCodeContainer>
                  </CodeContainer>
                </Col>
              );
            })}
          </Row>
        </FormCard>
      )}

      <FormCard title={t('finalshell.hostBlockTitle')}>
        <Paragraph style={{ marginBottom: 16, color: '#4b5563' }}>
          {t('finalshell.hostBlockDescription')}
        </Paragraph>
        
        <HostPathNotice>
          <InfoCircleOutlined className="info-icon" />
          <div className="info-content">
            <div className="info-title">{t('finalshell.hostFilePath')}</div>
            <div className="path-item">
              <span className="platform">Windows:</span> C:\Windows\System32\drivers\etc\hosts
            </div>
            <div className="path-item">
              <span className="platform">macOS/Linux:</span> /etc/hosts
            </div>
          </div>
        </HostPathNotice>
        
        <CodeLabel>{t('finalshell.hostRules')}</CodeLabel>
        <CodeContainer>
          <HostRulesContainer>
            {hostRules}
            <ButtonContainer>
              <CopyButton
                size="small"
                type="primary"
                ghost
                icon={copying['hostRules'] ? <CheckOutlined /> : <CopyOutlined />}
                onClick={() => copyToClipboard('hostRules', hostRules)}
              />
            </ButtonContainer>
          </HostRulesContainer>
        </CodeContainer>
      </FormCard>

      <FormCard title={t('finalshell.instructionsTitle')}>
        <StepItem>
          <StepNumber>1</StepNumber>
          <StepContent>{t('finalshell.usageSteps.step1')}</StepContent>
        </StepItem>
        <StepItem>
          <StepNumber>2</StepNumber>
          <StepContent>{t('finalshell.usageSteps.step2')}</StepContent>
        </StepItem>
        <StepItem>
          <StepNumber>3</StepNumber>
          <StepContent>{t('finalshell.usageSteps.step3')}</StepContent>
        </StepItem>
        <StepItem>
          <StepNumber>4</StepNumber>
          <StepContent>{t('finalshell.usageSteps.step4')}</StepContent>
        </StepItem>
        <StepItem>
          <StepNumber>5</StepNumber>
          <StepContent>{t('finalshell.usageSteps.step5')}</StepContent>
        </StepItem>
        <StepItem>
          <StepNumber>6</StepNumber>
          <StepContent>{t('finalshell.usageSteps.step6')}</StepContent>
        </StepItem>
        <StepItem>
          <StepNumber>7</StepNumber>
          <StepContent>{t('finalshell.usageSteps.step7')}</StepContent>
        </StepItem>
      </FormCard>
    </div>
  );
};

export default FinalShell; 