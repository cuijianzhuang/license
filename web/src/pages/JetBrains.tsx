import React, { useState, useEffect, useRef } from 'react';
import {Typography, Form, Button, Alert, Card, Space, Input, DatePicker, App} from 'antd';
import styled from 'styled-components';
import dayjs from 'dayjs';
import type { Dayjs } from 'dayjs';
import 'dayjs/locale/zh-cn';
import 'dayjs/locale/en';
import { 
  LoadingOutlined, 
  InfoCircleOutlined,
  CopyOutlined,
  CheckOutlined
} from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import PageHeader from '../components/PageHeader';
import ResultCard from '../components/ResultCard';
import { jetbrains } from '../api';
import { JetBrainsLicense } from '../types';
import { copyAndManageState } from '../utils/clipboard';

const { Paragraph, Text, Title } = Typography;

const FormCard = styled(Card)`
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  margin-bottom: 32px;
  border: 1px solid #e5e7eb;
  
  .ant-card-head {
    border-bottom: 1px solid #e5e7eb;
  }
`;

const SubmitButton = styled(Button)`
  width: 100%;
  height: 40px;
  border-radius: 8px;
  margin-top: 8px;
`;

const LicenseContent = styled.div`
  margin-top: 12px;
  background-color: #f9fafb;
  padding: 16px;
  padding-right: 48px;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
  white-space: pre-wrap;
  word-break: break-all;
  position: relative;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 14px;
  overflow-x: auto;
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

const LicenseResultCard = styled(Card)`
  margin-bottom: 32px;
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
  border: 1px solid #e5e7eb;
  overflow: hidden;

  .ant-card-head {
    background-color: #f9fafb;
    border-bottom: 1px solid #e5e7eb;
  }
`;

const LabelText = styled(Text)`
  font-weight: 500;
  color: #4b5563;
`;

// 创建自定义的切换按钮容器
const TabsContainer = styled.div`
  display: flex;
  margin-bottom: 24px;
  border-bottom: 1px solid #f0f0f0;
`;

// 创建自定义的切换按钮
const TabButton = styled.button<{ active: boolean }>`
  padding: 8px 16px;
  background: transparent;
  border: none;
  cursor: pointer;
  font-size: 14px;
  font-weight: ${props => props.active ? '500' : 'normal'};
  color: ${props => props.active ? '#1890ff' : 'rgba(0, 0, 0, 0.65)'};
  position: relative;
  transition: all 0.3s;
  
  &:focus {
    outline: none;
  }
  
  &::after {
    content: '';
    position: absolute;
    bottom: -1px;
    left: 0;
    right: 0;
    height: 2px;
    background: ${props => props.active ? '#1890ff' : 'transparent'};
  }
  
  &:hover {
    color: #1890ff;
  }
