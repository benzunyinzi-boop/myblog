package password

import "testing"

func TestHashAndVerify(t *testing.T) {
	h, err := Hash("s3cret-pw")
	if err != nil {
		t.Fatalf("hash: %v", err)
	}
	if !Verify(h, "s3cret-pw") {
		t.Errorf("verify with correct password should succeed")
	}
	if Verify(h, "wrong") {
		t.Errorf("verify with wrong password should fail")
	}
}

func TestHash_EmptyRejected(t *testing.T) {
	if _, err := Hash(""); err == nil {
		t.Errorf("empty password should be rejected")
	}
}
