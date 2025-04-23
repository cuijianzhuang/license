const jaJP = {
  app: {
    title: 'ライセンス認可サービス',
    language: '言語'
  },
  nav: {
    home: 'ホーム',
    jetbrains: 'JetBrains',
    gitlab: 'GitLab',
    finalshell: 'FinalShell',
    mobaxterm: 'MobaXterm',
    jrebel: 'JRebel'
  },
  common: {
    submit: '送信',
    reset: 'リセット',
    copy: 'コピー',
    copied: 'コピーしました',
    copyFail: 'コピーに失敗しました。手動でコピーしてください',
    generate: '生成',
    success: '成功',
    error: 'エラー',
    loading: '読み込み中',
    useNow: '今すぐ使用'
  },
  home: {
    welcome: 'ライセンス認可サービスへようこそ',
    description: '開発ツールの認可に関するニーズに対応する総合的なソリューションで、様々な一般的な開発ツールのサービスをサポートします',
    tools: {
      jetbrains: {
        title: 'JetBrains認可サービス',
        description: 'IntelliJ IDEA、PyCharm、WebStormなど、すべてのJetBrains製品の認可情報を取得します'
      },
      jrebel: {
        title: 'JRebel認可サービス',
        description: 'JRebelホットデプロイメントツールに認可サービスを提供'
      },
      gitlab: {
        title: 'GitLab認可サービス',
        description: 'GitLabのエンタープライズ認可を作成し、すべての高度な機能のロックを解除'
      },
      finalshell: {
        title: 'FinalShell認可サービス',
        description: 'FinalShell SSHツールの認可情報を取得'
      },
      mobaxterm: {
        title: 'MobaXterm認可サービス',
        description: 'MobaXtermの高度な機能のロックを解除し、専門的な認可を取得'
      }
    }
  },
  jetbrains: {
    title: 'JetBrains認可アシスタント',
    subTitle: 'すべてのJetBrains製品の認可情報を取得',
    activationMethod: '認可方法',
    codeActivation: '認可コード',
    serverActivation: 'オンラインサーバー認可',
    licenseeName: 'ライセンシー名',
    pleaseEnterLicenseeName: 'ライセンシー名を入力してください',
    effectiveDate: '有効日',
    effectiveDatePlaceholder: '日付と時間を選択',
    productCode: '製品コード',
    pleaseEnterProductCode: '製品コードを入力してください。複数の製品はカンマで区切ります',
    generateActivationCode: '認可コードを生成',
    serverActivationDescription: '認可サーバーを設定することでもJetBrains製品を使用できます。以下のpower.conf設定をJetBrains認可サーバー設定にコピーしてください：',
    serverAddress: 'サーバーアドレス',
    serverAddressTooltip: 'ブラウザから現在アクセスしているサーバーアドレスで、JetBrains製品の認可サーバーを設定するために使用されます',
    serverConfig: 'サーバー認可設定',
    loadingServerRule: 'サーバールールを読み込み中です。お待ちください...',
    serverRuleAutoload: 'サーバー認可オプションをクリックした後、サーバールールが自動的に読み込まれます',
    activationSuccess: '認可コードが正常に生成されました',
    product: '製品',
    unknownProduct: '不明な製品',
    powerConfConfig: 'power.conf設定',
    activationCode: '認可コード',
    copySuccess: 'コピーに成功しました',
    copyFail: 'コピーに失敗しました。手動でコピーしてください',
    serverRuleFetchError: 'サーバールールの取得に失敗しました',
    licenseGenerationError: '認可コードの生成に失敗しました',
    powerConfLabel: 'power.conf設定'
  },
  jrebel: {
    title: 'JRebel認可サービス',
    subTitle: 'JRebelホットデプロイメントツールに認可サービスを提供',
    description: 'JRebelは強力なJavaホットデプロイメントツールで、アプリケーションサーバーを再起動せずにコード変更をリアルタイムで確認できます。',
    usageNotice: '使用手順',
    authorizationDescription: 'このサービスはJRebel認可を提供します。専用サーバーアドレスを使用してJRebelを認可します。',
    serverConfig: 'サーバー認可設定',
    configurationDetails: '設定詳細',
    baseServerAddress: 'ベースサーバーアドレス',
    guid: '一意識別子（GUID）',
    regenerateGuid: 'GUIDを再生成',
    regenerateGuidButton: '再生成',
    configurationRules: '設定ルール',
    usageSteps: '使用手順',
    step1: 'IDE（例：IntelliJ IDEA）を開く',
    step2: 'JRebelプラグイン設定を見つける',
    step3: '"Team URL"認可方法を選択',
    step4: 'URLフィールドに上記のサーバーアドレスを入力',
    step5: 'メールフィールドに有効なメールアドレスを入力',
    step6: '"認可"をクリックして認可を完了'
  },
  gitlab: {
    title: 'GitLab認可サービス',
    subTitle: 'GitLabのエンタープライズライセンスを作成',
    description: '以下のフォームに記入してGitLabエンタープライズライセンスを生成します。生成されたライセンスはGitLabエンタープライズエディションのすべての機能を有効化するために使用できます。',
    usageNotice: '使用手順',
    warningDescription: '生成されたGitLabライセンスは学習およびテスト目的のみです。商用利用はできません。',
    form: {
      name: '名前',
      namePlaceholder: '名前を入力してください',
      email: 'メールアドレス',
      emailPlaceholder: 'メールアドレスを入力してください',
      emailInvalid: '有効なメールアドレスを入力してください',
      company: '会社/組織',
      companyPlaceholder: '会社または組織名を入力してください',
      expireTime: '有効期限',
      expireTimePlaceholder: '有効期限を選択してください',
      generateButton: 'ライセンスを生成'
    },
    success: {
      title: 'GitLabライセンスが正常に生成されました',
      name: '名前',
      email: 'メール',
      company: '会社/組織',
      expireTime: '有効期限',
      license: 'ライセンス',
      notSpecified: '指定なし',
      downloadStarted: 'ライセンスファイルのダウンロードが開始されました',
      downloadWarning: 'ライセンス生成が完了していない可能性があります。ダウンロードを確認してください',
      downloadFailed: 'ライセンス生成に失敗しました。再試行してください'
    },
    instructionsTitle: '使用手順',
    usageSteps: {
      step1: 'ライセンスを生成し、ZIPファイルをダウンロード',
      step2: 'ZIPファイルを解凍',
      step3: '/opt/gitlab/embedded/service/gitlab-rails/.license_encryption_key.pubを抽出した.license_encryption_key.pubに置き換え',
      step4: 'GitLabを再起動',
      step5: '管理者アカウントでログインし、左下の「管理エリア」に移動',
      step6: '設定 -> 一般 -> ライセンスの追加に進む',
      step7: 'license.gitlab-licenseファイルをアップロード',
      step8: '「ライセンスを追加」をクリック'
    }
  },
  finalshell: {
    title: 'FinalShell認可サービス',
    subTitle: 'FinalShell SSHツールの認可コードを生成',
    description: 'FinalShellは優れたSSHクライアントツールです。以下のフォームに記入してFinalShellの認可コードを生成し、すべてのプロ機能のロックを解除します。',
    usageNotice: '使用手順',
    warningDescription: '生成された認可コードは学習およびテスト目的のみです。正規ソフトウェアをサポートしてください。',
    machineCode: 'マシンコード',
    enterMachineCode: 'マシンコードを入力してください',
    machineCodeRequired: 'マシンコードが必要です',
    generateButton: '認可コードを生成',
    registrationSuccess: 'FinalShell認可コードが正常に生成されました',
    instructionsTitle: '使用手順',
    usageSteps: {
      step1: 'FinalShellソフトウェアを開く',
      step2: '右上の「有効化/アップグレード」をクリック',
      step3: '「オフライン有効化」を選択',
      step4: 'マシンコードをコピー',
      step5: '認可コードを生成',
      step6: '認可コードを貼り付け',
      step7: '「OK」をクリックして認可を完了'
    },
    versions: {
      advancedBelow396: 'バージョン < 3.9.6 アドバンスドエディション',
      proBelow396: 'バージョン < 3.9.6 プロフェッショナルエディション',
      advancedAbove396: '3.9.6 <= バージョン < 4.5 アドバンスドエディション',
      proAbove396: '3.9.6 <= バージョン < 4.5 プロフェッショナルエディション',
      advancedAbove45: 'バージョン >= 4.5 アドバンスドエディション',
      proAbove45: 'バージョン >= 4.5 プロフェッショナルエディション'
    }
  },
  mobaxterm: {
    title: 'MobaXterm認可サービス',
    subTitle: 'MobaXtermプロフェッショナルエディションの認可コードを生成',
    description: 'MobaXtermはXサーバーとネットワークツールを備えた強力なターミナルツールです。以下のフォームに記入してMobaXtermプロフェッショナルエディションの認可コードを生成します。',
    usageNotice: '使用手順',
    warningDescription: '生成された認可コードは学習およびテスト目的のみです。正規ソフトウェアをサポートしてください。',
    form: {
      username: 'ユーザー名',
      usernamePlaceholder: 'ユーザー名を入力してください',
      version: 'ソフトウェアバージョン',
      versionPlaceholder: 'バージョンを選択してください',
      count: 'ライセンス数',
      countPlaceholder: 'ライセンス数を入力してください',
      countInvalid: '正の整数を入力してください',
      generateButton: '認可コードを生成',
      getAuthCode: '認可コードを取得'
    },
    success: {
      title: 'MobaXterm認可コードが正常に生成されました',
      downloadStarted: 'ライセンスファイル"Custom.mxtpro"のダウンロードが開始されました'
    },
    instructionsTitle: '使用手順',
    usageSteps: {
      step1: '認可ファイルを生成',
      step2: 'Custom.mxtproファイルをMobaXtermインストールディレクトリに配置',
      step3: 'MobaXtermを起動'
    }
  }
};

export default jaJP; 