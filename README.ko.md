> 🌐 **[English](README.md)** | **한국어**

# LOADSTAR CLI

Go 기반 LOADSTAR 메타데이터 관리 도구. AI 에이전트와 사람이 공유하는 프로젝트 작업 단위(WayPoint)·인덱스(Map)·TODO·로그를 명령행에서 관리합니다.

> 📌 LOADSTAR가 처음이라면 먼저 [openLoadstar 전체 안내](https://github.com/openLoadstar/openLoadstar) 와 [방법론 명세(spec)](https://github.com/openLoadstar/spec) 를 참고하세요.

---

## 🛠️ 설치

### 사전 요구사항

- Go 1.21 이상

### 빌드

```bash
git clone https://github.com/openLoadstar/cli.git
cd cli
go build -o bin/loadstar.exe .
```

빌드된 바이너리(`bin/loadstar.exe`)를 PATH에 추가하거나 절대경로로 호출합니다.

### 빠른 동작 확인

```bash
loadstar --help
```

---

## 📋 명령어

| 명령 | 용도 |
|:---|:---|
| `loadstar init` | `.loadstar/` 디렉토리 구조 초기화 |
| `loadstar show [FILTER] [--recent]` | WayPoint 목록 (주소·STATUS·LAST_MODIFIED) — 키워드 필터 + 최근 수정순 정렬 |
| `loadstar validate` | 모든 요소의 참조 무결성 검증, 깨진 링크 보고 |
| `loadstar log [TIME_RANGE] [FILTER]` | 변경 로그 검색 — `7d`, `3h` 같은 기간 + 키워드 필터 |
| `loadstar log add <ADDR> <KIND> "<MSG>"` | 로그 항목 직접 추가 |
| `loadstar todo sync` | WayPoint STATUS 기반 TODO 자동 동기화 |
| `loadstar todo list` | 현재 PENDING / ACTIVE / BLOCKED 작업 목록 |
| `loadstar todo history [MAP_ADDR]` | 완료된 TECH_SPEC 항목 이력 |
| `loadstar question [FILTER] [--with-resolved]` | 미해결 OPEN_QUESTIONS 조회 |
| `loadstar question done <ADDR> <QID>` | RESOLVED 질문을 DONE으로 전환 |
| `loadstar question close <ADDR> <QID> [사유]` | 결정 파일 없이 질문 직접 종료 |
| `loadstar question stats` | OPEN/DEFERRED/RESOLVED/DONE 집계 |

> 자세한 옵션은 각 명령에 `--help` 를 붙여 확인하세요.

---

## 🚀 빠른 시작

### 새 프로젝트에 LOADSTAR 도입

```bash
cd /my/project

# 1. 메타데이터 디렉토리 초기화
loadstar init

# 2. 첫 WayPoint·Map은 AI에게 작성을 요청하거나 직접 .loadstar/WAYPOINT/ 에 추가
#    (자세한 절차는 openLoadstar README의 "AI 세션 진입 프롬프트" 참조)

# 3. 현재 상태 확인
loadstar show           # 전체 WayPoint 목록
loadstar show --recent  # 최근 수정순
loadstar show frontend  # "frontend" 키워드 필터

# 4. 깨진 참조 검증
loadstar validate
```

### 일상적인 메타 운영

```bash
# 진행 중인 작업 목록
loadstar todo list

# WayPoint STATUS 변경 후 동기화
loadstar todo sync

# 완료 이력 조회
loadstar todo history
loadstar todo history M://root/cli

# 변경 로그 검색
loadstar log 7d                    # 최근 7일
loadstar log cmd_show              # 키워드 필터
loadstar log 2d ISSUE              # 기간 + KIND 필터

# 사용자 결정이 필요한 미해결 질문 확인
loadstar question
loadstar question --with-resolved  # 결정 완료 항목까지 포함
```

### 메타 이벤트 직접 기록

```bash
loadstar log add W://root/cli/cmd_show MODIFIED "show 명령에 --recent 플래그 추가"
```

---

## 🧭 주소 체계

```
M://root/cli            →  .loadstar/MAP/root.cli.md
W://root/cli/cmd_show   →  .loadstar/WAYPOINT/root.cli.cmd_show.md
```

- **M (Map)**: WayPoint 묶음을 위한 인덱스 — STATUS 없음, 계층 경로만 표현. 하위 트리의 큰 의도를 담는 **GOAL** 슬롯을 가진다.
- **W (WayPoint)**: 작업의 최소 단위 — IDENTITY / CONNECTIONS / CODE_MAP / TODO / ISSUE 로 구성. 해당 WP가 달성해야 할 의도를 담는 **GOAL** 슬롯을 가진다.
- **D (Data WayPoint / dwp)**: 코드와 독립적으로 변화하는 데이터 아티팩트(설정, 참조 데이터)의 메타 레코드 — `.loadstar/DATA_WAYPOINT/`에 저장.

---

## 📂 디렉토리 구조

```
.loadstar/
├── MAP/              M:// 요소 (마크다운)
├── WAYPOINT/         W:// 요소 (마크다운)
├── DATA_WAYPOINT/    D:// 요소 — 데이터 아티팩트 메타데이터 (마크다운)
├── DECISIONS/        OPEN_QUESTIONS 결정 기록 (ADR)
├── COMMON/           프로젝트 설정
└── .clionly/         ⚠️ CLI 전용 — AI·사람 모두 직접 편집 금지
    ├── LOG/          변경 이력 로그
    └── TODO/         TODO_LIST·WP_SNAPSHOT (sync가 관리)
```

> `.clionly/` 직접 편집 시 LOG와 실제 메타 상태 간 정합성이 영구적으로 깨집니다.

---

## 🤖 AI 협업 워크플로우

1. **세션 시작 시** — AI는 `LOADSTAR_INIT.md` 와 SPEC을 로드하고, `loadstar show` / `loadstar todo list` / `loadstar question` 으로 현재 상태를 파악합니다.
2. **코드 수정 전** — 대상 WayPoint의 TODO 섹션에 `- [ ] 작업 내용` 을 등록합니다.
3. **수정 완료 후** — `- [x] YYYY-MM-DD 작업 내용` 으로 체크합니다.
4. **WayPoint 전체 완료** — STATUS를 `S_PRG → S_STB` 로 변경합니다.
5. **TODO 갱신** — `loadstar todo sync` 로 WP STATUS 기반 TODO_LIST를 자동 갱신합니다.
6. **검증** — 작업 종료 전 `loadstar validate` 로 깨진 참조 및 ORPHAN/DANGLING Data WayPoint가 없는지 확인합니다.

> Claude Code 환경에서는 PostToolUse Hook이 소스 편집 시 TODO 등록·갱신 리마인더를 자동 출력하도록 구성할 수 있습니다.

---

## 🧩 의존성

- [cobra](https://github.com/spf13/cobra) — CLI 프레임워크

표준 라이브러리 외 추가 의존성 최소화 정책.

---

## 🔗 관련 프로젝트

- 🌐 **[openLoadstar](https://github.com/openLoadstar/openLoadstar)** — 전체 생태계 안내
- 📖 **[spec](https://github.com/openLoadstar/spec)** — LOADSTAR 방법론 명세
- 🖥️ **[ui](https://github.com/openLoadstar/ui)** — Spring Boot + React 기반 Explorer UI
- 🔌 **[mcp](https://github.com/openLoadstar/mcp)** — Python MCP 서버 (Claude Desktop·Cursor 등 외부 AI 클라이언트 연동)

---

## 📮 기여 / 보안

- 🤝 **기여 가이드**: [openLoadstar/CONTRIBUTING.ko.md](https://github.com/openLoadstar/openLoadstar/blob/main/CONTRIBUTING.ko.md)
- 🔒 **보안 신고**: [openLoadstar/SECURITY.ko.md](https://github.com/openLoadstar/openLoadstar/blob/main/SECURITY.ko.md) — GitHub Security Advisories를 우선 사용해 주세요.
- 💬 **질문·아이디어**: [GitHub Discussions](https://github.com/openLoadstar/openLoadstar/discussions)

---

## 📄 License

[Apache License 2.0](./LICENSE)
