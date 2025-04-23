const ruRU = {
  app: {
    title: 'Сервис лицензионной авторизации',
    language: 'Язык'
  },
  nav: {
    home: 'Главная',
    jetbrains: 'JetBrains',
    gitlab: 'GitLab',
    finalshell: 'FinalShell',
    mobaxterm: 'MobaXterm',
    jrebel: 'JRebel'
  },
  common: {
    submit: 'Отправить',
    reset: 'Сбросить',
    copy: 'Копировать',
    copied: 'Скопировано',
    copyFail: 'Ошибка копирования, пожалуйста, скопируйте вручную',
    generate: 'Сгенерировать',
    success: 'Успешно',
    error: 'Ошибка',
    loading: 'Загрузка',
    useNow: 'Использовать сейчас'
  },
  home: {
    welcome: 'Добро пожаловать в сервис лицензионной авторизации',
    description: 'Комплексное решение для авторизации инструментов разработки, поддерживающее различные популярные инструменты',
    tools: {
      jetbrains: {
        title: 'Сервис авторизации JetBrains',
        description: 'Получите авторизацию для всех продуктов JetBrains, включая IntelliJ IDEA, PyCharm, WebStorm и др.'
      },
      jrebel: {
        title: 'Сервис авторизации JRebel',
        description: 'Предоставляет авторизацию для инструмента горячей перезагрузки JRebel'
      },
      gitlab: {
        title: 'Сервис авторизации GitLab',
        description: 'Создайте корпоративную авторизацию для GitLab, разблокировав все расширенные функции'
      },
      finalshell: {
        title: 'Сервис авторизации FinalShell',
        description: 'Получите авторизацию для SSH-инструмента FinalShell'
      },
      mobaxterm: {
        title: 'Сервис авторизации MobaXterm',
        description: 'Разблокируйте расширенные функции MobaXterm, получите профессиональную авторизацию'
      }
    }
  },
  jetbrains: {
    title: 'Помощник авторизации JetBrains',
    subTitle: 'Получите авторизацию для всех продуктов JetBrains',
    activationMethod: 'Метод авторизации',
    codeActivation: 'Код авторизации',
    serverActivation: 'Авторизация через онлайн-сервер',
    licenseeName: 'Имя лицензиата',
    pleaseEnterLicenseeName: 'Пожалуйста, введите имя лицензиата',
    effectiveDate: 'Дата вступления в силу',
    effectiveDatePlaceholder: 'Выберите дату и время',
    productCode: 'Код продукта',
    pleaseEnterProductCode: 'Пожалуйста, введите код продукта, разделяйте несколько продуктов запятыми',
    generateActivationCode: 'Сгенерировать код авторизации',
    serverActivationDescription: 'Вы также можете использовать продукты JetBrains, настроив сервер авторизации. Скопируйте конфигурацию power.conf ниже в настройки сервера авторизации JetBrains:',
    serverAddress: 'Адрес сервера',
    serverAddressTooltip: 'Текущий адрес сервера, используемый браузером, для настройки сервера авторизации продуктов JetBrains',
    serverConfig: 'Конфигурация серверной авторизации',
    loadingServerRule: 'Загрузка правил сервера, пожалуйста, подождите...',
    serverRuleAutoload: 'Правила сервера будут автоматически загружены после выбора опции серверной авторизации',
    activationSuccess: 'Код авторизации успешно сгенерирован',
    product: 'Продукт',
    unknownProduct: 'Неизвестный продукт',
    powerConfConfig: 'Конфигурация power.conf',
    activationCode: 'Код авторизации',
    copySuccess: 'Скопировано успешно',
    copyFail: 'Ошибка копирования, пожалуйста, скопируйте вручную',
    serverRuleFetchError: 'Не удалось получить правила сервера',
    licenseGenerationError: 'Не удалось сгенерировать код авторизации',
    powerConfLabel: 'Конфигурация power.conf'
  },
  jrebel: {
    title: 'Сервис авторизации JRebel',
    subTitle: 'Предоставляет авторизацию для инструмента горячей перезагрузки JRebel',
    description: 'JRebel — мощный инструмент горячей перезагрузки Java, позволяющий видеть изменения кода в реальном времени без перезапуска сервера приложений.',
    usageNotice: 'Инструкции по использованию',
    authorizationDescription: 'Этот сервис предоставляет авторизацию JRebel. Используйте выделенный адрес сервера для авторизации JRebel.',
    serverConfig: 'Конфигурация серверной авторизации',
    configurationDetails: 'Детали конфигурации',
    baseServerAddress: 'Базовый адрес сервера',
    guid: 'Уникальный идентификатор (GUID)',
    regenerateGuid: 'Перегенерировать GUID',
    regenerateGuidButton: 'Перегенерировать',
    configurationRules: 'Правила конфигурации',
    usageSteps: 'Инструкции по использованию',
    step1: 'Откройте вашу IDE (например, IntelliJ IDEA)',
    step2: 'Найдите настройки плагина JRebel',
    step3: 'Выберите метод авторизации "Team URL"',
    step4: 'Введите указанный выше адрес сервера в поле URL',
    step5: 'Введите любой действующий адрес электронной почты в поле email',
    step6: 'Нажмите "Авторизовать" для завершения авторизации'
  },
  gitlab: {
    title: 'Сервис авторизации GitLab',
    subTitle: 'Создать корпоративную лицензию для GitLab',
    description: 'Заполните форму ниже, чтобы сгенерировать корпоративную лицензию GitLab. Сгенерированная лицензия может использоваться для активации всех функций GitLab Enterprise Edition.',
    usageNotice: 'Инструкции по использованию',
    warningDescription: 'Сгенерированная лицензия GitLab предназначена только для обучения и тестирования, не для коммерческого использования.',
    form: {
      name: 'Имя',
      namePlaceholder: 'Пожалуйста, введите ваше имя',
      email: 'Адрес электронной почты',
      emailPlaceholder: 'Пожалуйста, введите адрес электронной почты',
      emailInvalid: 'Пожалуйста, введите действительный адрес электронной почты',
      company: 'Компания/Организация',
      companyPlaceholder: 'Пожалуйста, введите название компании или организации',
      expireTime: 'Дата истечения срока',
      expireTimePlaceholder: 'Пожалуйста, выберите дату истечения срока',
      generateButton: 'Сгенерировать лицензию'
    },
    success: {
      title: 'Лицензия GitLab успешно сгенерирована',
      name: 'Имя',
      email: 'Электронная почта',
      company: 'Компания/Организация',
      expireTime: 'Дата истечения срока',
      license: 'Лицензия',
      notSpecified: 'Не указано',
      downloadStarted: 'Началась загрузка файла лицензии',
      downloadWarning: 'Генерация лицензии может быть не завершена, пожалуйста, проверьте загрузку',
      downloadFailed: 'Сбой генерации лицензии, пожалуйста, попробуйте снова'
    },
    instructionsTitle: 'Инструкции по использованию',
    usageSteps: {
      step1: 'Сгенерируйте лицензию и загрузите ZIP-файл',
      step2: 'Распакуйте ZIP-файл',
      step3: 'Замените /opt/gitlab/embedded/service/gitlab-rails/.license_encryption_key.pub извлеченным файлом .license_encryption_key.pub',
      step4: 'Перезапустите GitLab',
      step5: 'Войдите в свою учетную запись администратора и перейдите в "Область администратора" внизу слева',
      step6: 'Перейдите в Настройки -> Общие -> Добавить лицензию',
      step7: 'Загрузите файл license.gitlab-license',
      step8: 'Нажмите "Добавить лицензию"'
    }
  },
  finalshell: {
    title: 'Сервис авторизации FinalShell',
    subTitle: 'Сгенерировать код авторизации для SSH-инструмента FinalShell',
    description: 'FinalShell — отличный SSH-клиент. Заполните форму ниже, чтобы сгенерировать код авторизации для FinalShell и разблокировать все профессиональные функции.',
    usageNotice: 'Инструкции по использованию',
    warningDescription: 'Сгенерированный код авторизации предназначен только для обучения и тестирования. Пожалуйста, поддержите оригинальное программное обеспечение.',
    machineCode: 'Код устройства',
    enterMachineCode: 'Пожалуйста, введите код устройства',
    machineCodeRequired: 'Требуется код устройства',
    generateButton: 'Сгенерировать код авторизации',
    registrationSuccess: 'Код авторизации FinalShell успешно сгенерирован',
    instructionsTitle: 'Инструкции по использованию',
    usageSteps: {
      step1: 'Откройте программу FinalShell',
      step2: 'Нажмите "Активировать/Обновить" в правом верхнем углу',
      step3: 'Выберите "Офлайн-активация"',
      step4: 'Скопируйте код устройства',
      step5: 'Сгенерируйте код авторизации',
      step6: 'Вставьте код авторизации',
      step7: 'Нажмите "ОК" для завершения активации'
    },
    versions: {
      advancedBelow396: 'Версия < 3.9.6 Расширенная редакция',
      proBelow396: 'Версия < 3.9.6 Профессиональная редакция',
      advancedAbove396: '3.9.6 <= Версия < 4.5 Расширенная редакция',
      proAbove396: '3.9.6 <= Версия < 4.5 Профессиональная редакция',
      advancedAbove45: 'Версия >= 4.5 Расширенная редакция',
      proAbove45: 'Версия >= 4.5 Профессиональная редакция'
    }
  },
  mobaxterm: {
    title: 'Сервис авторизации MobaXterm',
    subTitle: 'Сгенерировать код авторизации для MobaXterm Professional Edition',
    description: 'MobaXterm — мощный терминальный инструмент с X-сервером и сетевыми инструментами. Заполните форму ниже, чтобы сгенерировать код авторизации для MobaXterm Professional Edition.',
    usageNotice: 'Инструкции по использованию',
    warningDescription: 'Сгенерированный код авторизации предназначен только для обучения и тестирования. Пожалуйста, поддержите оригинальное программное обеспечение.',
    form: {
      username: 'Имя пользователя',
      usernamePlaceholder: 'Пожалуйста, введите имя пользователя',
      version: 'Версия программы',
      versionPlaceholder: 'Пожалуйста, выберите версию',
      count: 'Количество лицензий',
      countPlaceholder: 'Пожалуйста, введите количество лицензий',
      countInvalid: 'Пожалуйста, введите положительное целое число',
      generateButton: 'Сгенерировать код авторизации',
      getAuthCode: 'Получить код авторизации'
    },
    success: {
      title: 'Код авторизации MobaXterm успешно сгенерирован',
      downloadStarted: 'Началась загрузка лицензионного файла "Custom.mxtpro"'
    },
    instructionsTitle: 'Инструкции по использованию',
    usageSteps: {
      step1: 'Сгенерируйте лицензионный файл',
      step2: 'Поместите файл Custom.mxtpro в директорию установки MobaXterm',
      step3: 'Запустите MobaXterm'
    }
  }
};

export default ruRU; 