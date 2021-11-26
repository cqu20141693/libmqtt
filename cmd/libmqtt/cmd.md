### connect:

``` 
prod:
s: 172.20.245.147  
gk: gjpmAyDizqMCVivI gt: uG2q9gioJcpTV663 g2_sn g1_sn app: j4IzIBnqHCCbis76
```

#### real device

``` 
conn 172.30.203.21:1883 gateway_Sn yKbb7TlRuOMDOAqp G:2fw6oC2eVtDKraks 120

prod:
conn 172.20.245.147:1883 g1_sn gjpmAyDizqMCVivI G:uG2q9gioJcpTV663 120 
conn mqtt.chongctech.com:1883 g1_sn gjpmAyDizqMCVivI G:uG2q9gioJcpTV663 120
```

#### mirror

``` 

m-conn 172.30.203.21:1883 tPH6EZy6UbIxHxkg dVwEOrHihSRHtZTm http://172.30.203.22:11093/api/device/authenticate/config/getMirrorAuth 120

prod 
m-conn 172.20.245.147:1883 gjpmAyDizqMCVivI g1_sn http://172.20.245.154:11093/api/device/authenticate/config/getMirrorAuth?groupKey=gjpmAyDizqMCVivI&sn=g1_sn 120

```

#### signature

``` 
s-conn 172.30.203.21:1883 gateway_Sn yKbb7TlRuOMDOAqp SM3:2fw6oC2eVtDKraks 120
```

#### crypto

```  
s-conn 172.30.203.21:1883 gateway_Sn yKbb7TlRuOMDOAqp SM4:2fw6oC2eVtDKraks 120
```

#### subGroup

```
sg-conn 172.30.203.21:1883 tPH6EZy6UbIxHxkg 0sxsPsb5lBCEwXa0 http://172.30.203.22:11093/api/device/authenticate/config/getSubLoginAuthByGroupKey dVwEOrHihSRHtZTm,MQZddYPwBBepuXBx,WVZhKiRnDDzTyljo 120

```

#### app

```
conn 172.30.203.21:1883 0sxsPsb5lBCEwXa0 0sxsPsb5lBCEwXa0 A:gSMBUsjqOwdzHFWO 120
```

### publish:

```
p sys/data/string-channel/string#1#string#go-mqtt sys/data/int-channel/1#0#int#1 sys/data/long-channel/2#0#long#1635495018846 sys/data/float-channel/3#0#float#1.111 sys/data/double-channel/4#0#double#1.0111111 sys/data/struct-channel/6#0#json#{"name":"gow","age":25} 1
p sys/data/struct-channel/6#0#json#{"name":"gow","age":25} 1

```

#### proxy

##### data

```
p sys/data/yKbb7TlRuOMDOAqp/child2/struct-channel/6#0#json#{"name":"gow","age":25} 2 p
sys/data/yKbb7TlRuOMDOAqp/child1/struct-channel/6#0#json#{"name":"gow","age":25} 2

```

##### topology

```
p sys/topology/reportStatus#0#json#{"online":[{"groupKey":"yKbb7TlRuOMDOAqp","sn":"child2"}]} 1 
p sys/topology/fetch#0#json#{"needOtaAuthInfo":true} 1 p sys/topology/fetchDiff#0#json#{"needOtaAuthInfo":true,"useVersion":0} 1

```

##### shadow

```
p sys/shadow/fetchStatus#0#json#{"groupKey":"yKbb7TlRuOMDOAqp","sn":"gateway_Sn"} 1

```

### disconnect

## Mini

```
conn 172.30.203.22:1883 gateway_Sn yKbb7TlRuOMDOAqp G:2fw6oC2eVtDKraks 120

```