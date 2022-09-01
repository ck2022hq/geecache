echo ">>> start test"

curl "http://localhost:8001/geecache?group=ha&key=Tom" &
curl "http://localhost:8001/geecache?group=ha&key=Tom" &
curl "http://localhost:8001/geecache?group=ha&key=Tom" &

echo ">>> end test"
