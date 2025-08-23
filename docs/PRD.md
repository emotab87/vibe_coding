# RealWorld Conduit 앱 제품 요구사항 문서 (PRD)

## 1. 개요

### 1.1 제품 소개

- **제품명**: Conduit
- **제품 유형**: 소셜 블로깅 플랫폼 (Medium.com 클론)
- **목표**: 실제 운영 환경에서 사용 가능한 수준의 풀스택 웹 애플리케이션

### 1.2 프로젝트 배경

RealWorld 프로젝트는 단순한 "Todo" 데모를 넘어서는 실제 애플리케이션 구축 방법을 보여주는 오픈소스 프로젝트입니다. 다양한 프론트엔드와 백엔드 기술 스택을 조합하여 동일한 API 스펙을 공유하는 표준화된 애플리케이션을 구현합니다.

## 2. 핵심 기능 요구사항

### 2.1 사용자 인증 및 관리

- 사용자 등록 (회원가입)
- 사용자 로그인/로그아웃
- JWT 토큰 기반 인증 시스템
- 사용자 프로필 관리 및 수정

### 2.2 게시글 관리

- 게시글 작성, 수정, 삭제 (CRUD)
- 마크다운 지원
- 게시글 태그 기능
- 게시글 즐겨찾기 (좋아요) 기능
- 게시글 목록 조회 및 페이지네이션

### 2.3 소셜 기능

- 다른 사용자 팔로우/언팔로우
- 사용자 프로필 조회
- 게시글에 댓글 작성, 수정, 삭제
- 팔로우한 사용자들의 게시글 피드

### 2.4 검색 및 필터링

- 태그별 게시글 필터링
- 작성자별 게시글 필터링
- 즐겨찾기한 게시글 필터링
- 전체 태그 목록 조회

## 3. API 스펙

### 3.1 인증 관련 엔드포인트

```
POST /api/users/login      # 로그인
POST /api/users            # 회원가입
GET  /api/user             # 현재 사용자 정보
PUT  /api/user             # 사용자 정보 수정
```

### 3.2 프로필 관련 엔드포인트

```
GET    /api/profiles/:username         # 프로필 조회
POST   /api/profiles/:username/follow  # 팔로우
DELETE /api/profiles/:username/follow  # 언팔로우
```

### 3.3 게시글 관련 엔드포인트

```
GET    /api/articles           # 게시글 목록
GET    /api/articles/feed      # 팔로우 피드
GET    /api/articles/:slug     # 특정 게시글
POST   /api/articles           # 게시글 작성
PUT    /api/articles/:slug     # 게시글 수정
DELETE /api/articles/:slug     # 게시글 삭제
```

### 3.4 즐겨찾기 관련 엔드포인트

```
POST   /api/articles/:slug/favorite  # 즐겨찾기 추가
DELETE /api/articles/:slug/favorite  # 즐겨찾기 해제
```

### 3.5 댓글 관련 엔드포인트

```
GET    /api/articles/:slug/comments     # 댓글 목록
POST   /api/articles/:slug/comments     # 댓글 작성
DELETE /api/articles/:slug/comments/:id # 댓글 삭제
```

### 3.6 태그 관련 엔드포인트

```
GET /api/tags  # 전체 태그 목록
```

### 3.7 인증 헤더

```
Authorization: Token jwt.token.here
```

## 4. 프론트엔드 요구사항

### 4.1 필수 페이지

1. **홈 페이지**: 게시글 목록, 태그 클라우드
2. **로그인 페이지**: 사용자 로그인 폼
3. **회원가입 페이지**: 사용자 등록 폼
4. **설정 페이지**: 사용자 정보 수정
5. **프로필 페이지**: 사용자 프로필 및 게시글 목록
6. **게시글 상세 페이지**: 게시글 내용 및 댓글
7. **게시글 편집기**: 마크다운 에디터

### 4.2 라우팅 요구사항

- SPA (Single Page Application) 구조
- 브라우저 히스토리 API 지원
- 인증이 필요한 페이지에 대한 라우트 가드

### 4.3 상태 관리

