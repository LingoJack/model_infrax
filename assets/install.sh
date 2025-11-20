#!/bin/zsh
echo 'Setting environment...'
cp ~/.zshrc ~/.zshrc.bak

grep -v "export PATH=\"/Applications/model_infrax/\"" ~/.zshrc > ~/.zshrc.tmp

echo 'export PATH="/Applications/model_infrax/:$PATH"' >> ~/.zshrc.tmp

mv ~/.zshrc.tmp ~/.zshrc
source ~/.zshrc
echo '✅ Environment done (original zshrc file backup at ~/.zshrc.bak)'
echo '📝 model_infrax directory has been added to PATH while preserving existing PATH entries'

chmod +x /Applications/model_infrax/model_infrax
echo '✅ model_infrax init successfully'

chmod +x /Applications/model_infrax/jcode
echo '✅ jcode init successfully'