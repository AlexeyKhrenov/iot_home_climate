import network
import urequests
from dht import DHT11
from machine import Pin
import time
import gc

wlan = network.WLAN(network.STA_IF)
wlan.active(True)

ssid = '***'
password = '***'.replace(" ", "")
ip_address = '***' // address in format http://x.x.x.x:80/someapi
wlan.connect(ssid, password)

while not wlan.isconnected():
    print("Connecting to Wi-Fi...")
    time.sleep(1)
print("Connected to Wi-Fi:", wlan.ifconfig())

sensor = DHT11(Pin(16, Pin.OUT, Pin.PULL_DOWN))
 
while True:
    try:
        sensor.measure()
        temp = sensor.temperature()
        humidity = sensor.humidity()
        time.sleep(2)
        
        data = 'temp=' + str(temp) + '&humidity=' + str(humidity)

        headers = {'Content-Type': 'application/x-www-form-urlencoded'}
        response = urequests.post(ip_address, data=data, headers=headers)
        response.close()
        gc.collect()
    except Exception as e:
        print("Error:", e)

