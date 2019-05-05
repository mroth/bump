package main

import "testing"

func Test_parseGithubRemote(t *testing.T) {
	type args struct {
		remoteURL string
	}
	tests := []struct {
		name      string
		remoteURL string
		wantOwner string
		wantRepo  string
		wantOk    bool
	}{
		{
			name:      "GitHub_HTTPS",
			remoteURL: "https://github.com/mroth/bump.git",
			wantOwner: "mroth",
			wantRepo:  "bump",
			wantOk:    true,
		},
		{
			name:      "GitHub_SSH",
			remoteURL: "git@github.com:mroth/bump.git",
			wantOwner: "mroth",
			wantRepo:  "bump",
			wantOk:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOwner, gotRepo, gotOk := parseGithubRemote(tt.remoteURL)
			if gotOwner != tt.wantOwner {
				t.Errorf("parseGithubRemote() gotOwner = %v, want %v", gotOwner, tt.wantOwner)
			}
			if gotRepo != tt.wantRepo {
				t.Errorf("parseGithubRemote() gotRepo = %v, want %v", gotRepo, tt.wantRepo)
			}
			if gotOk != tt.wantOk {
				t.Errorf("parseGithubRemote() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_githubRepoDetect(t *testing.T) {
	owner, repo, err := githubRepoDetect(".")
	if err != nil {
		t.Fatal(err)
	}
	expect := "mroth/bump"
	actual := owner + "/" + repo
	if expect != actual {
		t.Errorf("want %v got %v", expect, actual)
	}
}

func Benchmark_detectRemoteURL_GoGit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_detectRemoteURL_GoGit(".")
	}
}

func Benchmark_detectRemoteURL_LocalGit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_detectRemoteURL_LocalGit(".")
	}
}

func Benchmark_parseGithubRemote(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseGithubRemote("https://github.com/mroth/bump.git")
	}
}
