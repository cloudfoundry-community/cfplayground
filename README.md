CF Playground
==============
The goal of this project is to provide an easily accessible environment for users who want to experience Cloud Foundry, without having to setup and to learn the operation of the platform.  CF Playground provides an interactive tutorial.

<p/>
![CF Playground in Action](https://github.com/cloudfoundry-community/cfplayground/raw/master/images/cfplayground.jpg)
<p/>

Setting up CF Playground
========================
The following instruction is for OSX/Linux, windows support is coming soon.
You will need to host your own Cloud Foundry environment (bosh-lite or any full deployment)

1) Ensure that [Go](http://golang.org/) version 1.2.x-64bit has been installed on the system

2) Setup the GOPATH

    export GOPATH=~/go
    export PATH=$GOPATH/bin:$PATH
  
3) Download CF PLayground

    go get github.com/cloudfoundry-community/cfplayground
    cd $GOPATH/src/github.com/cloudfoundry-community/cfplayground
  
  *(Ignore any warnings about "no buildable Go source files")
  
4) Create a config file under `config/` with the name `config.json` with the info to your Cloud Foundry environment, a sample config file is provided for reference `config/sameple_config.json`
  
5) Run CF Playground

    go run main.go
    
* If you are running CF Playground under Linux, download Linux [CF CLI Binary](https://github.com/cloudfoundry/cli/releases), rename and replace the `pcf` file under `assets/cf/` with the downloaded binary.
    
Limitation
==========
- No Windows support (coming soon)
- Arbitrary app pushing (work in progress)
- Temp user account/space clean up (work in progress)
- Restore user session (work in progress)
- The supported CF commands are:
    - cf push {app name} [-i NUM_INSTANCES] [-k DISK] [-m MEMORY] [-p PATH]
    - cf apps
    - cf app {app name}
    - cf delete {app name}
    - cf buildpacks
    - cf help [COMMAND]
    - cf scale {app name} [-i NUM_INSTANCES] [-k DISK] [-m MEMORY]
    - cf logs {app name} [--recent]
    - cf stop {app name}
    - cf start {app name}
    - cf restart {app name}
    - cf restage {app name}
    - cf services
    - cf domains
    - cf bind-service {service name} {app name}
    - cf unbind-service {service name} {app name}
    - cf delete-service {service name}
    - cf create-user-provided-service {service name} [-p CREDENTIALS]
    - cf marketplace
    - cf files {app name} [PATH]
    - cf env {app name}
    - cf map-route {app name} {domain} [-n HOSTNAME]
    - cf unmap-route {app name} {domain} [-n HOSTNAME]





