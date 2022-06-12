git status
git add .
echo "You will be saving the following changes:"
echo "Please enter a commit message:"
read commit_message
if [ -z "$commit_message" ]; then
  commit_message="Saved changes"
fi
git commit -m "$commit_message"
git push