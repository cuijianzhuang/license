import React, { useState } from 'react';
import { Typography, Form, Button, Input, DatePicker, Alert, message, Card } from 'antd';
import styled from 'styled-components';
import dayjs from 'dayjs';
import type { Dayjs } from 'dayjs';
import { useTranslation } from 'react-i18next';
import PageHeader from '../components/PageHeader';
import { gitlab } from '../api';

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

const GitLab: React.FC = () => {
  const { t } = useTranslation();
  const [loading, setLoading] = useState(false);
  const [form] = Form.useForm();
  
  // 计算本年度最后一天的23:59:59
  const getEndOfYear = () => {
    const currentYear = dayjs().year();
    return dayjs(`${currentYear}-12-31 23:59:59`);
  };

  const handleGenerateLicense = async (values: {
    name: string;
    email: string;
    company: string;
    expireTime: Dayjs;
  }) => {
    setLoading(true);
    try {
      const success = await gitlab.generateLicense(
        values.name,
        values.email,
        values.company,
        values.expireTime.format('YYYY-MM-DD HH:mm:ss')
      );
      
      if (success) {
        // 显示成功消息
        message.success(t('gitlab.success.downloadStarted'));
      } else {
        // 虽然不应该到达这里，但以防万一
        message.warning(t('gitlab.success.downloadWarning'));
      }
    } catch (error: any) {
      console.error('生成许可证失败:', error);
      // 显示更具体的错误消息
      const errorMsg = error.message || t('gitlab.success.downloadFailed');
      message.error(errorMsg);
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
      breadcrumbName: t('nav.gitlab'),
    },
  ];

  return (
    <div>
      <PageHeader
        title={t('gitlab.title')}
        subTitle={t('gitlab.subTitle')}
        breadcrumbs={breadcrumbs}
      />

      <Paragraph>
        {t('gitlab.description')}
      </Paragraph>

      <Alert
        message={t('gitlab.usageNotice')}
        description={t('gitlab.warningDescription')}
        type="info"
        showIcon
        style={{ marginBottom: 24 }}
      />

      <FormWrapper>
        <Form form={form} onFinish={handleGenerateLicense} layout="vertical">
          <Form.Item
            name="name"
            label={t('gitlab.form.name')}
            rules={[{ required: true, message: t('gitlab.form.namePlaceholder') }]}
          >
            <Input placeholder={t('gitlab.form.namePlaceholder')} />
          </Form.Item>

          <Form.Item
            name="email"
            label={t('gitlab.form.email')}
            rules={[
              { required: true, message: t('gitlab.form.emailPlaceholder') },
              { type: 'email', message: t('gitlab.form.emailInvalid') },
            ]}
          >
            <Input placeholder={t('gitlab.form.emailPlaceholder')} />
          </Form.Item>

          <Form.Item
            name="company"
            label={t('gitlab.form.company')}
            rules={[{ required: true, message: t('gitlab.form.companyPlaceholder') }]}
          >
            <Input placeholder={t('gitlab.form.companyPlaceholder')} />
          </Form.Item>

          <Form.Item
            name="expireTime"
            label={t('gitlab.form.expireTime')}
            rules={[{ required: true, message: t('gitlab.form.expireTimePlaceholder') }]}
            initialValue={getEndOfYear()}
          >
            <DatePicker
              format="YYYY-MM-DD HH:mm:ss"
              style={{ width: '100%' }}
              placeholder={t('gitlab.form.expireTimePlaceholder')}
              showTime={{ defaultValue: dayjs('23:59:59', 'HH:mm:ss') }}
              disabledDate={(current) => {
                return current && current < dayjs().endOf('day');
              }}
            />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>
              {t('gitlab.form.generateButton')}
            </Button>
          </Form.Item>
        </Form>
      </FormWrapper>

      <FormCard title={t('gitlab.instructionsTitle')}>
        <StepItem>
          <StepNumber>1</StepNumber>
          <StepContent>{t('gitlab.usageSteps.step1')}</StepContent>
        </StepItem>
        <StepItem>
          <StepNumber>2</StepNumber>
          <StepContent>{t('gitlab.usageSteps.step2')}</StepContent>
        </StepItem>
        <StepItem>
          <StepNumber>3</StepNumber>
          <StepContent>{t('gitlab.usageSteps.step3')}</StepContent>
        </StepItem>
        <StepItem>
          <StepNumber>4</StepNumber>
          <StepContent>{t('gitlab.usageSteps.step4')}</StepContent>
        </StepItem>
        <StepItem>
          <StepNumber>5</StepNumber>
          <StepContent>{t('gitlab.usageSteps.step5')}</StepContent>
        </StepItem>
        <StepItem>
          <StepNumber>6</StepNumber>
          <StepContent>{t('gitlab.usageSteps.step6')}</StepContent>
        </StepItem>
        <StepItem>
          <StepNumber>7</StepNumber>
          <StepContent>{t('gitlab.usageSteps.step7')}</StepContent>
        </StepItem>
        <StepItem>
          <StepNumber>8</StepNumber>
          <StepContent>{t('gitlab.usageSteps.step8')}</StepContent>
        </StepItem>
      </FormCard>
    </div>
  );
};

export default GitLab; 