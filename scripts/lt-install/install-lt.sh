echo Fetching files...
cd /tmp
rm install-lt.py > /dev/null 2>&1
curl -O https://static.yyjlincoln.com/scripts/install-lt.py
python3 install-lt.py --all
rm install-lt.py
echo Done.