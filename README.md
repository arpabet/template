# template

![build workflow](https://github.com/sprintframework/template/actions/workflows/build.yaml/badge.svg)

Template Server Application on Golang with Web Interface on Vue.Js.
Application ships as a single executable file for easy installation.

How to run UI locally
```
cd webapp
npm install
npm run dev
```

How to build App
```
go generate
make
```

How to run App
```
./template run
```

Properties:
```
access.token.minutes 20 by default
refresh.token.hours   24 by default
mail.sender   email like noreply@domainname
mail.support  email like support@domainname 
jwt.secret.key   token
mailgun.key from mailgun dashboard
```

