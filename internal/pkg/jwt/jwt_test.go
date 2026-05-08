package jwt

import (
	"testing"
	"time"
)

func mustSigner(t *testing.T, access, refresh time.Duration) *Signer {
	t.Helper()
	s, err := NewSigner(Config{
		Secret:     "test-secret-at-least-8",
		AccessTTL:  access,
		RefreshTTL: refresh,
		Issuer:     "myblog-test",
	})
	if err != nil {
		t.Fatalf("new signer: %v", err)
	}
	return s
}

func TestSigner_IssueAndVerifyAccess(t *testing.T) {
	s := mustSigner(t, time.Hour, 24*time.Hour)

	access, refresh, exp, err := s.Issue(7, "alice", "admin")
	if err != nil {
		t.Fatalf("issue: %v", err)
	}
	if access == "" || refresh == "" {
		t.Fatal("tokens should not be empty")
	}
	if time.Until(exp) <= 0 {
		t.Fatal("access token should not be expired")
	}

	claims, err := s.Verify(access, TypeAccess)
	if err != nil {
		t.Fatalf("verify access: %v", err)
	}
	if claims.UserID != 7 || claims.Username != "alice" || claims.Role != "admin" {
		t.Errorf("claims mismatch: %+v", claims)
	}
}

func TestSigner_Verify_WrongType(t *testing.T) {
	s := mustSigner(t, time.Hour, 24*time.Hour)
	access, refresh, _, _ := s.Issue(1, "u", "author")

	if _, err := s.Verify(refresh, TypeAccess); err == nil {
		t.Errorf("expect ErrWrongType when treating refresh as access, got nil")
	}
	if _, err := s.Verify(access, TypeRefresh); err == nil {
		t.Errorf("expect ErrWrongType when treating access as refresh, got nil")
	}
}

func TestSigner_Verify_Expired(t *testing.T) {
	s := mustSigner(t, 50*time.Millisecond, time.Hour)
	access, _, _, _ := s.Issue(1, "u", "author")
	time.Sleep(1100 * time.Millisecond) // jwt exp 以秒为粒度,sleep 略超 1s

	_, err := s.Verify(access, TypeAccess)
	if err != ErrExpiredToken {
		t.Errorf("expect ErrExpiredToken, got %v", err)
	}
}

func TestSigner_Verify_BadSignature(t *testing.T) {
	s := mustSigner(t, time.Hour, time.Hour)
	access, _, _, _ := s.Issue(1, "u", "author")

	// 换一把 secret 校验
	other, _ := NewSigner(Config{Secret: "another-secret-key", Issuer: "myblog-test"})
	if _, err := other.Verify(access, TypeAccess); err == nil {
		t.Errorf("expect invalid signature error, got nil")
	}
}
