RFC 3339 is used for datetime for both frontend and backend. YYYY-MM-DDTHH:mm:ss.sssZ
Measurements are send every second, that's around 100Mb of data per month.

Use it to exclude the data file from tracking by git:
git update-index --assume-unchanged backend/sample_climate.txt

## TODO
- Add websocket to send data to frontend
- Add charts to frontend
- Change sample_climate.txt for sample and some static file that is created
