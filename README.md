# PROMOTER - keep the application(s) running

_We will discuss about [`promoter`](https://github.com/pnkj-kmr/promoter/releases/) **WHAT**, **WHY** & **HOW** ?_

## WHAT?

An application which helps to keep the other application(s) running in server(s). It means that promoter helps to maintain the application state active(running) if application stopped in one of server then promoter make sure, application will resume it's state in other server.

Let's have an example - _"application APP_1 is deployed in Delhi and it's backup in Bengaluru APP_1 and we want if Delhi application goes down due to some reason so that operation should be resumed by Bengaluru application by starting backup application"_, here `promoter` fits well which start the application at Bengaluru and when the Delhi APP_1 comes up than promoter stops the Bengaluru APP_1 application. All operations are perfomed automatically.

## WHY?

Here we have 2-scenarios where promoter can help us.

![snapshot](./resource/p1.svg)

promoter comes with multitenancy along with

- High availability (HA)
- Disaster recovery (DC-DR)
- Presistance - Load balanacing (LB)

## HOW?

`promoter` setup configuration is very easy, there are only to 2-steps to take care at one server, if we want to deploy over distributed architecture, just repeat the same over mutli-servers.

- Configuration `app.yml`
- Run as a service

And download the `promoter` [**here**](https://github.com/pnkj-kmr/promoter/releases/)

#### **Configuration app.yml**

```
# app.yml
# promoter configuration file

config:
  # unique id of promoter
  # my_id very important to work properly,
  # make sure every promoter has unique id
  my_id: 1

  # cluster helps to communicate with other promoter(s)
  # it's kind of a secret key for communication
  cluster_id: XXX

  # priority of promoter node
  # helps to choose leader node
  # if more than one leader_nodes added
  # higher priority node will be leader
  # range 0 - 100
  priority: 100

  # promoter app refresh interval
  # default refresh_rate - 30 seconds
  # refresh_rate: 30

  # timeout for any computation
  # default - 90 seconds
  # timeout: 90

  # promoter roles
  # there are two possible roles
  #   - leader : decisions maker node
  #   - broker : a worker node
  role:
    - leader
    - broker

  # promoter starting address
  bind: 127.0.0.1:8080

  # all leader nodes bind address
  # it's very important
  # 1) leader node selection
  # 2) broker node hearbeat management
  leader_nodes:
    - 127.0.0.1:8080
    # - 127.0.0.1:8081

applications:
  # name of application
  # add multiple if needed

  app1:
    # application unique id
    app_id: 1
    # application priority helps to make
    # high chances to run first
    # over multiple node applications
    # range 0 - 1000
    priority: 101
    # application persistance stands that
    # no of application(s)
    # keep active(running) at a time
    # value > 0
    persist: 1
    # application status matching text
    # to decide application up/down
    status_match: "redis             started"
    # applcation status command
    status: |
      echo "status checking...."
      brew services list
    # applcation stop command
    stop: |
      echo "stopping..."
      brew services stop redis
    # applcation start command
    start: |
      echo "starting..."
      brew services run redis
```

#### **Run as a service**

[linux]: vi /etc/systemd/system/promoter.service

```
[Unit]
Description=Promoter
Documentation=https://github.com/pnkj-kmr/promoter
Requires=network.target
After=network.target

[Service]
Type=simple
WorkingDirectory=/opt/promoter
ExecStart=promoter -c app
ExecStop=
Restart=on-abnormal

[Install]
WantedBy=multi-user.target
```

### PROMOTER [Deployment Layouts]

![snapshot](./resource/p2.svg)

:)
