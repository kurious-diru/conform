// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

//nolint:testpackage
package commit

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"testing"

	"github.com/talos-systems/conform/internal/policy"
)

func RemoveAll(dir string) {
	if err := os.RemoveAll(dir); err != nil {
		log.Fatal(err)
	}
}

//nolint:gocognit
func TestConventionalCommitPolicy(t *testing.T) {
	//nolint:govet
	type testDesc struct {
		Name         string
		CreateCommit func() error
		ExpectValid  bool
	}

	for _, test := range []testDesc{
		{
			Name:         "Valid",
			CreateCommit: createValidCommit,
			ExpectValid:  true,
		},
		{
			Name:         "Breaking",
			CreateCommit: createBreakingCommit,
			ExpectValid:  true,
		},
		{
			Name:         "InvalidBreaking",
			CreateCommit: createInvalidBreakingCommit,
			ExpectValid:  false,
		},
		{
			Name:         "Invalid",
			CreateCommit: createInvalidCommit,
			ExpectValid:  false,
		},
		{
			Name:         "Empty",
			CreateCommit: createEmptyCommit,
			ExpectValid:  false,
		},
	} {
		func(test testDesc) {
			t.Run(test.Name, func(tt *testing.T) {
				dir, err := ioutil.TempDir("", "test")
				if err != nil {
					log.Fatal(err)
				}
				defer RemoveAll(dir)
				err = os.Chdir(dir)
				if err != nil {
					tt.Error(err)
				}
				err = initRepo()
				if err != nil {
					tt.Error(err)
				}

				err = test.CreateCommit()
				if err != nil {
					tt.Error(err)
				}
				report, err := runCompliance()
				if err != nil {
					t.Error(err)
				}

				if test.ExpectValid {
					if !report.Valid() {
						tt.Error("Report is invalid with valid conventional commit")
					}
				} else {
					if report.Valid() {
						tt.Error("Report is valid with invalid conventional commit")
					}
				}
			})
		}(test)
	}
}

func TestValidateDCO(t *testing.T) {
	type testDesc struct {
		Name          string
		CommitMessage string
		ExpectValid   bool
	}

	for _, test := range []testDesc{
		{
			Name:          "Valid DCO",
			CommitMessage: "something nice\n\nSigned-off-by: Foo Bar <foobar@example.org>\n\n",
			ExpectValid:   true,
		},
		{
			Name:          "Valid DCO with CRLF",
			CommitMessage: "something nice\r\n\r\nSigned-off-by: Foo Bar <foobar@example.org>\r\n\r\n",
			ExpectValid:   true,
		},
		{
			Name:          "No DCO",
			CommitMessage: "something nice\n\nnot signed\n",
			ExpectValid:   false,
		},
	} {
		// Fixes scopelint error.
		test := test
		t.Run(test.Name, func(tt *testing.T) {
			var report policy.Report
			c := Commit{msg: test.CommitMessage}
			report.AddCheck(c.ValidateDCO())

			if test.ExpectValid {
				if !report.Valid() {
					tt.Error("Report is invalid with valid DCP")
				}
			} else {
				if report.Valid() {
					tt.Error("Report is valid with invalid DCO")
				}
			}
		})
	}
}

func TestValidConventionalCommitPolicy(t *testing.T) {
	dir, err := ioutil.TempDir("", "test")
	if err != nil {
		log.Fatal(err)
	}

	defer RemoveAll(dir)

	err = os.Chdir(dir)
	if err != nil {
		t.Error(err)
	}

	err = initRepo()
	if err != nil {
		t.Error(err)
	}

	err = createValidCommit()
	if err != nil {
		t.Error(err)
	}

	report, err := runCompliance()
	if err != nil {
		t.Error(err)
	}

	if !report.Valid() {
		t.Errorf("Report is invalid with valid conventional commit")
	}
}

func TestInvalidConventionalCommitPolicy(t *testing.T) {
	dir, err := ioutil.TempDir("", "test")
	if err != nil {
		log.Fatal(err)
	}

	defer RemoveAll(dir)

	err = os.Chdir(dir)
	if err != nil {
		t.Error(err)
	}

	err = initRepo()
	if err != nil {
		t.Error(err)
	}

	err = createInvalidCommit()
	if err != nil {
		t.Error(err)
	}

	report, err := runCompliance()
	if err != nil {
		t.Error(err)
	}

	if report.Valid() {
		t.Errorf("Report is valid with invalid conventional commit")
	}
}

