# Installation Guide #

## Preconditions ##

This installation guide assumes that you have the Cloudfoundry CLI tool
installed and have already used it to push an application to Bluemix.

## Create Cloudant service on Bluemix ##

Use the Bluemix service catalog to create a Cloudant service instance and call
it **shsp-cloudant** (currently the service will be looked up in **VCAP_SERVICES**
using this name, it will be configurable in the near future).

## Create a Slack application ##

The focus of this integration is being able to host it yourself. There is no
central endpoint that could be used for the Slack Button installation
mechanism. Because of this you need to configure this integration yourself. Go
to https://api.slack.com/apps and click on **Create new app**. As Development
Slack Team use the one you want this integration to be installed for.

## Configure permissions and tokens ##
Go to **Features/OAuth & Permissions** and
add the permission scope **users:read**. This is needed to show the real names
of the voters in the poll results. Go to the top of the page and click on
**Install app to team** and then **Authorize**. Make a note somewhere of the
generated OAuth-Token. You will need it while creating the Cloudfoundry
manifest. Then go to **Settings / Basic Information** and make a note of the
**Verification Token**. This will be needed for the manifest, too.

## Clone the repo ##

Clone the git repo for this application (master will always be the current
stable release)

## Create Cloudfoundry manifest ##

To deploy the integration to Bluemix you need to create a manifest.yml file in
the root directory of the repo clone. You can use the included
manifest.yml.template as a guideline. The following things need to be changed in
it:

- You need to change the subdomain where the
application is hosted (**host**) 
- You need to set the token that Slack uses for
indentification of its requests (**env/SLACK_TOKEN**) 
- You need to set the OAuth
token the application uses for access to the Slack API (**env/SLACK_OAUTH_TOKEN**)

The other settings needn't to be changes.

## Push the application using cf push ##

Ensure that the manifest.yml you created in the last step is in the root dir
of the repo clone. If your CF CLI is configured correctly and you are logged
in into the CF endpoint, you can push the application by simply doing

	cf push

on the commandline while being in the root dir of the cloned repo.


## Set up slash commands ##
Go to the application on the Slack API page. Go to **Features/Slash Commands**.
Create the following two slash commands:

### New Poll ###
- Command: /poll
- Request URL: <Your CF application url>/newpoll
- Description: Create a new poll
- Usage hint: "A Question" OptionA "Option B" :cake:

### New anonymous poll ###
- Command: /pollanon
- Request URL: <Your CF application url>/newpollanon
- Description: Create a new anonymous poll
- Usage hint: "A Question" OptionA "Option B" :cake:

## Set up interactive messages ##
Go to the application on the Slack API page. Go to **Features/Interactive messages**.
Click on **Enable interactive messages** and set the following as request URL:

	<Your CF application url>/updatepoll

The integration is now ready for use.

