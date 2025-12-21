# ADR-002: CI (Continuous Integration) 전략

## 상태

Accepted

## 컨텍스트

코드 품질 보장 및 자동화된 테스트 실행을 위한 CI 시스템이 필요하다.
다양한 CI/CD 도구 옵션이 존재한다:

* Tekton
* GitHub Actions
* Jenkins
* GitLab CI
* CircleCI

## 결정

**Tekton을 CI 도구로 사용한다.**

### 선택 이유

1. **Kubernetes 네이티브**: Kubernetes 환경에서 네이티브하게 동작하여 인프라 일관성 확보
2. **클라우드 애그노스틱**: Kubernetes가 동작하는 어느 환경에서든 사용 가능
3. **확장성**: Kubernetes의 리소스 모델을 활용하여 높은 확장성 제공
4. **표준 기반**: Cloud Native Computing Foundation (CNCF) 프로젝트로 표준 기반
5. **재사용성**: Task와 Pipeline의 재사용이 용이하여 여러 프로젝트에서 활용 가능

### CI 파이프라인 구성

다음 단계들을 포함한다:

1. **코드 체크아웃**: GitHub에서 소스 코드 가져오기
2. **정적 분석**: golangci-lint를 통한 코드 품질 검사
3. **테스트 실행**: `go test`를 통한 단위 테스트 및 레이스 컨디션 검사
4. **빌드**: Docker 이미지 빌드
5. **이미지 푸시**: 컨테이너 레지스트리에 이미지 푸시

## 결과

### 긍정적 영향

* Kubernetes 환경과의 완벽한 통합
* 인프라와 CI/CD 인프라의 일관성 확보
* 확장 가능한 파이프라인 구성
* 다른 Kubernetes 기반 프로젝트와의 Task/Pipeline 재사용 가능

### 부정적 영향

* Tekton 학습 곡선 존재
* Kubernetes 클러스터가 필요하여 초기 설정 복잡도 증가
* GitHub Actions 대비 설정 복잡도 상대적으로 높음
