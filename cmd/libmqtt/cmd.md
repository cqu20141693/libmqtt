### topology

```
p sys/topology/reportStatus#0#json#{"online":[{"groupKey":"yKbb7TlRuOMDOAqp","sn":"child2"}]} 1
p sys/topology/fetch#0#json#{"needOtaAuthInfo":true} 1 
p sys/topology/fetchDiff#0#json#{"needOtaAuthInfo":true,"useVersion":0} 1

```

### shadow

```
p sys/shadow/fetchStatus#0#json#{} 1
p sys/shadow/yKbb7TlRuOMDOAqp/child1/fetchStatus#0#json#{} 1

prod
p sys/shadow/gjpmAyDizqMCVivI/child2/fetchStatus#0#json#{} 1
```

### test

#### prod

##### conn

``` 
conn 172.20.245.147:1883 g1_sn gjpmAyDizqMCVivI G:uG2q9gioJcpTV663 120 
conn mqtt.cc.com:1883 g1_sn gjpmAyDizqMCVivI G:uG2q9gioJcpTV663 120

conn mqtt.cc.com:1883 CC9912341234 7YaRYTy23fHQP7vV G:qUAzGNLrhYFSuJji 120

loginKey: F5PlN3JXKiwM4rHO sn:gateway1  gt:NXGX0SKqDDBTUHkQ
conn mqtt.cc.com:1883 gateway1 F5PlN3JXKiwM4rHO G:NXGX0SKqDDBTUHkQ 120
```

##### mirror

``` 
m-conn mqtt.cc.com:1883 F5PlN3JXKiwM4rHO gateway2 http://172.20.245.157:11093/api/device/authenticate/config/getMirrorAuth 120
```

##### app

``` 

```

##### push

``` 
p sys/data/string-channel/string#1#string#go-mqtt

child:
p sys/data/gjpmAyDizqMCVivI/child2/int-channel/1#0#int#10 1 
```

##### self subscribe

``` 
控制台添加订阅配置或者API添加
dest:
conn mqtt.cc.com:1883 gateway1 F5PlN3JXKiwM4rHO G:NXGX0SKqDDBTUHkQ 120

src:
conn mqtt.cc.com:1883 gateway2 F5PlN3JXKiwM4rHO G:NXGX0SKqDDBTUHkQ 120

```

##### child subscribe

```
控制台添加订阅配置或者API添加

p sys/topology/reportStatus#0#json#{"online":[{"groupKey":"F5PlN3JXKiwM4rHO","sn":"child1"}]} 1

```

#### dev

##### conn

``` 
conn 172.30.203.21:1883 gateway_Sn yKbb7TlRuOMDOAqp G:2fw6oC2eVtDKraks 120

```

##### signature

``` 
s-conn 172.30.203.21:1883 gateway_Sn yKbb7TlRuOMDOAqp SM3:2fw6oC2eVtDKraks 120
```

##### crypto

```  
s-conn 172.30.203.21:1883 gateway_Sn yKbb7TlRuOMDOAqp SM4:2fw6oC2eVtDKraks 120
```

##### push

``` 
p sys/data/string-channel/string#1#string#go-mqtt sys/data/int-channel/1#0#int#1 sys/data/long-channel/2#0#long#1635495018846 sys/data/float-channel/3#0#float#1.111 sys/data/double-channel/4#0#double#1.0111111 sys/data/struct-channel/6#0#json#{"name":"gow","age":25} 1

child:
p sys/data/yKbb7TlRuOMDOAqp/child2/int-channel/1#0#int#10 1
```

##### send cmd

``` 
控制台或者API
```

##### mirror

``` 
m-conn 172.30.203.21:1883 tPH6EZy6UbIxHxkg dVwEOrHihSRHtZTm http://172.30.203.22:11093/api/device/authenticate/config/getMirrorAuth 120
```

##### app

``` 
conn 172.30.203.21:1883 0sxsPsb5lBCEwXa0 0sxsPsb5lBCEwXa0 A:gSMBUsjqOwdzHFWO 120
```

##### subscribe group

``` 
sg-conn 172.30.203.21:1883 tPH6EZy6UbIxHxkg 0sxsPsb5lBCEwXa0 http://172.30.203.22:11093/api/device/authenticate/config/getSubLoginAuthByGroupKey dVwEOrHihSRHtZTm,MQZddYPwBBepuXBx,WVZhKiRnDDzTyljo 120
```

##### self subscribe

``` 
控制台添加订阅配置或者API添加
dest:
conn 172.30.203.21:1883 gateway_Sn yKbb7TlRuOMDOAqp G:2fw6oC2eVtDKraks 120

src:
conn 47.108.93.28:1883 gateway1 v7EZ22WxxpCkzRnw G:ODmP6sHYd2RigoZi 120

```

##### child subscribe

```
控制台添加订阅配置或者API添加

p sys/topology/reportStatus#0#json#{"online":[{"groupKey":"yKbb7TlRuOMDOAqp","sn":"child2"}]} 1

```