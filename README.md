# MTB Lohja Forum backup cron job

Runs disk snapshot daily by leveraging appEngine cron capabilities. 
Deploy the app with

    gcloud app deploy app.yaml
    gcloud app deploy cron.yaml
