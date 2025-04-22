import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';

// Import language resources
import enUS from '../locales/lang/en_US';
import zhCN from '../locales/lang/zh_CN';

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
      }
    },
    fallbackLng: 'zh', // Default language if detection fails
    interpolation: {
      escapeValue: false // React already safes from XSS
    },
    detection: {
      order: ['navigator', 'localStorage', 'cookie'],
      caches: ['localStorage', 'cookie'] // Cache user language preference
    }
  });

export default i18n; 