include setup/dev/Makefile

generate-certificate:
	openssl req -x509 -newkey rsa:2048 -keyout storage/private_key.pem -out storage/certificate.pem -sha256 -days 365 -nodes

generate-public:
	openssl rsa -in storage/private_key.pem -pubout -out storage/public_key.pem

authorize-auth:
	sudo chmod 777 -R storage