## pirogom/ddns/ddns_cli
====

Simple DDNS cli client

cli args -

-name [domain name]

-server [ddns server ip or domain ex:ddns.example.com:8080]

-cmd [Update or Delete]

-ip [ip for domain regist]

ex)

ddns_cli -name test.example.com -ip 192.168.0.1 -server ddns.example.com:8080 -cmd update

ddns_cli -name test.example.com -ip 192.168.0.1 -server ddns.example.com:8080 -cmd delete
