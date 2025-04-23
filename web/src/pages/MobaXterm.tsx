import React, { useState, useEffect, useRef } from 'react';
import { Typography, Form, Button, Input, Select, Alert, message, Card, Spin } from 'antd';
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

// Fallback versions in case API fails
const fallbackVersions = [
  '25.1',
  '25.0',
  '24.4',
  '24.3',
  '23.6',
  '23.5',
  '23.0',
  '22.3',
  '21.5',
  '21.0',
  '20.6',
  '20.0',
];

const MobaXterm: React.FC = () => {
  const { t } = useTranslation();
  const [loading, setLoading] = useState(false);
  const [versions, setVersions] = useState<string[]>(fallbackVersions);
  const [loadingVersions, setLoadingVersions] = useState(false);
  const fetchedRef = useRef(false);
  const [form] = Form.useForm();

  // Fetch versions from API when component mounts - with improved controls
  useEffect(() => {
    // Set initial form value with default version
    form.setFieldsValue({ 
      version: fallbackVersions[0],
      count: "1000"
    });

    // Only fetch if we haven't already
    if (fetchedRef.current) {
      return;
    }

    const fetchVersionsFromApi = async () => {
      // Prevent concurrent requests
      fetchedRef.current = true;
      setLoadingVersions(true);
      
      try {
        console.log('Fetching versions...');
        const fetchedVersions = await mobaxterm.fetchVersions();
        console.log('Fetched versions:', fetchedVersions);
        
        if (fetchedVersions && Array.isArray(fetchedVersions) && fetchedVersions.length > 0) {
          // 确保state更新
          setVersions([...fetchedVersions]);
          // 需要确保版本值已经被设置，并保留已有的用户名和数量
          setTimeout(() => {
            const currentValues = form.getFieldsValue();
            form.setFieldsValue({ 
              version: fetchedVersions[0],
              count: currentValues.count || "1000"
            });
            console.log('Updated form with version:', fetchedVersions[0]);
          }, 0);
        } else {
          console.log('Using fallback versions - invalid response');
          const currentValues = form.getFieldsValue();
          form.setFieldsValue({ 
            version: fallbackVersions[0],
            count: currentValues.count || "1000"
          });
        }
      } catch (error) {
        console.error('Failed to fetch versions:', error);
        // No need to show error message, just use fallbacks
      } finally {
        setLoadingVersions(false);
      }
    };

    fetchVersionsFromApi();
  }, [form, t]);

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
      a.download = 'Custom.mxtpro';
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
        message={t('mobaxterm.usageNotice')}
        description={t('mobaxterm.warningDescription')}
        type="info"
        showIcon
        style={{ marginBottom: 24 }}
      />

      <FormWrapper>
        <Form form={form} onFinish={handleGenerateLicense} layout="vertical">
          <Form.Item
            name="username"
            label={t('mobaxterm.form.username')}
            rules={[{ required: true, message: t('mobaxterm.form.usernamePlaceholder') }]}
          >
            <Input placeholder={t('mobaxterm.form.usernamePlaceholder')} />
          </Form.Item>

          <Form.Item
            name="version"
            label={t('mobaxterm.form.version')}
            rules={[{ required: true, message: t('mobaxterm.form.versionPlaceholder') }]}
            initialValue={fallbackVersions[0]}
          >
            <Select 
              placeholder={t('mobaxterm.form.versionPlaceholder')}
              loading={loadingVersions}
              notFoundContent={loadingVersions ? <Spin size="small" /> : null}
              showSearch
              optionFilterProp="children"
              onDropdownVisibleChange={(open) => {
                if (open) {
                  console.log('Dropdown opened, available versions:', versions);
                }
              }}
            >
              {versions.map((version) => (
                <Option key={version} value={version}>
                  {version}
                </Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item
            label={t('mobaxterm.form.count')}
            name="count"
            initialValue="1000"
            rules={[
              {
                required: true,
                message: t('mobaxterm.form.countPlaceholder')
              },
              {
                pattern: /^[1-9]\d*$/,
                message: t('mobaxterm.form.countInvalid')
              }
            ]}
          >
            <Input 
              placeholder={t('mobaxterm.form.countPlaceholder')} 
              type="number"
              min={1}
              step={1}
              onKeyDown={(e) => {
                // Prevent typing e, +, - or decimal point
                if (['+', '-', 'e', '.'].includes(e.key)) {
                  e.preventDefault();
                }
              }}
              onChange={(e) => {
                // Remove any leading zeros
                if (e.target.value.startsWith('0')) {
                  e.target.value = e.target.value.replace(/^0+/, '');
                }
                // If empty after removing zeros, set to empty
                if (e.target.value === '') {
                  e.target.value = '';
                }
              }}
            />
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
      </FormCard>
    </div>
  );
};

export default MobaXterm; 