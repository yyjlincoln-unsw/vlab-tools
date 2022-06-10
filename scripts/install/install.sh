echo "Installing $1..."
mkdir -p ~/bin
# This is also a shadow - it will look into
# the discoverable directories and find the 
# file during runtime.
curl -L -o ~/bin/$1 https://yyj.app/shadow
chmod +x ~/bin/$1