func TestEmptyConventionalCommitPolicy(t *testing.T) {
	dir, err := ioutil.TempDir("", "test")
	if err != nil {
		log.Fatal(err)
	}

	defer RemoveAll(dir)

	err = os.Chdir(dir)
	if err != nil {
		t.Error(err)
	}

	err = initRepo()
	if err != nil {
		t.Error(err)
	}

	err = createEmptyCommit()
	if err != nil {
		t.Error(err)
	}

	report, err := runCompliance()
	if err != nil {
		t.Error(err)
	}

	if report.Valid() {
		t.Error("Report is valid with invalid conventional commit")
	}
}

func TestValidConventionalCommitPolicyRegex(t *testing.T) {
	dir, err := ioutil.TempDir("", "test")
	if err != nil {
		log.Fatal(err)
	}

	defer RemoveAll(dir)

	err = os.Chdir(dir)
	if err != nil {
		t.Error(err)
	}

	err = initRepo()
	if err != nil {
		t.Error(err)
	}

	err = createValidCommitRegex()
	if err != nil {
		t.Error(err)
	}

	report, err := runCompliance()
	if err != nil {
		t.Error(err)
	}

	if !report.Valid() {
		t.Error("Report is invalid with valid conventional commit")
	}
}

func TestInvalidConventionalCommitPolicyRegex(t *testing.T) {
	dir, err := ioutil.TempDir("", "test")
	if err != nil {
		log.Fatal(err)
	}

	defer RemoveAll(dir)

	err = os.Chdir(dir)
	if err != nil {
		t.Error(err)
	}

	err = initRepo()
	if err != nil {
		t.Error(err)
	}

	err = createInvalidCommitRegex()
	if err != nil {
		t.Error(err)
	}

	report, err := runCompliance()
	if err != nil {
		t.Error(err)
	}

	if report.Valid() {
		t.Error("Report is valid with invalid conventional commit")
	}
}

func runCompliance() (*policy.Report, error) {
	c := &Commit{
		Conventional: &Conventional{
			Types:  []string{"type"},
			Scopes: []string{"scope", "^valid"},
		},
	}

	return c.Compliance(&policy.Options{})
}

func initRepo() error {
	_, err := exec.Command("git", "init").Output()
	if err != nil {
		return err
	}

	_, err = exec.Command("touch", "test").Output()
	if err != nil {
		return err
	}

	_, err = exec.Command("git", "add", "test").Output()

	return err
}

func createValidCommit() error {
	_, err := exec.Command("git", "-c", "user.name='test'", "-c", "user.email='test@talos-systems.io'", "commit", "-m", "type(scope): description").Output()

	return err
}

func createBreakingCommit() error {
	_, err := exec.Command("git", "-c", "user.name='test'", "-c", "user.email='test@talos-systems.io'", "commit", "-m", "feat!: description").Output()

	return err
}

func createInvalidBreakingCommit() error {
	_, err := exec.Command("git", "-c", "user.name='test'", "-c", "user.email='test@talos-systems.io'", "commit", "-m", "feat$: description").Output()

	return err
}

func createInvalidCommit() error {
	_, err := exec.Command("git", "-c", "user.name='test'", "-c", "user.email='test@talos-systems.io'", "commit", "-m", "invalid commit").Output()

	return err
}

func createEmptyCommit() error {
	_, err := exec.Command("git", "-c", "user.name='test'", "-c", "user.email='test@talos-systems.io'", "commit", "--allow-empty-message", "-m", "").Output()

	return err
}

func createValidCommitRegex() error {
	_, err := exec.Command("git", "-c", "user.name='test'", "-c", "user.email='test@talos-systems.io'", "commit", "-m", "type(valid-1): description").Output()

	return err
}

func createInvalidCommitRegex() error {
	_, err := exec.Command("git", "-c", "user.name='test'", "-c", "user.email='test@talos-systems.io'", "commit", "-m", "type(invalid-1): description").Output()

	return err
}
