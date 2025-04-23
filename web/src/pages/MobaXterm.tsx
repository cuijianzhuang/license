import React, { useState } from 'react';
import { Typography, Form, Button, Input, Select, Alert, message, Card } from 'antd';
import styled from 'styled-components';
import { useTranslation } from 'react-i18next';
import PageHeader from '../components/PageHeader';
import { mobaxterm } from '../api';

const { Paragraph } = Typography;
const { Option } = Select;

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

const versions = [
  '23.6',
  '23.5',
  '23.4',
  '23.3',
  '23.2',
  '23.1',
  '23.0',
  '22.3',
  '22.2',
  '22.1',
  '22.0',
  '21.5',
  '21.4',
  '21.3',
  '21.2',
  '21.1',
  '21.0',
  '20.6',
  '20.5',
  '20.4',
  '20.3',
  '20.2',
  '20.1',
  '20.0',
];

const MobaXterm: React.FC = () => {
  const { t } = useTranslation();
  const [loading, setLoading] = useState(false);
  const [form] = Form.useForm();

  const handleGenerateLicense = async (values: { 
    username: string; 
    version: string;
    count: string;
  }) => {
    setLoading(true);
    try {
      // Create FormData
      const formData = new FormData();
      formData.append('name', values.username);
      formData.append('version', values.version);
      formData.append('count', values.count);

      // Get file blob response
      const blob = await mobaxterm.generateLicense(formData);
      
      // Create download
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = 'mobaxterm-license.txt';
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      URL.revokeObjectURL(url);
      
      message.success(t('mobaxterm.success.downloadStarted'));
    } catch (error) {
      console.error('Failed to generate license:', error);
      message.error(t('common.error'));
    } finally {
      setLoading(false);
    }
  };

  const breadcrumbs = [
    {
      path: '/',
      breadcrumbName: t('nav.home'),
    },
    {
      path: '',
      breadcrumbName: t('nav.mobaxterm'),
    },
  ];

  return (
    <div>
      <PageHeader
        title={t('mobaxterm.title')}
        subTitle={t('mobaxterm.subTitle')}
        breadcrumbs={breadcrumbs}
      />

      <Paragraph>
        {t('mobaxterm.description')}
      </Paragraph>

      <Alert
        message={t('mobaxterm.warning')}
        description={t('mobaxterm.warningDescription')}
        type="warning"
        showIcon
        style={{ marginBottom: 24 }}
      />

      <FormWrapper>
        <Form form={form} onFinish={handleGenerateLicense} layout="vertical">
          <Form.Item
            name="username"
            label={t('mobaxterm.form.username')}
            rules={[{ required: true, message: t('mobaxterm.form.usernamePlaceholder') }]}
            initialValue="MobaXterm User"
          >
            <Input placeholder={t('mobaxterm.form.usernamePlaceholder')} />
          </Form.Item>

          <Form.Item
            name="version"
            label={t('mobaxterm.form.version')}
            rules={[{ required: true, message: t('mobaxterm.form.versionPlaceholder') }]}
            initialValue={versions[0]}
          >
            <Select placeholder={t('mobaxterm.form.versionPlaceholder')}>
              {versions.map((version) => (
                <Option key={version} value={version}>
                  {version}
                </Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item
            name="count"
            label={t('mobaxterm.form.count')}
            rules={[{ required: true, message: t('mobaxterm.form.countPlaceholder') }]}
            initialValue="1000"
          >
            <Input placeholder={t('mobaxterm.form.countPlaceholder')} />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>
              {t('mobaxterm.form.generateButton')}
            </Button>
          </Form.Item>
        </Form>
      </FormWrapper>

      <FormCard title={t('mobaxterm.instructionsTitle')}>
        <StepItem>
          <StepNumber>1</StepNumber>
          <StepContent>{t('mobaxterm.usageSteps.step1')}</StepContent>
        </StepItem>
        <StepItem>
          <StepNumber>2</StepNumber>
          <StepContent>{t('mobaxterm.usageSteps.step2')}</StepContent>
        </StepItem>
        <StepItem>
          <StepNumber>3</StepNumber>
          <StepContent>{t('mobaxterm.usageSteps.step3')}</StepContent>
        </StepItem>
        <StepItem>
          <StepNumber>4</StepNumber>
          <StepContent>{t('mobaxterm.usageSteps.step4')}</StepContent>
        </StepItem>
      </FormCard>
    </div>
  );
};

export default MobaXterm; 