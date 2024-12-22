curl -X POST http://localhost:8080/climate \
	-d "temp=23&humidity=53" \
	-w "Response : %{http_code}\n"

