# GitHub Actions workflows

The playbook generates these files into `.github/workflows/`. They are kept
here because the token that published this demo cannot write workflow files.
To enable CI, move them back:

    git mv ci/github-workflows/*.yml .github/workflows/
