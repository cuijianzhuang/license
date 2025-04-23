const enUS = {
  app: {
    title: 'License Authorization Service',
    language: 'Language'
  },
  nav: {
    home: 'Home',
    jetbrains: 'JetBrains',
    gitlab: 'GitLab',
    finalshell: 'FinalShell',
    mobaxterm: 'MobaXterm',
    jrebel: 'JRebel'
  },
  common: {
    submit: 'Submit',
    reset: 'Reset',
    copy: 'Copy',
    copied: 'Copied',
    generate: 'Generate',
    success: 'Success',
    error: 'Error',
    loading: 'Loading',
    useNow: 'Use Now'
  },
  home: {
    welcome: 'Welcome to License Authorization Service',
    description: 'One-stop solution for development tool authorization needs, supporting services for various common development tools',
    tools: {
      jetbrains: {
        title: 'JetBrains Authorization Service',
        description: 'Obtain authorization information for all JetBrains products, including IntelliJ IDEA, PyCharm, WebStorm, etc.'
      },
      jrebel: {
        title: 'JRebel Authorization Service',
        description: 'Provide authorization service for JRebel hot deployment tool'
      },
      gitlab: {
        title: 'GitLab Authorization Service',
        description: 'Create enterprise authorization for GitLab, unlocking all advanced features'
      },
      finalshell: {
        title: 'FinalShell Authorization Service',
        description: 'Obtain authorization information for FinalShell SSH tools'
      },
      mobaxterm: {
        title: 'MobaXterm Authorization Service',
        description: 'Unlock MobaXterm advanced features, obtain professional authorization'
      }
    }
  },
  jetbrains: {
    title: 'JetBrains Authorization Assistant',
    subTitle: 'Obtain authorization information for all JetBrains products',
    activationMethod: 'Authorization Method',
    codeActivation: 'Authorization Code',
    serverActivation: 'Online Server Authorization',
    licenseeName: 'Licensee Name',
    pleaseEnterLicenseeName: 'Please enter licensee name',
    effectiveDate: 'Effective Date',
    effectiveDatePlaceholder: 'Select date and time',
    productCode: 'Product Code',
    pleaseEnterProductCode: 'Please enter product code, separate multiple products with commas',
    generateActivationCode: 'Generate Authorization Code',
    serverActivationDescription: 'You can also use JetBrains products by configuring an authorization server. Copy the power.conf configuration below to the JetBrains authorization server settings:',
    serverAddress: 'Server Address',
    serverAddressTooltip: 'Current server address accessed by the browser, used to configure the authorization server for JetBrains products',
    serverConfig: 'Server Authorization Configuration',
    loadingServerRule: 'Loading server rules, please wait...',
    serverRuleAutoload: 'Server rules will be automatically loaded after clicking the server authorization option',
    activationSuccess: 'Authorization Code Generated Successfully',
    product: 'Product',
    unknownProduct: 'Unknown Product',
    powerConfConfig: 'power.conf Configuration',
    activationCode: 'Authorization Code',
    copySuccess: 'Copied successfully',
    copyFail: 'Copy failed, please copy manually',
    serverRuleFetchError: 'Failed to fetch server rules',
    licenseGenerationError: 'Failed to generate authorization code',
    powerConfLabel: 'power.conf Configuration'
  },
  jrebel: {
    title: 'JRebel Authorization Service',
    subTitle: 'Provide authorization service for JRebel hot deployment tool',
    description: 'JRebel is a powerful Java hot deployment tool that allows you to see code changes in real-time without restarting the application server.',
    usageNotice: 'Usage Instructions',
    authorizationDescription: 'This service provides JRebel authorization. Use the dedicated server address to authorize JRebel.',
    serverConfig: 'Server Authorization Configuration',
    configurationDetails: 'Configuration Details',
    baseServerAddress: 'Base Server Address',
    guid: 'Unique Identifier (GUID)',
    regenerateGuid: 'Regenerate GUID',
    regenerateGuidButton: 'Regenerate',
    configurationRules: 'Configuration Rules',
    usageSteps: 'Usage Instructions',
    step1: 'Open your IDE (e.g., IntelliJ IDEA)',
    step2: 'Find the JRebel plugin settings',
    step3: 'Choose "Team URL" authorization method',
    step4: 'Enter the server address above in the URL field',
    step5: 'Enter any valid email address in the email field',
    step6: 'Click "Authorize" to complete the authorization'
  },
  gitlab: {
    title: 'GitLab Authorization Service',
    subTitle: 'Create enterprise license for GitLab',
    description: 'Fill in the form below to generate a GitLab enterprise license. The generated license can be used to activate all features of GitLab Enterprise Edition.',
    usageNotice: 'Usage Instructions',
    warningDescription: 'The generated GitLab license is for learning and testing purposes only, not for commercial use.',
    form: {
      name: 'Name',
      namePlaceholder: 'Please enter your name',
      email: 'Email Address',
      emailPlaceholder: 'Please enter email address',
      emailInvalid: 'Please enter a valid email address',
      company: 'Company/Organization',
      companyPlaceholder: 'Please enter company or organization name',
      expireTime: 'Expiration Date',
      expireTimePlaceholder: 'Please select expiration date',
      generateButton: 'Generate License'
    },
    success: {
      title: 'GitLab License Generated Successfully',
      name: 'Name',
      email: 'Email',
      company: 'Company/Organization',
      expireTime: 'Expiration Date',
      license: 'License',
      notSpecified: 'Not specified',
      downloadStarted: 'License file download has started',
      downloadWarning: 'License generation may not be complete, please check download',
      downloadFailed: 'License generation failed, please try again'
    },
    instructionsTitle: 'Usage Instructions',
    usageSteps: {
      step1: 'Generate the license and download the ZIP file',
      step2: 'Extract the ZIP file',
      step3: 'Replace /opt/gitlab/embedded/service/gitlab-rails/.license_encryption_key.pub with the extracted .license_encryption_key.pub',
      step4: 'Restart GitLab',
      step5: 'Log in to your admin account and navigate to the "Admin Area" in the bottom left',
      step6: 'Go to Settings -> General -> Add license',
      step7: 'Upload the license.gitlab-license file',
      step8: 'Click "Add license"'
    }
  },
  finalshell: {
    title: 'FinalShell Authorization Service',
    subTitle: 'Generate authorization code for FinalShell SSH tool',
    description: 'FinalShell is an excellent SSH client tool. Fill in the form below to generate an authorization code for FinalShell to unlock all professional features.',
    usageNotice: 'Usage Instructions',
    warningDescription: 'The generated authorization code is for learning and testing purposes only. Please support the original software.',
    machineCode: 'Machine Code',
    enterMachineCode: 'Please enter machine code',
    machineCodeRequired: 'Machine code is required',
    generateButton: 'Generate Authorization Code',
    registrationSuccess: 'FinalShell Authorization Code Generated Successfully',
    instructionsTitle: 'Usage Instructions',
    usageSteps: {
      step1: 'Open FinalShell software',
      step2: 'Click "Activate/Upgrade" in the top right corner',
      step3: 'Select "Offline Activation"',
      step4: 'Copy the machine code',
      step5: 'Generate the authorization code',
      step6: 'Paste the authorization code',
      step7: 'Click "OK" to complete the activation'
    },
    versions: {
      advancedBelow396: 'Version < 3.9.6 Advanced Edition',
      proBelow396: 'Version < 3.9.6 Professional Edition',
      advancedAbove396: 'Version >= 3.9.6 Advanced Edition',
      proAbove396: 'Version >= 3.9.6 Professional Edition'
    }
  },
  mobaxterm: {
    title: 'MobaXterm Authorization Service',
    subTitle: 'Generate authorization code for MobaXterm Professional Edition',
    description: 'MobaXterm is a powerful terminal tool with X server and network tools. Fill in the form below to generate an authorization code for MobaXterm Professional Edition.',
    usageNotice: 'Usage Instructions',
    warningDescription: 'The generated authorization code is for learning and testing purposes only. Please support the original software.',
    form: {
      username: 'Username',
      usernamePlaceholder: 'Please enter username',
      version: 'Software Version',
      versionPlaceholder: 'Please select version',
      count: 'License Count',
      countPlaceholder: 'Please enter license count',
      countInvalid: 'Please enter a positive integer',
      generateButton: 'Generate Authorization Code',
      getAuthCode: 'Get Authorization Code'
    },
    success: {
      title: 'MobaXterm Authorization Code Generated Successfully',
      downloadStarted: 'License file "Custom.mxtpro" download has started'
    },
    instructionsTitle: 'Usage Instructions',
    usageSteps: {
      step1: 'Generate the license file',
      step2: 'Place the Custom.mxtpro file in your MobaXterm installation directory',
      step3: 'Start MobaXterm'
    }
  }
};

export default enUS; 