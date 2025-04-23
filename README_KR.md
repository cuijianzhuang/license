# 라이선스 관리 서비스

![Docker Pulls](https://img.shields.io/docker/pulls/ghostx/license-server)

[简体中文](README_CN.md) | [繁體中文](README_TW.md) | [Русский](README_RU.md) | [English](README.md) | [日本語](README_JP.md) | [한국어](README_KR.md)

Go 언어 기반 라이선스 관리 서비스입니다.

## 특징

- 다양한 소프트웨어 제품을 위한 라이선스 생성 및 검증
- JetBrains 제품, GitLab, FinalShell, MobaXterm 및 JRebel 지원
- Gin 프레임워크로 구축된 RESTful API 인터페이스
- cron을 사용한 예약 작업
- GORM(MySQL/SQLite 지원)을 통한 데이터베이스 저장
- RSA를 사용한 안전한 암호화
- 다국어 지원(간체 중국어, 번체 중국어, 러시아어, 영어, 일본어, 한국어)

## 시스템 요구사항

- Go 1.24 이상
- MySQL 데이터베이스(개발 환경에서는 SQLite 사용 가능)
- Docker(선택 사항, 컨테이너화 배포용)

## 설치 방법

### 방법 1: 직접 설치

1. 저장소 복제
   ```
   git clone https://github.com/nannanStrawberry314/license.git
   cd license
   ```

2. 종속성 설치
   ```
   go mod download
   ```

3. 환경 변수 구성(.env.example을 .env로 복사하고 필요에 따라 수정)

4. 빌드 및 실행
   ```
   go build -o license-server
   ./license-server
   ```

### 방법 2: Docker 배포

1. Docker 이미지 빌드
   ```
   docker build -t license-server .
   ```

2. docker-compose를 사용하여 실행
   ```
   docker-compose up -d
   ```

## 구성

구성은 환경 변수 및 `.env` 파일을 통해 관리됩니다:

- `HTTP_HOST`: 서버를 바인딩할 호스트 주소
- `HTTP_PORT`: 수신 대기할 포트
- `DB_TYPE`: 데이터베이스 유형(mysql 또는 sqlite)
- `DB_DSN`: 데이터베이스 연결 문자열

## GitHub Actions 구성

이 프로젝트에는 자동 빌드, 테스트 및 릴리스를 위한 GitHub Actions 워크플로우가 포함되어 있습니다. 이러한 워크플로우를 사용하려면 다음 저장소 암호를 구성해야 합니다:

- `HUB_USER`: Docker Hub 사용자 이름
- `HUB_PASS`: Docker Hub 비밀번호
- `HUB_REPO`: Docker Hub 저장소 이름
- `PUBLIC_REPO_TOKEN`: (선택 사항) 공개 저장소에 대한 쓰기 권한이 있는 개인 액세스 토큰
- `PUBLIC_REPO`: (선택 사항) `사용자이름/저장소` 형식의 공개 저장소 경로

워크플로우에는 다음이 포함됩니다:
- 각 푸시에 대한 빌드 및 테스트
- 태그 푸시 또는 수동 트리거 시 Docker 이미지 빌드 및 게시
- GitHub 릴리스 생성
- 공개 저장소에 릴리스 동기화(수동 트리거 시에만 실행, `PUBLIC_REPO_TOKEN` 및 `PUBLIC_REPO` 필요)

## API 엔드포인트

API는 라이선스 관리를 위한 다양한 엔드포인트를 제공합니다:

- `POST /v1/generate`: 새 라이선스 생성
- `POST /v1/validate`: 기존 라이선스 검증
- `GET /v1/status`: 서비스 상태 확인

자세한 사용 지침은 API 문서를 참조하십시오.

## 개발

이 프로젝트에 기여하려면:

1. 저장소 포크
2. 기능 브랜치 생성
3. 변경 사항 적용
4. 풀 리퀘스트 제출

## 면책 조항

### 프로젝트 목적

이 프로젝트는 기술 연구 및 학술 학습 기회를 제공하기 위해 교육 목적으로만 개발 및 공유되었습니다. 이 프로젝트에서 논의되는 기술과 방법은 학습 및 연구 목적으로만 사용됩니다.

### 사용 제한

이 프로젝트의 저자는 소프트웨어 불법 복제 또는 불법 활동을 장려하기 위한 것이 아니라 지식 공유와 기술 발전을 촉진하기 위해 이 프로젝트를 게시했습니다. **소프트웨어를 크랙하거나, 활성화하거나, 불법적으로 사용하기 위해 이 프로젝트를 사용하는 것은 엄격히 금지됩니다**. 사용자는 이 프로젝트를 사용할 때 모든 관련 지역 및 국제 법률과 규정을 준수해야 합니다.

### 상업적 사용 금지

이 프로젝트는 **상업적 용도나 이익을 창출하는 활동에 사용하는 것을 엄격히 금지**합니다. 직접적 또는 간접적으로 경제적 이익을 창출할 수 있는 상황에서 프로젝트의 어떤 부분도 사용해서는 안 됩니다.

### 책임 면제

이 프로젝트는 명시적이든 묵시적이든 어떠한 보증 없이 "있는 그대로" 제공됩니다. 프로젝트 저자는 이 프로젝트 콘텐츠 사용으로 인해 발생할 수 있는 어떤 형태의 손해나 법적 결과에 대해서도 책임을 지지 않습니다. 이 프로젝트를 사용함으로써 귀하는 이러한 조건을 이해하고 동의하며, 이 프로젝트 사용과 관련된 모든 위험을 감수해야 합니다.

## 스타 히스토리

[![Star History Chart](https://api.star-history.com/svg?repos=nannanStrawberry314/license&type=Date)](https://www.star-history.com/#nannanStrawberry314/license&Date) 