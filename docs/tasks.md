# RealWorld Conduit 앱 MVP 구현 작업 목록

이 문서는 RealWorld Conduit 앱의 MVP(최소 실행 가능한 제품)를 구현하기 위한 핵심 작업 목록입니다.

## 1. 프로젝트 기반 설정

### 1.1 프로젝트 구조 초기화
- [ ] 백엔드 디렉토리 구조 생성 (`backend/`)
- [ ] 프론트엔드 디렉토리 구조 생성 (`frontend/`)
- [ ] `.gitignore` 파일 생성 및 설정
- [ ] `README.md` 작성 (기본 실행 가이드)

### 1.2 개발 환경 설정
- [ ] `docker-compose.yml` 작성 (개발 환경)
- [ ] `.env.example` 파일 생성
- [ ] 데이터베이스 연결 설정

## 2. 백엔드 구현 (Go) - MVP

### 2.1 프로젝트 기본 구조
- [ ] Go 모듈 초기화 (`go mod init`)
- [ ] 기본 디렉토리 구조 생성:
  - `cmd/` - 애플리케이션 진입점
  - `internal/` - 내부 패키지
  - `migrations/` - 데이터베이스 마이그레이션
- [ ] 기본 의존성 추가 (gorilla/mux, database/sql, JWT 등)

### 2.2 데이터베이스 설계 및 구현
- [ ] SQLite 데이터베이스 스키마 설계
- [ ] 핵심 테이블 생성:
  - `users` 테이블
  - `articles` 테이블
  - `comments` 테이블
- [ ] 데이터베이스 연결 설정

### 2.3 인증 시스템 (핵심)
- [ ] JWT 토큰 생성/검증 유틸리티
- [ ] 비밀번호 해싱 (bcrypt)
- [ ] 사용자 등록 API (`POST /api/users`)
- [ ] 사용자 로그인 API (`POST /api/users/login`)
- [ ] 현재 사용자 정보 조회 API (`GET /api/user`)

### 2.4 게시글 관리 (핵심)
- [ ] 게시글 목록 조회 API (`GET /api/articles`)
- [ ] 게시글 상세 조회 API (`GET /api/articles/:slug`)
- [ ] 게시글 작성 API (`POST /api/articles`)
- [ ] 게시글 수정 API (`PUT /api/articles/:slug`)
- [ ] 게시글 삭제 API (`DELETE /api/articles/:slug`)

### 2.5 댓글 시스템 (기본)
- [ ] 댓글 목록 조회 API (`GET /api/articles/:slug/comments`)
- [ ] 댓글 작성 API (`POST /api/articles/:slug/comments`)
- [ ] 댓글 삭제 API (`DELETE /api/articles/:slug/comments/:id`)

### 2.6 기본 에러 처리
- [ ] 기본 에러 처리 미들웨어
- [ ] 표준화된 API 응답 형식

## 3. 프론트엔드 구현 (React + Vite) - MVP

### 3.1 프로젝트 초기 설정
- [ ] Vite + React + TypeScript 프로젝트 생성
- [ ] 필수 의존성 설치:
  - React Router (라우팅)
  - Tailwind CSS (스타일링)
  - axios (API 호출)
- [ ] 기본 설정 및 구조

### 3.2 라우팅 설정 (핵심)
- [ ] React Router 설정
- [ ] 기본 라우트 정의:
  - `/` - 홈페이지
  - `/login` - 로그인
  - `/register` - 회원가입
  - `/article/:slug` - 게시글 상세
  - `/editor` - 게시글 편집기

### 3.3 인증 시스템 (핵심)
- [ ] 인증 상태 관리 (Context API)
- [ ] JWT 토큰 저장/관리 (localStorage)
- [ ] 로그인 폼 컴포넌트
- [ ] 회원가입 폼 컴포넌트
- [ ] 로그아웃 기능

### 3.4 기본 UI 컴포넌트
- [ ] 헤더/네비게이션 컴포넌트
- [ ] 로딩 스피너
- [ ] 기본 레이아웃

### 3.5 홈페이지 및 게시글 목록 (핵심)
- [ ] 홈페이지 레이아웃
- [ ] 게시글 목록 컴포넌트
- [ ] 게시글 카드 컴포넌트

### 3.6 게시글 관리 (핵심)
- [ ] 게시글 상세 페이지
- [ ] 게시글 편집기 (기본 텍스트 에디터)
- [ ] 게시글 작성/수정 기능

### 3.7 댓글 시스템 (기본)
- [ ] 댓글 목록 컴포넌트
- [ ] 댓글 작성 폼
- [ ] 댓글 삭제 기능

### 3.8 API 통신
- [ ] API 서비스 함수 작성:
  - 인증 관련 API
  - 게시글 관련 API
  - 댓글 관련 API

## 4. MVP 완성을 위한 통합

### 4.1 기본 통합 테스트
- [ ] 프론트엔드-백엔드 API 연동 확인
- [ ] 기본 기능 동작 테스트

### 4.2 기본 배포
- [ ] Docker로 로컬 환경 실행 확인
- [ ] 기본 README 작성

## MVP 개발 순서 (학습 최적화)

### Week 1: 백엔드 기초
1. Go 프로젝트 설정 및 기본 구조
2. SQLite 데이터베이스 연결
3. 사용자 인증 API (회원가입, 로그인)
4. 기본 게시글 CRUD API

### Week 2: 프론트엔드 기초  
1. React 프로젝트 설정
2. 라우팅 및 기본 레이아웃
3. 로그인/회원가입 페이지
4. 홈페이지 및 게시글 목록

### Week 3: 기능 완성
1. 게시글 상세 페이지 및 편집기
2. 댓글 시스템 구현
3. 프론트엔드-백엔드 통합
4. 기본 스타일링 및 사용성 개선

## 학습 목표별 중요도

### 🔥 필수 학습 (MVP 핵심)
- Go 기본 문법 및 웹 서버
- React 컴포넌트 및 상태 관리  
- REST API 설계 및 구현
- 데이터베이스 기본 CRUD

### ⚡ 권장 학습 (추가 구현 시)
- JWT 인증 심화
- React Router 고급 기능
- 에러 처리 및 검증
- Docker 컨테이너화

### 💡 선택 학습 (시간 여유 시)
- 팔로우/즐겨찾기 기능
- 태그 시스템
- 무한 스크롤
- 반응형 디자인

## 참고 자료

- [RealWorld 프로젝트 공식 문서](https://docs.realworld.show/)
- [데모 API](https://api.realworld.build/api) - 테스트용
- [Go 웹 개발 가이드](https://golang.org/doc/)
- [React 공식 문서](https://react.dev/)

---

**MVP 목표**: 기본적인 블로그 플랫폼으로서 게시글을 작성하고 읽을 수 있으며, 사용자 인증이 가능한 수준까지 구현하는 것이 1차 목표입니다.