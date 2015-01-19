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

1) Ensure that [Go](http://golang.org/) version 1.2+ is installed on the system

2) Setup the GOPATH

    export GOPATH=~/go
    export PATH=$GOPATH/bin:$PATH
  
3) Download CF PLayground

    go get github.com/cloudfoundry-community/cfplayground
    cd $GOPATH/src/github.com/cloudfoundry-community/cfplayground
  
  *(Ignore any warnings about "no buildable Go source files")
  
4) Create a config file `config.json` under `config/` with the info of your Cloud Foundry environment, a sample config file is provided for reference `config/sameple_config.json`

_*_ if no `config.json` is found, `boshlite_config.json` will be used to target a local boshlite environment.
  
5) Run CF Playground

    go run main.go
    
* If you are running CF Playground under Linux, download Linux [CF CLI Binary](https://github.com/cloudfoundry/cli/releases), rename and replace the `pcf` file under `assets/cf/` with the downloaded binary.
    
Limitation
==========
- No Windows support (coming soon)
- Arbitrary app pushing (functioning, improvement to be made)
- Temp user account/space clean up (work in progress)
- Restore user session (functioning, improvement to be made)
- The supported CF commands are:
    - cf push
    - cf apps
    - cf app {app name}
    - cf delete {app name}





