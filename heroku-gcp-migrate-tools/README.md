This is a small utility that reads your heroku configs and convert them 
into GCP secrets. They can be referred to in your k8s yaml, 
to save some of your typeworks when migrating.

Prerequisite
------------

You must have heroku toolbelt installed and connected to your heroku 
account.

Download
--------

Please find them in the dist folder. At the moment, only darwin/amd64 is 
provided.

Running
-------

```
heroku-gcp-migrate-tool <heroku app name>
```

[sohoffice](https://medium.com/sohoffice), Happy coding ~
