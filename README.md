## Beanstalk Printer Queue Worker

Copy .env.example as .env and if need be replace localhost with your Beanstalkd Server IP.

Also set the REMOTE_FILE to the url where file will be downloaded from.

### Build Binary
```
cd src && go build -o ../bin/SeminarPrinter
```