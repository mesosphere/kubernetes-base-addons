#!/bin/bash
# NOTE: this script was originally created to be used by dispatch CI runs
set -euxo pipefail

if test -z "git config user.email"; then
    git config user.email "ci@mesosphere.com";
fi

if test -z "git config user.name"; then
    git config user.name "CI";
fi

# Setup git
# Replace https://github.com/ with "git@github.com:" in ~/.gitconfig

git config --global url.git@github.com:.insteadOf https://github.com/

# Steps to make sure go mod will download from private git repositories.

# add SSH_KEY to the ssh-agent
# SSH_KEY_BASE64 is provided by Dispatch.

eval "$(ssh-agent -s)";
mkdir /root/.ssh;
echo $SSH_KEY_BASE64 |  tr -d "[:space:]" | base64 -d | install -b -m 600 /dev/stdin /root/.ssh/id_rsa

ssh-add /root/.ssh/id_rsa;

# trust github.com
ssh-keyscan github.com >> /root/.ssh/known_hosts;
