test:
	# GET
	curl -i http://localhost:3000/healthz
	curl -i "http://localhost:3000/users?q=Ada"
	curl -i http://localhost:3000/users/1

	# POST
	curl -i -X POST http://localhost:3000/users -H "Content-Type: application/json" -H "Authorization: Bearer JWT" -d '{"name":"Charlie"}'

	# PUT
	curl -i -X PUT http://localhost:3000/users/1 -H "Content-Type: application/json" -d '{"name":"Alice Updated"}'

	# PATCH
	curl -i -X PATCH http://localhost:3000/users/1 -H "Content-Type: application/json" -d '{"name":"Alice Patched"}'

	# DELETE
	curl -i -X DELETE http://localhost:3000/users/1

	# HEAD
	curl -I http://localhost:3000/ping

	# GET
	curl -i "http://localhost:3000/echo?foo=bar"

	# POST
	curl -i -X POST http://localhost:3000/echo \
	  -H "Content-Type: application/json" \
	  -d '{"hello":"world"}'

	# PUT
	curl -i -X PUT http://localhost:3000/echo -d "x=1"

	# PATCH
	curl -i -X PATCH http://localhost:3000/echo -d "y=2"

	# DELETE
	curl -i -X DELETE http://localhost:3000/echo

	# OPTIONS
	curl -i -X OPTIONS http://localhost:3000/echo

	# HEAD
	curl -I http://localhost:3000/echo
