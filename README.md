# RealWorld Conduit App - Vibe Coding Edition

> ### Go + React codebase containing real world examples (CRUD, auth, advanced patterns, etc) that adheres to the [RealWorld](https://github.com/gothinkster/realworld) spec and API.

이 프로젝트는 **바이브 코딩(Vibe Coding)** 기법과 **아르민 로나허(Armin Ronacher)의 추천 기술 스택**을 사용하여 RealWorld 앱을 구현하는 것을 목표로 합니다.

### [Demo](https://demo.realworld.build/) &nbsp;&nbsp;&nbsp;&nbsp; [RealWorld](https://github.com/gothinkster/realworld)

This codebase was created to demonstrate a fully fledged fullstack application built with **Go + React** including CRUD operations, authentication, routing, pagination, and more.

We've gone to great lengths to adhere to the **Go** and **React** community styleguides & best practices.

For more information on how to this works with other frontends/backends, head over to the [RealWorld](https://github.com/gothinkster/realworld) repo.

## 🎯 프로젝트 목적

### 1. 바이브 코딩 실험
- **직관적 개발**: 사전 엄밀한 설계보다 직감과 느낌에 의존하는 개발 방식 실험
- **AI 보조 개발**: ChatGPT, Claude 등 생성형 AI를 적극 활용한 코드 생성 및 문제 해결
- **빠른 프로토타이핑**: 완전한 이해보다는 동작하는 결과물 우선의 개발 접근법

### 2. 아르민 로나허의 개발 철학 적용
- **의존성 최소화**: "The more I build software, the more I despise dependencies"
- **직접 구현 우선**: 라이브러리 추가보다 코드 생성 및 직접 구현 선호
- **실용적 접근**: 이론적 완벽함보다 실제 동작하는 솔루션 중시

### 3. 학습 및 실험 목적
- RealWorld 표준 스펙을 통한 실제 애플리케이션 개발 경험
- Go 백엔드와 React 프론트엔드의 실무 패턴 학습
- 현대적 개발 방법론과 전통적 개발 방법론의 균형점 찾기

## 🛠 기술 스택

### 백엔드 (Go)
- **언어**: Go 1.21+
- **웹 프레임워크**: Gorilla Mux (의존성 최소화)
- **데이터베이스**: SQLite (직접 SQL 사용, ORM 최소화)
- **인증**: JWT 토큰 (직접 구현)
- **패스워드**: bcrypt

### 프론트엔드 (React)
- **프레임워크**: React 18+ with TypeScript
- **빌드 도구**: Vite
- **라우팅**: React Router
- **스타일링**: Tailwind CSS
- **API 통신**: Axios
- **상태 관리**: React Context API + useState/useReducer

### 개발 환경
- **컨테이너화**: Docker & Docker Compose
- **개발 도구**: Hot reload, ESLint, Prettier

## 🚀 Getting Started

### 사전 요구사항
- Docker & Docker Compose
- Go 1.21+ (로컬 개발 시)
- Node.js 18+ (로컬 개발 시)

### 빠른 시작

```bash
# 저장소 클론
git clone https://github.com/emotab87/vibe_coding.git
cd vibe_coding

# Docker로 전체 환경 실행
docker-compose up

# 또는 개별 실행
# 백엔드
cd backend && go run cmd/main.go

# 프론트엔드 (새 터미널)
cd frontend && npm install && npm run dev
```

앱이 다음 주소에서 실행됩니다:
- **프론트엔드**: http://localhost:3000
- **백엔드 API**: http://localhost:8080/api

### API 문서
백엔드 서버 실행 후 http://localhost:8080/api-docs 에서 API 문서를 확인할 수 있습니다.

## 📁 프로젝트 구조

```
.
├── backend/                 # Go 백엔드
│   ├── cmd/                # 애플리케이션 진입점
│   ├── internal/           # 내부 패키지
│   │   ├── handlers/       # HTTP 핸들러
│   │   ├── models/         # 데이터 모델
│   │   ├── middleware/     # 미들웨어
│   │   └── database/       # 데이터베이스 로직
│   └── migrations/         # DB 마이그레이션
├── frontend/               # React 프론트엔드
│   ├── src/
│   │   ├── components/     # 재사용 컴포넌트
│   │   ├── pages/          # 페이지 컴포넌트
│   │   ├── hooks/          # 커스텀 훅
│   │   ├── services/       # API 서비스
│   │   └── types/          # TypeScript 타입
│   └── public/
├── docs/                   # 문서
├── docker-compose.yml      # Docker 설정
└── README.md
```

## 🎨 구현된 기능 (MVP)

### ✅ 핵심 기능
- [x] 사용자 인증 (회원가입, 로그인, 로그아웃)
- [x] 게시글 CRUD (작성, 조회, 수정, 삭제)
- [x] 댓글 시스템 (작성, 조회, 삭제)
- [x] 반응형 UI

### 🔄 개발 예정 기능
- [ ] 프로필 관리
- [ ] 팔로우/언팔로우
- [ ] 게시글 즐겨찾기
- [ ] 태그 시스템
- [ ] 게시글 검색 및 필터링

## 🧪 바이브 코딩 실험 결과

### 적용된 바이브 코딩 원칙
1. **AI 보조 개발**: 대부분의 코드를 Claude/ChatGPT와 함께 작성
2. **직감적 구현**: 완벽한 설계서 없이 기능부터 구현
3. **빠른 반복**: 동작 확인 후 점진적 개선
4. **의존성 최소화**: 필요한 경우 라이브러리보다 직접 구현 선택

### 학습된 교훈
- AI 보조 개발의 효과적 활용법
- 전통적 개발과 바이브 코딩의 균형점
- 실무에서의 적용 가능성과 한계

## 📚 학습 리소스

### RealWorld 관련
- [RealWorld 스펙](https://realworld-docs.netlify.app/)
- [다양한 구현체](https://codebase.show/projects/realworld)
- [데모 API](https://api.realworld.build/api)

### 바이브 코딩 & 아르민 로나허
- [Armin Ronacher's Blog](https://lucumr.pocoo.org/)
- [Welcoming The Next Generation of Programmers](https://lucumr.pocoo.org/2025/7/20/the-next-generation/)
- [Agentic Coding Recommendations](https://lucumr.pocoo.org/2025/6/12/agentic-coding/)

## 🤝 기여하기

이 프로젝트는 학습 및 실험 목적이지만, 기여를 환영합니다!

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 라이선스

이 프로젝트는 MIT 라이선스 하에 배포됩니다. 자세한 내용은 [LICENSE](LICENSE) 파일을 참조하세요.

## 🙋‍♂️ 연락처

- **프로젝트 링크**: [https://github.com/emotab87/vibe_coding](https://github.com/emotab87/vibe_coding)
- **RealWorld**: [https://github.com/gothinkster/realworld](https://github.com/gothinkster/realworld)

---

> 💡 **참고**: 이 프로젝트는 바이브 코딩과 아르민 로나허의 개발 철학을 실험하는 교육적 목적으로 만들어졌습니다. 실제 프로덕션 환경에서는 더 엄격한 설계와 테스트가 필요할 수 있습니다.