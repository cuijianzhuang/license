import React, { useState } from 'react';
import { Typography, Form, Button, Input, DatePicker, InputNumber, Alert } from 'antd';
import styled from 'styled-components';
import moment from 'moment';
import PageHeader from '../components/PageHeader';
import ResultCard from '../components/ResultCard';
import { gitlab } from '../api';
import { GitLabLicense } from '../types';

const { Paragraph } = Typography;

const FormWrapper = styled.div`
  max-width: 600px;
  margin-bottom: 32px;
`;

const GitLab: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [license, setLicense] = useState<GitLabLicense | null>(null);
  const [form] = Form.useForm();

  const handleGenerateLicense = async (values: {
    company: string;
    email: string;
    userCount: number;
    expiresAt: moment.Moment;
  }) => {
    setLoading(true);
    try {
      const data = await gitlab.generateLicense(
        values.company,
        values.email,
        values.userCount,
        values.expiresAt.format('YYYY-MM-DD')
      );
      setLicense(data);
    } catch (error) {
      console.error('生成许可证失败:', error);
    } finally {
      setLoading(false);
    }
  };

  const breadcrumbs = [
    {
      path: '/',
      breadcrumbName: '首页',
    },
    {
      path: '',
      breadcrumbName: 'GitLab 许可证生成',
    },
  ];

  return (
    <div>
      <PageHeader
        title="GitLab 许可证生成"
        subTitle="为GitLab创建企业版许可证"
        breadcrumbs={breadcrumbs}
      />

      <Paragraph>
        填写以下表单信息，生成GitLab企业版许可证。生成的许可证可用于激活GitLab企业版的所有功能。
      </Paragraph>

      <Alert
        message="注意事项"
        description="生成的GitLab许可证仅供学习和测试使用，请勿用于商业环境。"
        type="warning"
        showIcon
        style={{ marginBottom: 24 }}
      />

      <FormWrapper>
        <Form form={form} onFinish={handleGenerateLicense} layout="vertical">
          <Form.Item
            name="company"
            label="公司/组织名称"
            rules={[{ required: true, message: '请输入公司/组织名称' }]}
          >
            <Input placeholder="请输入公司或组织名称" />
          </Form.Item>

          <Form.Item
            name="email"
            label="邮箱地址"
            rules={[
              { required: true, message: '请输入邮箱地址' },
              { type: 'email', message: '请输入有效的邮箱地址' },
            ]}
          >
            <Input placeholder="请输入邮箱地址" />
          </Form.Item>

          <Form.Item
            name="userCount"
            label="用户数量"
            rules={[{ required: true, message: '请输入用户数量' }]}
            initialValue={100}
          >
            <InputNumber min={1} max={10000} style={{ width: '100%' }} />
          </Form.Item>

          <Form.Item
            name="expiresAt"
            label="过期日期"
            rules={[{ required: true, message: '请选择过期日期' }]}
            initialValue={moment().add(1, 'year')}
          >
            <DatePicker
              format="YYYY-MM-DD"
              style={{ width: '100%' }}
              disabledDate={(current) => {
                return current && current < moment().endOf('day');
              }}
            />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>
              生成许可证
            </Button>
          </Form.Item>
        </Form>
      </FormWrapper>

      {license && (
        <ResultCard
          title="GitLab许可证生成成功"
          data={{
            '公司/组织': license.company || '未指定',
            '邮箱': license.email || '未指定',
            '用户数量': String(license.userCount) || '未指定',
            '过期日期': license.expiresAt || '未指定',
            '许可证': license.license,
          }}
          fileName="gitlab-license.txt"
        />
      )}
    </div>
  );
};

export default GitLab; 