- 전역 상태 관리 (사용자 인증 상태, 게시글 데이터 등)
- 로컬 상태 관리 (폼 입력, UI 상태 등)

### 4.4 UI/UX 요구사항

- 반응형 웹 디자인
- shadcn/ui + Tailwind CSS 기반 스타일링
  - 재사용 가능한 컴포넌트 라이브러리
  - 유틸리티 우선 CSS 접근법
  - 다크 모드 지원
  - 접근성 기본 제공
- 로딩 상태 표시 (스켈레톤 UI, 스피너)
- 에러 처리 및 사용자 피드백 (토스트, 알림)
- 페이지네이션 (무한 스크롤 또는 페이지 번호)

## 5. 기술적 요구사항

### 5.1 API 통신

- RESTful API 설계 원칙 준수
- JSON 데이터 형식
- HTTP 상태 코드 적절한 활용
- CORS 지원
- 명시적이고 서술적인 엔드포인트 네이밍

### 5.2 보안 요구사항

- JWT 토큰 기반 인증
- XSS 공격 방지
- 입력 데이터 검증 및 새니타이징
- 민감한 정보 로그 기록 방지
- 권한 검사의 로컬화 및 가시화

### 5.3 성능 요구사항

- 페이지네이션을 통한 대용량 데이터 처리
- API 응답 시간 최적화
- 클라이언트 사이드 캐싱 (TanStack Query)
- 이미지 및 정적 자산 최적화

### 5.4 코드 품질 요구사항

- **단순성과 가독성 우선**: 복잡한 추상화보다 명확한 코드 작성
- **명시적 함수명**: 길고 서술적인 함수명 사용으로 의도 명확화
- **의존성 최소화**: 필요한 라이브러리만 선택적 사용
- **직접 SQL 사용**: ORM보다 명시적 SQL 쿼리 작성
- **에러 처리 명확화**: AI가 이해하기 쉬운 명확한 에러 메시지

### 5.5 관찰 가능성 요구사항

- **포괄적인 로깅**: 디버깅과 모니터링을 위한 구조화된 로그
- **상태 추적**: 시스템 상태를 쉽게 파악할 수 있는 구조
- **성능 모니터링**: 응답 시간 및 리소스 사용량 추적

## 6. 데이터 모델

### 6.1 사용자 (User)

```json
{
  "id": "number",
  "username": "string",
  "email": "string",
  "bio": "string",
  "image": "string",
  "token": "string"
}
```

### 6.2 게시글 (Article)

```json
{
  "slug": "string",
  "title": "string",
  "description": "string",
  "body": "string",
  "tagList": ["string"],
  "createdAt": "datetime",
  "updatedAt": "datetime",
  "favorited": "boolean",
  "favoritesCount": "number",
  "author": "Profile"
}
```

### 6.3 프로필 (Profile)

```json
{
  "username": "string",
  "bio": "string",
  "image": "string",
  "following": "boolean"
}
```

### 6.4 댓글 (Comment)

```json
{
  "id": "number",
  "createdAt": "datetime",
  "updatedAt": "datetime",
  "body": "string",
  "author": "Profile"
}
```

## 7. 개발 환경 및 배포

### 7.1 컨테이너화 환경

- **Docker 기반 개발 환경**
  - Dockerfile 및 docker-compose.yml 구성
  - 프론트엔드, 백엔드, 데이터베이스 서비스 분리
  - 개발/스테이징/프로덕션 환경 일관성 보장
  - 볼륨 마운트를 통한 개발 환경 최적화

### 7.2 개발 환경

- **로컬 개발 지원**
  - docker-compose up으로 전체 환경 구동
  - 핫 리로드 기능 (Vite, Go Air 등)
  - 개발자 도구 지원
  - 환경 변수 관리 (.env 파일)

- **데이터베이스**
  - SQLite 파일 기반 데이터베이스
  - 볼륨 마운트를 통한 데이터 영속성
  - 스키마 마이그레이션 스크립트

### 7.3 테스트

- **API 테스트**
  - Postman 컬렉션을 통한 API 엔드포인트 테스트
  - Go 표준 테스트 프레임워크 (백엔드)
  - 테스트 데이터베이스 분리

