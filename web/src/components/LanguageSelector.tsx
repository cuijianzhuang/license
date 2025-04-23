import React, { useEffect, useState } from 'react';
import { Select } from 'antd';
import { useTranslation } from 'react-i18next';
import { GlobalOutlined } from '@ant-design/icons';

const { Option } = Select;

const LanguageSelector: React.FC = () => {
  const { i18n, t } = useTranslation();
  const [currentLanguage, setCurrentLanguage] = useState<string>('en'); // Default to English
  
  // Initialize language from localStorage on component mount
  useEffect(() => {
    // Priority 1: Check localStorage for previously configured language
    const savedLanguage = localStorage.getItem('i18nextLng');
    
    if (savedLanguage) {
      let lang = 'en';
      if (savedLanguage.startsWith('zh')) lang = 'zh';
      if (savedLanguage.startsWith('ja')) lang = 'ja';
      if (savedLanguage.startsWith('ko')) lang = 'ko';
      
      setCurrentLanguage(lang);
      
      // Ensure i18n language matches localStorage
      if (i18n.language !== lang) {
        i18n.changeLanguage(lang);
      }
    } else {
      // Priority 2: Check browser language
      const browserLang = navigator.language || (navigator as any).userLanguage;
      let detectedLang = 'en';
      
      if (browserLang && browserLang.startsWith('zh')) detectedLang = 'zh';
      if (browserLang && browserLang.startsWith('ja')) detectedLang = 'ja';
      if (browserLang && browserLang.startsWith('ko')) detectedLang = 'ko';
      
      setCurrentLanguage(detectedLang);
      i18n.changeLanguage(detectedLang);
      
      // Save to localStorage for future visits
      localStorage.setItem('i18nextLng', detectedLang);
    }
  }, [i18n]);

  const handleChange = (value: string) => {
    setCurrentLanguage(value);
    i18n.changeLanguage(value);
    
    // Save the selected language to localStorage
    localStorage.setItem('i18nextLng', value);
  };

  return (
    <Select
      value={currentLanguage}
      style={{ width: 120 }}
      onChange={handleChange}
      dropdownStyle={{ zIndex: 1100 }}
      prefix={<GlobalOutlined />}
    >
      <Option value="zh">中文</Option>
      <Option value="en">English</Option>
      <Option value="ja">日本語</Option>
      <Option value="ko">한국어</Option>
    </Select>
  );
};

export default LanguageSelector; 