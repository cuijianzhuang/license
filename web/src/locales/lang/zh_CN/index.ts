const zhCN = {
  app: {
    title: '许可证授权服务',
    language: '语言'
  },
  nav: {
    home: '首页',
    jetbrains: 'JetBrains',
    gitlab: 'GitLab',
    finalshell: 'FinalShell',
    mobaxterm: 'MobaXterm',
    jrebel: 'JRebel'
  },
  common: {
    submit: '提交',
    reset: '重置',
    copy: '复制',
    copied: '已复制',
    generate: '生成',
    success: '成功',
    error: '错误',
    loading: '加载中',
    useNow: '立即使用'
  },
  home: {
    welcome: '欢迎使用许可证授权服务',
    description: '一站式开发工具授权需求解决方案，支持各种常用开发工具的服务',
    tools: {
      jetbrains: {
        title: 'JetBrains授权服务',
        description: '获取所有JetBrains产品的授权信息，包括IntelliJ IDEA、PyCharm、WebStorm等'
      },
      gitlab: {
        title: 'GitLab授权服务',
        description: '创建GitLab企业授权，解锁所有高级功能'
      },
      finalshell: {
        title: 'FinalShell授权服务',
        description: '获取FinalShell SSH工具的授权信息'
      },
      mobaxterm: {
        title: 'MobaXterm授权服务',
        description: '解锁MobaXterm高级功能，获取专业授权'
      },
      jrebel: {
        title: 'JRebel授权服务',
        description: '为JRebel热部署工具提供授权服务'
      }
    }
  },
  jetbrains: {
    title: 'JetBrains授权服务',
    subTitle: '获取所有JetBrains产品的授权信息',
    activationMethod: '授权方式',
    codeActivation: '授权码',
    serverActivation: '在线服务器授权',
    licenseeName: '授权名称',
    pleaseEnterLicenseeName: '请输入授权名称',
    effectiveDate: '生效日期',
    effectiveDatePlaceholder: '选择日期和时间',
    productCode: '产品代码',
    pleaseEnterProductCode: '请输入产品代码，多个产品用逗号分隔',
    generateActivationCode: '生成授权码',
    serverActivationDescription: '您也可以通过配置授权服务器来使用JetBrains产品。复制下面的power.conf配置到JetBrains授权服务器设置中：',
    serverAddress: '服务器地址',
    serverAddressTooltip: '浏览器当前访问的服务器地址，用于配置JetBrains产品的授权服务器',
    serverConfig: '服务器授权配置',
    loadingServerRule: '正在加载服务器规则，请稍候...',
    serverRuleAutoload: '点击服务器授权选项后，服务器规则将自动加载',
    activationSuccess: '授权码生成成功',
    product: '产品',
    unknownProduct: '未知产品',
    powerConfConfig: 'power.conf配置',
    activationCode: '授权码',
    copySuccess: '复制成功',
    copyFail: '复制失败，请手动复制',
    serverRuleFetchError: '获取服务器规则失败',
    licenseGenerationError: '生成授权码失败',
    powerConfLabel: 'power.conf配置'
  },
  gitlab: {
    title: 'GitLab授权服务',
    subTitle: '为GitLab创建企业版许可证',
    description: '填写以下表单信息，生成GitLab企业版许可证。生成的许可证可用于激活GitLab企业版的所有功能。',
    warning: '注意事项',
    warningDescription: '生成的GitLab许可证仅供学习和测试使用，请勿用于商业环境。',
    form: {
      name: '授权',
      namePlaceholder: '请输入授权名称',
      email: '电子邮件',
      emailPlaceholder: '请输入电子邮件',
      emailInvalid: '请输入有效的电子邮件',
      company: '公司',
      companyPlaceholder: '请输入公司名称',
      expireTime: '过期日期',
      expireTimePlaceholder: '请选择过期日期',
      generateButton: '生成许可证'
    },
    success: {
      title: 'GitLab许可证生成成功',
      name: '授权',
      email: '电子邮件',
      company: '公司',
      expireTime: '过期日期',
      license: '许可证',
      notSpecified: '未指定',
      downloadStarted: '许可证文件下载已开始',
      downloadWarning: '许可证生成可能未完成，请检查下载',
      downloadFailed: '许可证生成失败，请重试'
    }
  },
  mobaxterm: {
    title: 'MobaXterm 授权服务',
    subTitle: '生成MobaXterm专业版激活码',
    description: 'MobaXterm是一款功能强大的终端工具，集成X服务器和网络工具。填写以下表单生成MobaXterm专业版激活码。',
    warning: '注意事项',
    warningDescription: '生成的激活码仅供学习和测试使用，请支持正版软件。',
    form: {
      username: '用户名',
      usernamePlaceholder: '请输入用户名',
      version: '软件版本',
      versionPlaceholder: '请选择版本',
      count: '授权数量',
      countPlaceholder: '请输入授权数量',
      generateButton: '生成激活码'
    },
    success: {
      title: 'MobaXterm激活码生成成功',
      downloadStarted: '许可证文件下载已开始'
    },
    instructionsTitle: '使用说明',
    usageSteps: {
      step1: '打开MobaXterm软件',
      step2: '点击右上角的"?"按钮，然后选择"Register"',
      step3: '输入上面生成的用户名和激活码',
      step4: '点击"确定"完成激活'
    }
  },
  jrebel: {
    title: 'JRebel授权服务',
    subTitle: '为JRebel热部署工具提供授权服务',
    description: 'JRebel是一款强大的Java热部署工具，能让您实时查看代码修改，无需重启应用服务器。',
    activationTitle: '激活说明',
    activationDescription: '本服务提供JRebel激活。使用专用服务器地址激活JRebel。',
    serverConfig: '服务器授权配置',
    configurationDetails: '配置详情',
    baseServerAddress: '基础服务器地址',
    guid: '唯一标识符',
    regenerateGuid: '重新生成GUID',
    regenerateGuidButton: '重新生成',
    configurationRules: '配置规则',
    activationSteps: 'JRebel激活步骤',
    step1: '打开你的IDE（例如IntelliJ IDEA）',
    step2: '找到JRebel插件设置',
    step3: '选择"Team URL"激活方式',
    step4: '在URL字段中输入上方服务器地址',
    step5: '在邮箱字段中输入任意有效的邮箱地址',
    step6: '点击"激活"完成激活'
  },
  finalshell: {
    title: 'FinalShell 授权服务',
    subTitle: '生成FinalShell SSH工具的注册码',
    description: 'FinalShell是一款优秀的SSH客户端工具，填写以下表单生成FinalShell的注册码，解锁所有专业功能。',
    warning: '注意事项',
    warningDescription: '生成的注册码仅供学习和测试使用，请支持正版软件。',
    machineCode: '机器码',
    enterMachineCode: '请输入机器码',
    machineCodeRequired: '请输入机器码',
    generateButton: '生成注册码',
    registrationSuccess: 'FinalShell注册码生成成功',
    instructionsTitle: '使用说明',
    usageSteps: {
      step1: '打开FinalShell软件',
      step2: '点击"帮助" > "注册"',
      step3: '输入上面生成的用户名和注册码',
      step4: '点击"确定"完成注册'
    },
    versions: {
      advancedBelow396: '版本号 < 3.9.6 高级版',
      proBelow396: '版本号 < 3.9.6 专业版',
      advancedAbove396: '版本号 >= 3.9.6 高级版',
      proAbove396: '版本号 >= 3.9.6 专业版'
    }
  }
};

export default zhCN; 