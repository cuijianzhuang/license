import i18n from 'i18next';
import {initReactI18next} from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';

// Import language resources
import zhCN from '../locales/lang/zh_CN';
import zhTW from '../locales/lang/zh_TW';
import ruRU from "../locales/lang/ru_RU";
import enUS from '../locales/lang/en_US';
import jaJP from '../locales/lang/ja_JP';
import koKR from '../locales/lang/ko_KR';

// Get saved language from localStorage or default to browser language
const getSavedLanguage = () => {
    const savedLanguage = localStorage.getItem('i18nextLng');
    if (savedLanguage) {
        if (savedLanguage.startsWith('zh-CN') || savedLanguage === 'zh-Hans' || savedLanguage === 'zh_CN') return 'zh-CN';
        if (savedLanguage.startsWith('zh-TW') || savedLanguage === 'zh-Hant' || savedLanguage === 'zh_TW') return 'zh-TW';
        if (savedLanguage.startsWith('ja')) return 'ja';
        if (savedLanguage.startsWith('ko')) return 'ko';
        if (savedLanguage.startsWith('ru')) return 'ru';
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
            'zh-CN': {
                translation: zhCN
            },
            'zh-TW': {
                translation: zhTW
            },
            ru: {
                translation: ruRU
            },
            en: {
                translation: enUS
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