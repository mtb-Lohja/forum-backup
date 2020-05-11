# MTB Lohja Forum backup cron job

**Note: This is legacy; replaced with gcloud snapshot schedules**

Runs disk snapshot daily by leveraging appEngine cron capabilities.
Deploy the app with

    gcloud app deploy app.yaml
    gcloud app deploy cron.yaml
