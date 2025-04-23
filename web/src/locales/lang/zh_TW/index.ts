const zhTW = {
  app: {
    title: '許可證授權服務',
    language: '語言'
  },
  nav: {
    home: '首頁',
    jetbrains: 'JetBrains',
    gitlab: 'GitLab',
    finalshell: 'FinalShell',
    mobaxterm: 'MobaXterm',
    jrebel: 'JRebel'
  },
  common: {
    submit: '提交',
    reset: '重置',
    copy: '複製',
    copied: '已複製',
    copyFail: '複製失敗，請手動複製',
    generate: '生成',
    success: '成功',
    error: '錯誤',
    loading: '加載中',
    useNow: '立即使用'
  },
  home: {
    welcome: '歡迎使用許可證授權服務',
    description: '一站式開發工具授權需求解決方案，支持各種常用開發工具的服務',
    tools: {
      jetbrains: {
        title: 'JetBrains授權服務',
        description: '獲取所有JetBrains產品的授權信息，包括IntelliJ IDEA、PyCharm、WebStorm等'
      },
      jrebel: {
        title: 'JRebel授權服務',
        description: '為JRebel熱部署工具提供授權服務'
      },
      gitlab: {
        title: 'GitLab授權服務',
        description: '創建GitLab企業授權，解鎖所有高級功能'
      },
      finalshell: {
        title: 'FinalShell授權服務',
        description: '獲取FinalShell SSH工具的授權信息'
      },
      mobaxterm: {
        title: 'MobaXterm授權服務',
        description: '解鎖MobaXterm高級功能，獲取專業授權'
      }
    }
  },
  jetbrains: {
    title: 'JetBrains授權服務',
    subTitle: '獲取所有JetBrains產品的授權信息',
    activationMethod: '授權方式',
    codeActivation: '授權碼',
    serverActivation: '在線服務器授權',
    licenseeName: '授權名稱',
    pleaseEnterLicenseeName: '請輸入授權名稱',
    effectiveDate: '生效日期',
    effectiveDatePlaceholder: '選擇日期和時間',
    productCode: '產品代碼',
    pleaseEnterProductCode: '請輸入產品代碼，多個產品用逗號分隔',
    generateActivationCode: '生成授權碼',
    serverActivationDescription: '您也可以通過配置授權服務器來使用JetBrains產品。複製下面的power.conf配置到JetBrains授權服務器設置中：',
    serverAddress: '服務器地址',
    serverAddressTooltip: '瀏覽器當前訪問的服務器地址，用於配置JetBrains產品的授權服務器',
    serverConfig: '服務器授權配置',
    loadingServerRule: '正在加載服務器規則，請稍候...',
    serverRuleAutoload: '點擊服務器授權選項後，服務器規則將自動加載',
    activationSuccess: '授權碼生成成功',
    product: '產品',
    unknownProduct: '未知產品',
    powerConfConfig: 'power.conf配置',
    activationCode: '授權碼',
    copySuccess: '複製成功',
    copyFail: '複製失敗，請手動複製',
    serverRuleFetchError: '獲取服務器規則失敗',
    licenseGenerationError: '生成授權碼失敗',
    powerConfLabel: 'power.conf配置'
  },
  jrebel: {
    title: 'JRebel授權服務',
    subTitle: '為JRebel熱部署工具提供授權服務',
    description: 'JRebel是一款強大的Java熱部署工具，能讓您實時查看代碼修改，無需重啟應用服務器。',
    usageNotice: '使用說明',
    authorizationDescription: '本服務提供JRebel授權。使用專用服務器地址授權JRebel。',
    serverConfig: '服務器授權配置',
    configurationDetails: '配置詳情',
    baseServerAddress: '基礎服務器地址',
    guid: '唯一標識符',
    regenerateGuid: '重新生成GUID',
    regenerateGuidButton: '重新生成',
    configurationRules: '配置規則',
    usageSteps: '使用說明',
    step1: '打開你的IDE（例如IntelliJ IDEA）',
    step2: '找到JRebel插件設置',
    step3: '選擇"Team URL"授權方式',
    step4: '在URL字段中輸入上方服務器地址',
    step5: '在郵箱字段中輸入任意有效的郵箱地址',
    step6: '點擊"授權"完成授權'
  },
  gitlab: {
    title: 'GitLab授權服務',
    subTitle: '為GitLab創建企業版許可證',
    description: '填寫以下表單信息，生成GitLab企業版許可證。生成的許可證可用於激活GitLab企業版的所有功能。',
    usageNotice: '使用說明',
    warningDescription: '生成的GitLab許可證僅供學習和測試使用，請勿用於商業環境。',
    form: {
      name: '授權',
      namePlaceholder: '請輸入授權名稱',
      email: '電子郵件',
      emailPlaceholder: '請輸入電子郵件',
      emailInvalid: '請輸入有效的電子郵件',
      company: '公司',
      companyPlaceholder: '請輸入公司名稱',
      expireTime: '過期日期',
      expireTimePlaceholder: '請選擇過期日期',
      generateButton: '生成許可證'
    },
    success: {
      title: 'GitLab許可證生成成功',
      name: '授權',
      email: '電子郵件',
      company: '公司',
      expireTime: '過期日期',
      license: '許可證',
      notSpecified: '未指定',
      downloadStarted: '許可證文件下載已開始',
      downloadWarning: '許可證生成可能未完成，請檢查下載',
      downloadFailed: '許可證生成失敗，請重試'
    },
    instructionsTitle: '使用說明',
    usageSteps: {
      step1: '生成許可證，下載壓縮包',
      step2: '解壓壓縮包',
      step3: '將 .license_encryption_key.pub 替換到 /opt/gitlab/embedded/service/gitlab-rails/.license_encryption_key.pub',
      step4: '重啟GitLab',
      step5: '登錄管理員賬號，導航到左下角【管理員】',
      step6: '選擇【設置】->【通用】->【添加許可證】',
      step7: '將 license.gitlab-license 上傳',
      step8: '點擊添加許可'
    }
  },
  finalshell: {
    title: 'FinalShell 授權服務',
    subTitle: '生成FinalShell SSH工具的授權碼',
    description: 'FinalShell是一款優秀的SSH客戶端工具，填寫以下表單生成FinalShell的授權碼，解鎖所有專業功能。',
    usageNotice: '使用說明',
    warningDescription: '生成的授權碼僅供學習和測試使用，請支持正版軟件。',
    machineCode: '機器碼',
    enterMachineCode: '請輸入機器碼',
    machineCodeRequired: '請輸入機器碼',
    generateButton: '生成授權碼',
    registrationSuccess: 'FinalShell授權碼生成成功',
    instructionsTitle: '使用說明',
    usageSteps: {
      step1: '打開FinalShell軟件',
      step2: '點擊右上角"激活/升級"',
      step3: '選擇"離線激活"',
      step4: '複製機器碼',
      step5: '生成授權碼',
      step6: '粘貼授權碼',
      step7: '點擊"確定"完成授權'
    },
    versions: {
      advancedBelow396: '版本號 < 3.9.6 高級版',
      proBelow396: '版本號 < 3.9.6 專業版',
      advancedAbove396: '3.9.6 <= 版本號 < 4.5 高級版',
      proAbove396: '3.9.6 <= 版本號 < 4.5 專業版',
      advancedAbove45: '版本號 >= 4.5 高級版',
      proAbove45: '版本號 >= 4.5 專業版'
    }
  },
  mobaxterm: {
    title: 'MobaXterm 授權服務',
    subTitle: '生成 MobaXterm 專業版授權碼',
    description: 'MobaXterm 是一款功能強大的終端工具，集成了 X 伺服器和網路工具。填寫下方表單生成 MobaXterm 專業版授權碼。',
    usageNotice: '使用說明',
    warningDescription: '生成的授權碼僅供學習和測試使用。請支持正版軟件。',
    form: {
      username: '用戶名',
      usernamePlaceholder: '請輸入用戶名',
      version: '軟件版本',
      versionPlaceholder: '請選擇版本',
      count: '許可證數量',
      countPlaceholder: '請輸入許可證數量',
      countInvalid: '請輸入正整數',
      generateButton: '生成授權碼',
      getAuthCode: '獲取授權碼'
    },
    success: {
      title: 'MobaXterm 授權碼生成成功',
      downloadStarted: '授權文件"Custom.mxtpro"下載已開始'
    },
    instructionsTitle: '使用說明',
    usageSteps: {
      step1: '生成授權文件',
      step2: '將Custom.mxtpro 放到 MobaXterm 安裝目錄',
      step3: '啟動 MobaXterm'
    }
  }
};

export default zhTW; 