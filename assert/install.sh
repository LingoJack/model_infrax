#!/bin/zsh
echo 'Setting environment...'
cp ~/.zshrc ~/.zshrc.bak
grep -v "export PATH=\"/Applications/jen/" ~/.zshrc > ~/.zshrc.tmp
echo 'export PATH="/Applications/jen/"' >> ~/.zshrc.tmp
mv ~/.zshrc.tmp ~/.zshrc
source ~/.zshrc
echo '✅ Environment done (original zshrc file backup at ~/.zshrc.bak)'