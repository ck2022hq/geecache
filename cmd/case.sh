echo ">>> start test"

curl "http://localhost:9999/geecache?group=ha&key=Tom" &
curl "http://localhost:9999/geecache?group=hb&key=Tom" &
curl "http://localhost:9999/geecache?group=hc&key=Tom" &

curl "http://localhost:9999/geecache?group=ha&key=Jack" &
curl "http://localhost:9999/geecache?group=hb&key=Jack" &
curl "http://localhost:9999/geecache?group=hc&key=Jack" &

curl "http://localhost:9999/geecache?group=ha&key=Sam" &
curl "http://localhost:9999/geecache?group=hb&key=Sam" &
curl "http://localhost:9999/geecache?group=hc&key=Sam" &

curl "http://localhost:9999/geecache?group=ha&key=Ss" &
curl "http://localhost:9999/geecache?group=hb&key=Ss" &
curl "http://localhost:9999/geecache?group=hc&key=Ss" &

echo ">>> end test"