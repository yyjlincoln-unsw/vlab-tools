echo "Getting the latest changes..."
git pull --rebase
git add .
git status
echo "You will be saving the above changes:"
echo "Please enter a commit message:"
read commit_message
if [ -z "$commit_message" ]; then
  commit_message="Saved changes"
fi
git commit -m "$commit_message"
git push