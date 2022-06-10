mkdir -p ~/bin
curl -L -o ~/bin/lt https://yyj.app/shadow
chmod +x ~/bin/lt
echo "Adding to PATH..."
./bin/lt post-install --all
