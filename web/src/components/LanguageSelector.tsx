import React, { useEffect, useState } from 'react';
import { Select } from 'antd';
import { useTranslation } from 'react-i18next';

// Helper function to map language code to HTML lang attribute
const mapToHtmlLang = (language: string): string => {
  if (language.startsWith('zh-CN') || language === 'zh-Hans' || language === 'zh_CN') return 'zh-CN';
  if (language.startsWith('zh-TW') || language === 'zh-Hant' || language === 'zh_TW') return 'zh-TW';
  if (language.startsWith('ja')) return 'ja';
  if (language.startsWith('ko')) return 'ko';
  if (language.startsWith('ru')) return 'ru';
  return 'en';
};

const LanguageSelector: React.FC = () => {
  const { i18n } = useTranslation();
  const [currentLanguage, setCurrentLanguage] = useState<string>('en'); // Default to English
  const [isMobile, setIsMobile] = useState<boolean>(window.innerWidth <= 768);
  
  // 添加窗口大小变化监听
  useEffect(() => {
    const handleResize = () => {
      setIsMobile(window.innerWidth <= 768);
    };
    
    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, []);
  
  // Initialize language from localStorage on component mount
  useEffect(() => {
    // Priority 1: Check localStorage for previously configured language
    const savedLanguage = localStorage.getItem('i18nextLng');
    
    if (savedLanguage) {
      let lang = 'en';
      if (savedLanguage.startsWith('zh-CN') || savedLanguage === 'zh-Hans' || savedLanguage === 'zh_CN') lang = 'zh-CN';
      if (savedLanguage.startsWith('zh-TW') || savedLanguage === 'zh-Hant' || savedLanguage === 'zh_TW') lang = 'zh-TW';
      if (savedLanguage.startsWith('ja')) lang = 'ja';
      if (savedLanguage.startsWith('ko')) lang = 'ko';
      if (savedLanguage.startsWith('ru')) lang = 'ru';
      
      setCurrentLanguage(lang);
      
      // Ensure i18n language matches localStorage
      if (i18n.language !== lang) {
        i18n.changeLanguage(lang);
      }

      // Ensure HTML lang attribute is correctly set
      document.documentElement.lang = mapToHtmlLang(lang);
    } else {
      // Priority 2: Check browser language
      const browserLang = navigator.language || (navigator as any).userLanguage;
      let detectedLang = 'en';
      
      if ((browserLang && browserLang.startsWith('zh-CN')) || browserLang === 'zh-Hans') detectedLang = 'zh-CN';
      if ((browserLang && browserLang.startsWith('zh-TW')) || browserLang === 'zh-Hant') detectedLang = 'zh-TW';
      if (browserLang && browserLang.startsWith('ja')) detectedLang = 'ja';
      if (browserLang && browserLang.startsWith('ko')) detectedLang = 'ko';
      if (browserLang && browserLang.startsWith('ru')) detectedLang = 'ru';
      
      setCurrentLanguage(detectedLang);
      i18n.changeLanguage(detectedLang);
      
      // Save to localStorage for future visits
      localStorage.setItem('i18nextLng', detectedLang);

      // Update HTML lang attribute
      document.documentElement.lang = mapToHtmlLang(detectedLang);
    }
  }, [i18n]);

  const handleChange = (value: string) => {
    setCurrentLanguage(value);
    i18n.changeLanguage(value);
    
    // Save the selected language to localStorage
    localStorage.setItem('i18nextLng', value);

    // Update HTML lang attribute
    document.documentElement.lang = mapToHtmlLang(value);
  };

  // 获取语言显示名称
  const getLanguageDisplay = (lang: string): string => {
    switch (lang) {
      case 'zh-CN': return isMobile ? '简中' : '简体中文';
      case 'zh-TW': return isMobile ? '繁中' : '繁體中文';
      case 'ru': return isMobile ? 'RU' : 'Русский';
      case 'en': return isMobile ? 'EN' : 'English';
      case 'ja': return isMobile ? 'JP' : '日本語';
      case 'ko': return isMobile ? 'KR' : '한국어';
      default: return 'English';
    }
  };

  return (
    <Select
      value={currentLanguage}
      style={{ width: 100, maxWidth: '100%' }}
      onChange={handleChange}
      dropdownStyle={{ zIndex: 1100 }}
      popupMatchSelectWidth={false}
      options={[
        { value: 'zh-CN', label: getLanguageDisplay('zh-CN') },
        { value: 'zh-TW', label: getLanguageDisplay('zh-TW') },
        { value: 'ru', label: getLanguageDisplay('ru') },
        { value: 'en', label: getLanguageDisplay('en') },
        { value: 'ja', label: getLanguageDisplay('ja') },
        { value: 'ko', label: getLanguageDisplay('ko') }
      ]}
    />
  );
};

export default LanguageSelector; 