echo Fetching files...
cd /tmp
rm install-cloud-autotest.py > /dev/null 2>&1
curl -O https://static.yyjlincoln.com/scripts/install-cloud-autotest.py
python3 install-cloud-autotest.py --all
rm install-cloud-autotest.py
echo Done.