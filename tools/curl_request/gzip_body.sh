echo '{"URL":"https://ya.ru/search/?text=curl+send+gzip+body&lr=47"}' | gzip > body.gz
curl -v -i http://localhost:8080/api/shorten -H'Content-Encoding: gzip' -H'Content-type: application/x-gzip' --data-binary @body.gz
rm body.gz
echo "\n"