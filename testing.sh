HOST=$1
curl $HOST/api/account/register -d '{"username": "john", "password": "john123456"}'
curl -X POST $HOST/api/account/token -d '{"username": "john", "password": "john123456"}' | tee token.json
curl -H "Authorization: Bearer `jq -r .token token.json`" $HOST/api/account/sessions
rm cookies.txt
curl -c cookies.txt $HOST/api/account/login -d '{"username": "john", "password": "john123456"}'
curl -b cookies.txt $HOST/api/account/sessions
