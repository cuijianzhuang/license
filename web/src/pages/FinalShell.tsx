import React, { useState } from 'react';
import { Typography, Form, Button, Input, Alert, Card, Row, Col, Divider } from 'antd';
import styled from 'styled-components';
import { useTranslation } from 'react-i18next';
import { CopyOutlined, CheckOutlined } from '@ant-design/icons';
import PageHeader from '../components/PageHeader';
import { finalshell } from '../api';

const { Paragraph } = Typography;

const FormWrapper = styled.div`
  max-width: 600px;
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

const CodeLabel = styled.div`
  font-weight: 600;
  margin-bottom: 4px;
  color: #4b5563;
`;

const FinalShell: React.FC = () => {
  const { t } = useTranslation();
  const [loading, setLoading] = useState(false);
  const [authorizationCodes, setAuthorizationCodes] = useState<string[]>([]);
  const [copying, setCopying] = useState<{[key: string]: boolean}>({});
  const [form] = Form.useForm();

  const handleGenerateLicense = async (values: { machineCode: string }) => {
    setLoading(true);
    try {
      const data = await finalshell.generateLicense(values.machineCode);
      setAuthorizationCodes(data);
    } catch (error) {
      console.error('生成许可证失败:', error);
    } finally {
      setLoading(false);
    }
  };

  // 复制到剪贴板
  const copyToClipboard = (key: string, text: string) => {
    navigator.clipboard.writeText(text).then(() => {
      setCopying({ ...copying, [key]: true });
      
      setTimeout(() => {
        setCopying({ ...copying, [key]: false });
      }, 2000);
    });
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

  // 解析授权码字符串
  const parseAuthCode = (codeWithLabel: string): { label: string, code: string } => {
    const matches = codeWithLabel.match(/(.*?):\s*(.*)/);
    if (matches && matches.length > 2) {
      return { label: matches[1], code: matches[2] };
    }
    return { label: '', code: codeWithLabel };
  };

  // 获取翻译的版本标签
  const getVersionLabel = (label: string): string => {
    if (label.includes('< 3.9.6') && label.includes('高级版')) {
      return t('finalshell.versions.advancedBelow396');
    } else if (label.includes('< 3.9.6') && label.includes('专业版')) {
      return t('finalshell.versions.proBelow396');
    } else if (label.includes('>= 3.9.6') && label.includes('高级版')) {
      return t('finalshell.versions.advancedAbove396');
    } else if (label.includes('>= 3.9.6') && label.includes('专业版')) {
      return t('finalshell.versions.proAbove396');
    }
    return label;
  };

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

      <Alert
        message={t('finalshell.usageNotice')}
        description={t('finalshell.warningDescription')}
        type="info"
        showIcon
        style={{ marginBottom: 24 }}
      />

      <FormWrapper>
        <Form form={form} onFinish={handleGenerateLicense} layout="vertical">
          <Form.Item
            name="machineCode"
            label={t('finalshell.machineCode')}
            rules={[{ required: true, message: t('finalshell.machineCodeRequired') }]}
          >
            <Input placeholder={t('finalshell.enterMachineCode')} />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>
              {t('finalshell.generateButton')}
            </Button>
          </Form.Item>
        </Form>
      </FormWrapper>

      {authorizationCodes.length > 0 && (
        <FormCard title={t('finalshell.registrationSuccess')}>
          {authorizationCodes.map((codeWithLabel, index) => {
            const { label, code } = parseAuthCode(codeWithLabel);
            const versionLabel = getVersionLabel(label);
            const copyKey = `code-${index}`;

            return (
              <div key={index} style={{marginBottom: 16}}>
                <CodeLabel>{versionLabel}</CodeLabel>
                <AuthorizationCodeContainer>
                  {code}
                  <CopyButton
                    size="small"
                    type="primary"
                    ghost
                    icon={copying[copyKey] ? <CheckOutlined /> : <CopyOutlined />}
                    onClick={() => copyToClipboard(copyKey, code)}
                  />
                </AuthorizationCodeContainer>
              </div>
            );
          })}
        </FormCard>
      )}

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