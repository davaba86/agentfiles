# Homebrew Distribution

`agentfiles` should be distributed through a separate public Homebrew tap repository:

```text
github.com/davaba86/homebrew-tap
```

Users can install from that tap with:

```sh
brew install davaba86/tap/agentfiles
```

## Release flow

Create a GitHub release from a tag in the public source repository:

```sh
gh release create v0.1.0 \
  --title "v0.1.0" \
  --notes "Initial release of agentfiles."
```

GitHub automatically exposes the source archive at:

```text
https://github.com/davaba86/agentfiles/archive/refs/tags/v0.1.0.tar.gz
```

Compute the archive checksum:

```sh
curl -L https://github.com/davaba86/agentfiles/archive/refs/tags/v0.1.0.tar.gz | shasum -a 256
```

## Formula

Add this formula to `davaba86/homebrew-tap` as `Formula/agentfiles.rb`.

```ruby
class Agentfiles < Formula
  desc "Standardize AI coding-agent instruction files across repositories"
  homepage "https://github.com/davaba86/agentfiles"
  url "https://github.com/davaba86/agentfiles/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "PUT_RELEASE_TARBALL_SHA256_HERE"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-X main.version=#{version}"), "./cmd/agentfiles"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/agentfiles version")
  end
end
```

## Public repository notes

- Do not commit personal filesystem paths in public docs.
- Do not commit GitHub tokens or Homebrew tap credentials.
- If release automation later updates `homebrew-tap`, store the token as a GitHub Actions secret.
- Keep the formula in the tap repository, not in this source repository.