- **프론트엔드 테스트**
  - Vitest를 사용한 유닛 테스트
  - React Testing Library를 사용한 컴포넌트 테스트
  - Playwright를 사용한 E2E 테스트

### 7.4 배포

- **컨테이너 배포**
  - 멀티 스테이지 Docker 빌드
  - 프로덕션 최적화된 이미지 생성
  - docker-compose를 통한 오케스트레이션

- **CI/CD**
  - GitHub Actions 또는 GitLab CI 지원
  - 자동 테스트 실행
  - 컨테이너 이미지 빌드 및 배포
  - 환경별 배포 파이프라인

## 8. 외부 리소스

### 8.1 공개 API

- 데모 API: `https://api.realworld.build/api`
- 실제 API: `https://api.realworld.io/api`

### 8.2 개발 리소스

- GitHub 저장소: https://github.com/gothinkster/realworld
- 공식 문서: https://realworld-docs.netlify.app/
- 다양한 구현체: https://codebase.show/projects/realworld

## 9. 성공 지표

- 모든 API 엔드포인트 정상 동작
- Postman 테스트 컬렉션 통과
- 반응형 웹 디자인 구현
- 크로스 브라우저 호환성
- 접근성 가이드라인 준수

## 10. 기술 스택

### 10.1 백엔드 기술
- **언어**: Go
  - 명시적 컨텍스트 시스템
  - 단순성과 가독성
  - AI 에이전트 친화적 구조
  - 낮은 생태계 변동성
  - 구조적 인터페이스 지원

### 10.2 프론트엔드 기술
- **프레임워크**: React 18+
- **빌드 도구**: Vite
- **라우팅**: TanStack Router
- **상태 관리**: TanStack Query (서버 상태) + React useState/useReducer (로컬 상태)
- **스타일링**: Tailwind CSS + shadcn/ui
  - 유틸리티 우선 CSS 프레임워크
  - 재사용 가능한 컴포넌트 라이브러리
  - 타입스크립트 지원
  - 접근성 기본 제공

### 10.3 데이터베이스
- **데이터베이스**: SQLite
  - 개발 환경 단순성
  - 파일 기반 데이터베이스
  - 배포 용이성
  - 트랜잭션 지원
- **쿼리 방식**: 직접 SQL 사용 (ORM 최소화)

### 10.4 실행 환경
- **컨테이너화**: Docker
  - 일관된 개발/배포 환경
  - 의존성 격리
  - 크로스 플랫폼 지원
  - docker-compose를 통한 개발 환경 구성

### 10.5 개발 도구
- **패키지 매니저**: npm/yarn
- **타입 시스템**: TypeScript (프론트엔드)
- **코드 포맷팅**: Prettier + ESLint
- **테스트**: Go 표준 테스트 (백엔드), Vitest (프론트엔드)

## 11. 개발 철학 및 원칙

### 11.1 코드 작성 원칙
- **단순성 우선**: 복잡한 추상화보다 명확하고 단순한 코드
- **명시적 함수명**: 길고 서술적인 함수명 사용
- **의존성 최소화**: 필요한 라이브러리만 추가
- **코드 생성 우선**: 새 라이브러리 추가보다 코드 생성 선호

### 11.2 AI 친화적 개발
- **명확한 에러 메시지**: AI가 이해하기 쉬운 에러 처리
- **포괄적인 로깅**: 디버깅과 모니터링을 위한 상세 로그
- **관찰 가능성**: 시스템 상태를 쉽게 파악할 수 있는 구조
- **보안 고려**: 오용 방지를 위한 보호 장치

### 11.3 데이터베이스 접근 방식
- **직접 SQL 사용**: 복잡한 ORM 대신 명시적 SQL 쿼리
- **권한 검사 로컬화**: 중요한 검사 로직을 가시적이고 로컬하게 유지
- **명시적 트랜잭션**: 데이터 일관성을 위한 명시적 트랜잭션 관리

## 12. 향후 확장 가능성

- 실시간 알림 시스템
- 이미지 업로드 기능
- 소셜 미디어 연동
- 검색 기능 고도화
- 모바일 앱 버전