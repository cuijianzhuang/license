const koKR = {
  app: {
    title: '라이선스 인증 서비스',
    language: '언어'
  },
  nav: {
    home: '홈',
    jetbrains: 'JetBrains',
    gitlab: 'GitLab',
    finalshell: 'FinalShell',
    mobaxterm: 'MobaXterm',
    jrebel: 'JRebel'
  },
  common: {
    submit: '제출',
    reset: '초기화',
    copy: '복사',
    copied: '복사됨',
    generate: '생성',
    success: '성공',
    error: '오류',
    loading: '로딩 중',
    useNow: '지금 사용'
  },
  home: {
    welcome: '라이선스 인증 서비스에 오신 것을 환영합니다',
    description: '개발 도구 인증 요구 사항에 대한 원스톱 솔루션으로, 다양한 일반 개발 도구에 대한 서비스를 지원합니다',
    tools: {
      jetbrains: {
        title: 'JetBrains 인증 서비스',
        description: 'IntelliJ IDEA, PyCharm, WebStorm 등 모든 JetBrains 제품의 인증 정보를 얻으세요'
      },
      jrebel: {
        title: 'JRebel 인증 서비스',
        description: 'JRebel 핫 배포 도구에 대한 인증 서비스 제공'
      },
      gitlab: {
        title: 'GitLab 인증 서비스',
        description: 'GitLab의 엔터프라이즈 인증을 생성하고 모든 고급 기능 잠금 해제'
      },
      finalshell: {
        title: 'FinalShell 인증 서비스',
        description: 'FinalShell SSH 도구의 인증 정보 얻기'
      },
      mobaxterm: {
        title: 'MobaXterm 인증 서비스',
        description: 'MobaXterm 고급 기능 잠금 해제 및 전문 인증 획득'
      }
    }
  },
  jetbrains: {
    title: 'JetBrains 인증 도우미',
    subTitle: '모든 JetBrains 제품에 대한 인증 정보 획득',
    activationMethod: '인증 방법',
    codeActivation: '인증 코드',
    serverActivation: '온라인 서버 인증',
    licenseeName: '라이선스 사용자 이름',
    pleaseEnterLicenseeName: '라이선스 사용자 이름을 입력하세요',
    effectiveDate: '유효 날짜',
    effectiveDatePlaceholder: '날짜와 시간 선택',
    productCode: '제품 코드',
    pleaseEnterProductCode: '제품 코드를 입력하세요. 여러 제품은 쉼표로 구분하세요',
    generateActivationCode: '인증 코드 생성',
    serverActivationDescription: '인증 서버를 구성하여 JetBrains 제품을 사용할 수도 있습니다. 아래의 power.conf 구성을 JetBrains 인증 서버 설정에 복사하세요:',
    serverAddress: '서버 주소',
    serverAddressTooltip: '브라우저에서 현재 접속 중인 서버 주소로, JetBrains 제품의 인증 서버를 구성하는 데 사용됩니다',
    serverConfig: '서버 인증 구성',
    loadingServerRule: '서버 규칙 로딩 중, 잠시 기다려 주세요...',
    serverRuleAutoload: '서버 인증 옵션을 클릭한 후 서버 규칙이 자동으로 로드됩니다',
    activationSuccess: '인증 코드가 성공적으로 생성되었습니다',
    product: '제품',
    unknownProduct: '알 수 없는 제품',
    powerConfConfig: 'power.conf 구성',
    activationCode: '인증 코드',
    copySuccess: '성공적으로 복사되었습니다',
    copyFail: '복사 실패, 수동으로 복사해 주세요',
    serverRuleFetchError: '서버 규칙 가져오기 실패',
    licenseGenerationError: '인증 코드 생성 실패',
    powerConfLabel: 'power.conf 구성'
  },
  jrebel: {
    title: 'JRebel 인증 서비스',
    subTitle: 'JRebel 핫 배포 도구에 대한 인증 서비스 제공',
    description: 'JRebel은 애플리케이션 서버를 다시 시작하지 않고도 코드 변경 사항을 실시간으로 볼 수 있는 강력한 Java 핫 배포 도구입니다.',
    usageNotice: '사용 안내',
    authorizationDescription: '이 서비스는 JRebel 인증을 제공합니다. 전용 서버 주소를 사용하여 JRebel을 인증하세요.',
    serverConfig: '서버 인증 구성',
    configurationDetails: '구성 세부 정보',
    baseServerAddress: '기본 서버 주소',
    guid: '고유 식별자(GUID)',
    regenerateGuid: 'GUID 재생성',
    regenerateGuidButton: '재생성',
    configurationRules: '구성 규칙',
    usageSteps: '사용 안내',
    step1: 'IDE 열기(예: IntelliJ IDEA)',
    step2: 'JRebel 플러그인 설정 찾기',
    step3: '"Team URL" 인증 방법 선택',
    step4: 'URL 필드에 위의 서버 주소 입력',
    step5: '이메일 필드에 유효한 이메일 주소 입력',
    step6: '"인증"을 클릭하여 인증 완료'
  },
  gitlab: {
    title: 'GitLab 인증 서비스',
    subTitle: 'GitLab용 엔터프라이즈 라이선스 생성',
    description: '아래 양식을 작성하여 GitLab 엔터프라이즈 라이선스를 생성합니다. 생성된 라이선스는 GitLab 엔터프라이즈 에디션의 모든 기능을 활성화하는 데 사용할 수 있습니다.',
    usageNotice: '사용 안내',
    warningDescription: '생성된 GitLab 라이선스는 학습 및 테스트 목적으로만 사용됩니다. 상업적 용도로는 사용하지 마세요.',
    form: {
      name: '이름',
      namePlaceholder: '이름을 입력하세요',
      email: '이메일 주소',
      emailPlaceholder: '이메일 주소를 입력하세요',
      emailInvalid: '유효한 이메일 주소를 입력하세요',
      company: '회사/조직',
      companyPlaceholder: '회사 또는 조직 이름을 입력하세요',
      expireTime: '만료 날짜',
      expireTimePlaceholder: '만료 날짜를 선택하세요',
      generateButton: '라이선스 생성'
    },
    success: {
      title: 'GitLab 라이선스가 성공적으로 생성되었습니다',
      name: '이름',
      email: '이메일',
      company: '회사/조직',
      expireTime: '만료 날짜',
      license: '라이선스',
      notSpecified: '지정되지 않음',
      downloadStarted: '라이선스 파일 다운로드가 시작되었습니다',
      downloadWarning: '라이선스 생성이 완료되지 않았을 수 있습니다. 다운로드를 확인하세요',
      downloadFailed: '라이선스 생성에 실패했습니다. 다시 시도하세요'
    },
    instructionsTitle: '사용 안내',
    usageSteps: {
      step1: '라이선스를 생성하고 ZIP 파일 다운로드',
      step2: 'ZIP 파일 압축 해제',
      step3: '/opt/gitlab/embedded/service/gitlab-rails/.license_encryption_key.pub를 추출된 .license_encryption_key.pub로 교체',
      step4: 'GitLab 재시작',
      step5: '관리자 계정으로 로그인하고 왼쪽 하단의 "관리자 영역"으로 이동',
      step6: '설정 -> 일반 -> 라이선스 추가로 이동',
      step7: 'license.gitlab-license 파일 업로드',
      step8: '"라이선스 추가" 클릭'
    }
  },
  finalshell: {
    title: 'FinalShell 인증 서비스',
    subTitle: 'FinalShell SSH 도구의 인증 코드 생성',
    description: 'FinalShell은 우수한 SSH 클라이언트 도구입니다. 아래 양식을 작성하여 FinalShell의 인증 코드를 생성하고 모든 전문 기능 잠금을 해제하세요.',
    usageNotice: '사용 안내',
    warningDescription: '생성된 인증 코드는 학습 및 테스트 목적으로만 사용됩니다. 원본 소프트웨어를 지원해 주세요.',
    machineCode: '기계 코드',
    enterMachineCode: '기계 코드를 입력하세요',
    machineCodeRequired: '기계 코드가 필요합니다',
    generateButton: '인증 코드 생성',
    registrationSuccess: 'FinalShell 인증 코드가 성공적으로 생성되었습니다',
    instructionsTitle: '사용 안내',
    usageSteps: {
      step1: 'FinalShell 소프트웨어 열기',
      step2: '오른쪽 상단의 "활성화/업그레이드" 클릭',
      step3: '"오프라인 활성화" 선택',
      step4: '기계 코드 복사',
      step5: '인증 코드 생성',
      step6: '인증 코드 붙여넣기',
      step7: '"확인"을 클릭하여 인증 완료'
    },
    versions: {
      advancedBelow396: '버전 < 3.9.6 고급 에디션',
      proBelow396: '버전 < 3.9.6 전문가 에디션',
      advancedAbove396: '버전 >= 3.9.6 고급 에디션',
      proAbove396: '버전 >= 3.9.6 전문가 에디션'
    }
  },
  mobaxterm: {
    title: 'MobaXterm 인증 서비스',
    subTitle: 'MobaXterm 프로페셔널 에디션의 인증 코드 생성',
    description: 'MobaXterm은 X 서버 및 네트워크 도구가 통합된 강력한 터미널 도구입니다. 아래 양식을 작성하여 MobaXterm 프로페셔널 에디션의 인증 코드를 생성하세요.',
    usageNotice: '사용 안내',
    warningDescription: '생성된 인증 코드는 학습 및 테스트 목적으로만 사용됩니다. 원본 소프트웨어를 지원해 주세요.',
    form: {
      username: '사용자 이름',
      usernamePlaceholder: '사용자 이름을 입력하세요',
      version: '소프트웨어 버전',
      versionPlaceholder: '버전을 선택하세요',
      count: '라이선스 수',
      countPlaceholder: '라이선스 수를 입력하세요',
      countInvalid: '양의 정수를 입력하세요',
      generateButton: '인증 코드 생성',
      getAuthCode: '인증 코드 받기'
    },
    success: {
      title: 'MobaXterm 인증 코드가 성공적으로 생성되었습니다',
      downloadStarted: '라이선스 파일 "Custom.mxtpro" 다운로드가 시작되었습니다'
    },
    instructionsTitle: '사용 안내',
    usageSteps: {
      step1: '인증 파일 생성',
      step2: 'Custom.mxtpro 파일을 MobaXterm 설치 디렉토리에 배치',
      step3: 'MobaXterm 시작'
    }
  }
};

export default koKR; 