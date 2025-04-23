import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';

// Import language resources
import enUS from '../locales/lang/en_US';
import zhCN from '../locales/lang/zh_CN';
import jaJP from '../locales/lang/ja_JP';
import koKR from '../locales/lang/ko_KR';

// Get saved language from localStorage or default to browser language
const getSavedLanguage = () => {
  const savedLanguage = localStorage.getItem('i18nextLng');
  if (savedLanguage) {
    if (savedLanguage.startsWith('zh')) return 'zh';
    if (savedLanguage.startsWith('ja')) return 'ja';
    if (savedLanguage.startsWith('ko')) return 'ko';
    return 'en';
  }
  return undefined; // Let language detector decide
};

// Configure i18next
i18n
  // Detect user language
  .use(LanguageDetector)
  // Pass the i18n instance to react-i18next
  .use(initReactI18next)
  // Initialize i18next
  .init({
    resources: {
      en: {
        translation: enUS
      },
      zh: {
        translation: zhCN
      },
      ja: {
        translation: jaJP
      },
      ko: {
        translation: koKR
      }
    },
    lng: getSavedLanguage(), // Try to use saved language first
    fallbackLng: 'en', // Default language is English if detection fails
    interpolation: {
      escapeValue: false // React already safes from XSS
    },
    detection: {
      order: ['localStorage', 'navigator'], // 1. Check localStorage, 2. Check browser language
      lookupLocalStorage: 'i18nextLng',
      caches: ['localStorage'] // Cache user language preference
    }
  });

export default i18n; 