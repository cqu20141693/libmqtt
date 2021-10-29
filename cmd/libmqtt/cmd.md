###connect:
#### real device
conn 172.30.203.21:1883 gateway_Sn yKbb7TlRuOMDOAqp G:2fw6oC2eVtDKraks 120
#### mirror
m-conn 172.30.203.21:1883 http://172.30.203.22:11093/api/device/authenticate/config/getMirrorAuth?groupKey=tPH6EZy6UbIxHxkg&sn=dVwEOrHihSRHtZTm 120
publish:
p sys/data/string-channel/string#1#string#go-mqtt sys/data/int-channel/1#0#int#1 sys/data/long-channel/2#0#long#1635495018846 sys/data/float-channel/3#0#float#1.111 sys/data/double-channel/4#0#double#1.0111111 sys/data/struct-channel/6#0#json#{"name":"gow","age":25} 1
p sys/data/struct-channel/6#0#json#{"name":"gow","age":25} 2

