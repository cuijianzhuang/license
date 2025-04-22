import React, { useState } from 'react';
import { Typography, Form, Button, Input, DatePicker, Alert, message } from 'antd';
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
        message={t('gitlab.warning')}
        description={t('gitlab.warningDescription')}
        type="warning"
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
    </div>
  );
};

export default GitLab; 