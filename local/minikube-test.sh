#!/bin/bash

set -e
localport=8080
typename=svc/birthday-app
remoteport=8080
kubectl port-forward -n revolut $typename $localport:$remoteport > /dev/null 2>&1 &
pid=$!
trap '{
    kill $pid
}' EXIT
while ! nc -vz localhost $localport > /dev/null 2>&1
do
    sleep 0.1
done

echo -n 'user not found - '
curl -w '%{http_code}' -s -o /dev/null http://localhost:$localport/hello/nouser | (grep 404 > /dev/null && echo "OK") || (echo "FAIL" && exit 1)


today=$(date +%Y-%m-%d)
echo -n "put today $today - "
curl -w '%{http_code}' -H 'Content-Type: application/json' -s -o /dev/null -X PUT http://localhost:$localport/hello/admin -d '{"dateOfBirth": "'${today}'"}'  | (grep 400 > /dev/null && echo "OK") || (echo "FAIL" && exit 1)

tomorrow=$(date -v +1d +%Y-%m-%d)
echo -n "put tomorrow $tomorrow - "
curl -w '%{http_code}' -H 'Content-Type: application/json' -s -o /dev/null -X PUT http://localhost:$localport/hello/admin -d '{"dateOfBirth": "'${tomorrow}'"}'  | (grep 400 > /dev/null && echo "OK") || (echo "FAIL" && exit 1)

yesterday=$(date -v -20y -v +40d +%Y-%m-%d)
echo -n "put valid not today $yesterday - "
curl -w '%{http_code}' -H 'Content-Type: application/json' -s -o /dev/null -X PUT http://localhost:$localport/hello/admin -d '{"dateOfBirth": "'${yesterday}'"}'  | (grep 204 > /dev/null && echo "OK") || (echo "FAIL" && exit 1)

echo -n 'get valid not today - '
curl -s http://localhost:$localport/hello/admin | (grep 'Your birthday is in 40 days' > /dev/null && echo "OK") || (echo "FAIL" && exit 1)

today=$(date -v -20y +%Y-%m-%d)
echo -n "put today in the past $today - "
curl -w '%{http_code}' -H 'Content-Type: application/json' -s -o /dev/null -X PUT http://localhost:$localport/hello/admin -d '{"dateOfBirth": "'${today}'"}'  | (grep 204 > /dev/null && echo "OK") || (echo "FAIL" && exit 1)

echo -n 'get today in the past - '
curl -s http://localhost:$localport/hello/admin | (grep 'Happy birthday!' > /dev/null && echo "OK") || (echo "FAIL" && exit 1)