`;

const JetBrains: React.FC = () => {
  const { t, i18n } = useTranslation();
  const { notification } = App.useApp();
  const [loading, setLoading] = useState(false);
  const [license, setLicense] = useState<JetBrainsLicense | null>(null);
  const [rawResponse, setRawResponse] = useState<string | null>(null);
  const [serverRule, setServerRule] = useState<string>('');
  const [loadingServerRule, setLoadingServerRule] = useState(false);
  const [activationMethod, setActivationMethod] = useState<'code' | 'server'>('code');
  const [copying, setCopying] = useState<{[key: string]: boolean}>({});
  const [codeForm] = Form.useForm();
  const [serverAddress, setServerAddress] = useState<string>('');
  const previousMethodRef = useRef<'code' | 'server'>(activationMethod);

  // Set dayjs locale based on i18n language
  useEffect(() => {
    // Map i18n language code to dayjs locale
    const locale = i18n.language.startsWith('zh') ? 'zh-cn' : 'en';
    dayjs.locale(locale);
  }, [i18n.language]);

  // Get browser's host address when component mounts
  useEffect(() => {
    const protocol = window.location.protocol;
    const host = window.location.host;
    setServerAddress(`${protocol}//${host}`);
  }, []);
  
  // Reset form and clear results when activation method changes
  useEffect(() => {
    // Only perform reset if the method actually changed (not on initial render)
    if (previousMethodRef.current !== activationMethod) {
      // Clear license results
      setLicense(null);
      setRawResponse(null);
      
      // Reset the form fields
      codeForm.resetFields();
      
      // Update the previous method ref
      previousMethodRef.current = activationMethod;
    }
  }, [activationMethod, codeForm]);
  
  // Fetch server rules when user switches to server authorization mode
  useEffect(() => {
    if (activationMethod === 'server' && !serverRule && !loadingServerRule) {
      setLoadingServerRule(true);
      
      const fetchServerRule = async () => {
        try {
          const serverRuleText = await jetbrains.getLicenseServerRule();
          setServerRule(serverRuleText);
        } catch (error) {
          console.error(`${t('jetbrains.serverRuleFetchError')}:`, error);
        } finally {
          setLoadingServerRule(false);
        }
      };
  
      fetchServerRule();
    }
  }, [activationMethod, serverRule, loadingServerRule, t]);

  const handleTabChange = (method: 'code' | 'server') => {
    setActivationMethod(method);
  };

  const handleGenerateLicense = async (values: { 
    licenseeName?: string, 
    effectiveDate?: Dayjs,
    codes?: string 
  }) => {
    setLoading(true);
    try {
      // Format the date if it exists
      const formattedDate = values.effectiveDate 
        ? values.effectiveDate.format('YYYY-MM-DD HH:mm:ss') 
        : undefined;
        
      const data = await jetbrains.generateLicense(
        values.licenseeName, 
        formattedDate, 
        values.codes
      );

      // Store the raw response data for our custom rendering
      setRawResponse(typeof data === 'string' ? data : String(data));
      
      // Also keep the license object for compatibility
      setLicense({
        code: '',
        product: values.codes?.split(',')[0] || t('jetbrains.unknownProduct')
      });
    } catch (error) {
      console.error(`${t('jetbrains.licenseGenerationError')}:`, error);
    } finally {
      setLoading(false);
    }
  };

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

  // Extract power.conf content from raw response
  const extractPowerConf = (): string | null => {
    if (!rawResponse) return null;
    
    // Find the power.conf section
    const startMarker = "================== power.conf ==================";
    const endMarker = "================== power.conf ==================";
    
    const startIdx = rawResponse.indexOf(startMarker);
    if (startIdx === -1) return null;
    
    const contentStart = startIdx + startMarker.length;
    const endIdx = rawResponse.indexOf(endMarker, contentStart);
    if (endIdx === -1) return null;
    
    return rawResponse.substring(contentStart, endIdx).trim();
  };

  // Extract activation code from raw response
  const extractActivationCode = (): string | null => {
    if (!rawResponse) return null;
    
    // Find the activation code section
    const startMarker = "================== activation code ==================";
    const endMarker = "================== activation code ==================";
    
    const startIdx = rawResponse.indexOf(startMarker);
    if (startIdx === -1) return null;
    
    const contentStart = startIdx + startMarker.length;
    const endIdx = rawResponse.indexOf(endMarker, contentStart);
    if (endIdx === -1) return null;
    
    return rawResponse.substring(contentStart, endIdx).trim();
  };

  const breadcrumbs = [
    {
      path: '/',
      breadcrumbName: t('nav.home'),
    },
    {
      path: '',
      breadcrumbName: t('nav.jetbrains'),
    },
  ];

  const onFinish = (values: any) => {
    if (!values.licenseeName) {
      notification.error({
        message: t('jetbrains.pleaseEnterLicenseeName'),
        duration: 3,
      });
      return;
    }

    // If generating authorization code by product
    if (activationMethod === 'code') {
      handleGenerateLicense({
        licenseeName: values.licenseeName,
        effectiveDate: values.effectiveDate,
        codes: values.manualCodes
      });
    } else if (activationMethod === 'server') {
      // Server authorization method, no need to pass product codes
      handleGenerateLicense({
        licenseeName: values.licenseeName,
        effectiveDate: values.effectiveDate
      });
    }
  };

  // Function to disable past dates in DatePicker
  const disablePastDates = (current: Dayjs) => {
    // Can not select days before today
    return current && current < dayjs().startOf('day');
  };

  return (
    <div>
      <PageHeader
        title={t('jetbrains.title')}
        subTitle={t('jetbrains.subTitle')}
        breadcrumbs={breadcrumbs}
      />

      <FormCard title={t('jetbrains.title')}>
        <TabsContainer>
          <TabButton 
            active={activationMethod === 'code'} 
            onClick={() => handleTabChange('code')}
          >
            {t('jetbrains.codeActivation')}
          </TabButton>
          <TabButton 
            active={activationMethod === 'server'} 
            onClick={() => handleTabChange('server')}
          >
            {t('jetbrains.serverActivation')}
          </TabButton>
        </TabsContainer>

        {activationMethod === 'code' ? (
          <Form form={codeForm} onFinish={onFinish} layout="vertical" preserve={false}>
            <Form.Item
              name="licenseeName"
              label={t('jetbrains.licenseeName')}
              rules={[{ required: true, message: t('jetbrains.pleaseEnterLicenseeName') }]}
            >
              <Input placeholder={t('jetbrains.pleaseEnterLicenseeName')} />
            </Form.Item>

            <Form.Item
              name="effectiveDate"
              label={t('jetbrains.effectiveDate')}
            >
              <DatePicker 
                showTime
                style={{ width: '100%' }} 
                placeholder={t('jetbrains.effectiveDatePlaceholder')}
                format="YYYY-MM-DD HH:mm:ss"
                disabledDate={disablePastDates}
                showNow={false}
                locale={i18n.language.startsWith('zh') ? 
                  require('antd/es/date-picker/locale/zh_CN').default : 
                  require('antd/es/date-picker/locale/en_US').default
                }
              />
            </Form.Item>
            
            <Form.Item
              name="manualCodes"
              label={t('jetbrains.productCode')}
            >
              <Input.TextArea 
                placeholder={t('jetbrains.pleaseEnterProductCode')} 
                rows={3}
              />
            </Form.Item>

            <Form.Item>
              <SubmitButton
                type="primary"
                htmlType="submit"
                loading={loading}
              >
                {t('jetbrains.generateActivationCode')}
              </SubmitButton>
            </Form.Item>
          </Form>
        ) : (
          <div>
            <Paragraph>
              {t('jetbrains.serverActivationDescription')}
            </Paragraph>

            {serverRule ? (
              <ResultCard
                title={t('jetbrains.serverConfig')}
                data={{
                  [t('jetbrains.serverAddress')]: serverAddress,
                  [t('jetbrains.powerConfLabel')]: serverRule,
                }}
              />
            ) : (
              <Alert
                message={loadingServerRule ? t('jetbrains.loadingServerRule') : t('jetbrains.serverRuleAutoload')}
                type="info"
                showIcon
                icon={loadingServerRule ? <LoadingOutlined /> : <InfoCircleOutlined />}
              />
            )}
          </div>
        )}
      </FormCard>

      {/* Only show license results if they exist AND we're in the activation method that created them */}
      {license && rawResponse && activationMethod === 'code' && (
        <LicenseResultCard title={<Title level={5} style={{ margin: 0 }}>{t('jetbrains.activationSuccess')}</Title>}>
          <Space direction="vertical" style={{ width: '100%' }}>
            <div>
              <LabelText>{t('jetbrains.product')}:</LabelText>
              <LicenseContent>
                {license.product || t('jetbrains.unknownProduct')}
              </LicenseContent>
            </div>
            
            {extractPowerConf() && (
              <div style={{ marginTop: 16 }}>
                <LabelText>{t('jetbrains.powerConfLabel')}:</LabelText>
                <LicenseContent>
                  {extractPowerConf()}
                  <ButtonContainer>
                    <CopyButton
                      size="small"
                      type="primary"
                      ghost
                      icon={copying['powerConf'] ? <CheckOutlined /> : <CopyOutlined />}
                      onClick={() => copyToClipboard('powerConf', extractPowerConf() || '')}
                    />
                  </ButtonContainer>
                </LicenseContent>
              </div>
            )}
            
            {extractActivationCode() && (
              <div style={{ marginTop: 16 }}>
                <LabelText>{t('jetbrains.activationCode')}:</LabelText>
                <LicenseContent>
                  {extractActivationCode()}
                  <ButtonContainer>
                    <CopyButton
                      size="small"
                      type="primary"
                      ghost
                      icon={copying['activationCode'] ? <CheckOutlined /> : <CopyOutlined />}
                      onClick={() => copyToClipboard('activationCode', extractActivationCode() || '')}
                    />
                  </ButtonContainer>
                </LicenseContent>
              </div>
            )}
          </Space>
        </LicenseResultCard>
      )}
    </div>
  );
};

export default JetBrains; 