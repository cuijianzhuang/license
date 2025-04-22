import React from 'react';
import { Select } from 'antd';
import { useTranslation } from 'react-i18next';
import { GlobalOutlined } from '@ant-design/icons';

const { Option } = Select;

const LanguageSelector: React.FC = () => {
  const { i18n, t } = useTranslation();

  const handleChange = (value: string) => {
    i18n.changeLanguage(value);
    // Save the selected language to localStorage
    localStorage.setItem('i18nextLng', value);
  };

  return (
    <Select
      defaultValue={i18n.language.startsWith('zh') ? 'zh' : 'en'}
      style={{ width: 120 }}
      onChange={handleChange}
      dropdownStyle={{ zIndex: 1100 }}
      prefix={<GlobalOutlined />}
    >
      <Option value="zh">中文</Option>
      <Option value="en">English</Option>
    </Select>
  );
};

export default LanguageSelector; 