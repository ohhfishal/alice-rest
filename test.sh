

echo "POST {user}"
curl -X POST localhost:8000/api/v1/event/test -H "Content-Type: application/json" --data '{"description": "foo"}'
echo ""

echo "GET {user}"
curl -X GET localhost:8000/api/v1/event/test
echo ""

echo "GET {user}/{id}"
curl -X GET localhost:8000/api/v1/event/test/0 
echo ""

echo "PATCH {user}/{id}"
curl -X PATCH localhost:8000/api/v1/event/test/0 -H "Content-Type: application/json" --data '{"description": "updated"}'
echo ""

echo "GET {user}/{id}"
curl -X GET localhost:8000/api/v1/event/test/0 
echo ""

echo "DELETE {user}/{id}"
curl -X DELETE localhost:8000/api/v1/event/test/0 
echo ""
