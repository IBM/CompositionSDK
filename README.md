This repo is a fork of RedHat operator-sdk and will be re-used to develop composition-sdk.

## Build sdk
run `make all`. an executable named `operator-sdk` will be created under `build/` directory.

## Walkthrough example

#### Pre-requisites

1.  Set your namespace:
    ```
    export NAMESPACE=...
    ```

1.  deploy the following operators (instructions are available inside each of them):
    1.  [gnforchestrator](https://github.com/IBM/gnforchestrator)
    1.  [ping-operator](https://github.com/IBM/gnforchestrator/tree/main/demos/ping-pong/ping-operator)
    1.  [pong-operator](https://github.com/IBM/gnforchestrator/tree/main/demos/ping-pong/pong-operator)


## Demo steps

1.  Deploy the [sleep sample](https://github.com/istio/istio/tree/master/samples/sleep) from the [istio.io](http://istio.io).
    You will use this sample to send `curl` commands to the Multi-Language-Hello application.
    (No features of [istio.io](http://istio.io) are used in this demo, you will only use the sample from the project's
    [repository](https://github.com/istio/istio)).

    ```
    $ kubectl apply -f https://raw.githubusercontent.com/istio/istio/master/samples/sleep/sleep.yaml -n $NAMESPACE
    ```

1.  Export the sleep pod as an environment variable:

    ```
    $ export SLEEP_POD=$(kubectl get pod -l app=sleep -n $NAMESPACE -o jsonpath={.items..metadata.name})
    ```

1.  Generate a new composition operator using [composition-sdk](https://github.com/IBM/CompositionSDK).  
    The sdk expects to get as an input  `api-version` and `kind` of the generated CRD and a network service template:  
    <details>
    
    <summary>nsvc_template.yaml</summary>
    
    ```yaml
    apiVersion: gnforchestrator.ibm.com/v2alpha1
    kind: NetworkService
    metadata:
      name: example-pingpong
      labels:
        service: pingpong
    spec:
      properties:
        message: "Hello"
      components:
        ping:
          template:
            apiVersion: ping.example.com/v1alpha1
            kind: Ping
            metadata:
              name: "[% meta.name %]-ping"
              namespace: "[% meta.namespace %]"
            spec:
              pingVersion: v1.0
              pongAddress: "[% pong.status.ip %]"
              pongPort: 6006
        pong:
          template:
            apiVersion: pong.example.com/v1alpha1
            kind: Pong
            metadata:
              name: "[% meta.name %]-pong"
              namespace: "[% meta.namespace %]"
            spec:
              pongVersion: v1.4
              message: "[% message %]"
      statusTemplate:
        ip: "[% ping.status.ip %]"
        port: "[% ping.status.port %]"

    ```
    
    </details>
    
    ```
    $ operator-sdk composition test-operator \
    --api-version=pingpong.example.com/v1alpha1 \
    --kind=Pingpong \
     --generate-playbook \
     --nsvc-template=nsvc_template.yaml
    ```
    
    This command generates a skeletal test-operator application in the current directory.

1.  Implement your Pingpong CRD and CRs inside the generated operator under `deploy/crds` directory.  

1.  Install the new operator as explained in the auto generated `README.md` file.
    ```
    export REGISTRY=<YOUR REGISTRY>
    export IMAGE=$REGISTRY/$(basename $(pwd)):v0.0.1
    make docker-push "IMAGE=$IMAGE" "NAMESPACE=$NAMESPACE"
    make install "IMAGE=$IMAGE" "NAMESPACE=$NAMESPACE"
    ```

1.  Deploy Pingpong:

    <details>
        
    <summary>deploy/crds/pingpong.example.com_v1alpha1_pingpong_cr_message.yaml</summary>
    
    ```yaml
    apiVersion: pingpong.example.com/v1alpha1
    kind: Pingpong
    metadata:
      name: example-pingpong
    spec:
      message: HelloWorldTest
      pingVersion: Ping1
      pongVersion: Pong1   
    ```
    
    </details>
    
    ```
    $ kubectl apply -f deploy/crds/pingpong.example.com_v1alpha1_pingpong_cr_message.yaml -n $NAMESPACE
    ```  
    
1.  Watch the network services resources being created:
    ```
    $ watch kubectl get pingpong,nsvc,ping,pong,pod -n $NAMESPACE
    ```
    
1.  Create aliases to get the IP address and port of the pingpong resource:

    ```
    alias example-ip='kubectl get pingpong example-pingpong -n $NAMESPACE -o jsonpath={.status.ip}'
    alias example-port='kubectl get pingpong example-pingpong -n $NAMESPACE -o jsonpath={.status.port}'
    ```

1.  Run `curl` to ping hello endpoint:

    ```
    $ kubectl exec -it $SLEEP_POD -n $NAMESPACE -- curl $(example-ip):$(example-port)/hello
    Hello from Ping VNF
    ```

1.  Run `curl` to perform ping pong 3 times:

    ```
    $ kubectl exec -it $SLEEP_POD -n $NAMESPACE -- curl $(example-ip):$(example-port)/ping/3
    ping version v1.0: pong version v1.4 message HelloWorldTest.pong version v1.4 message HelloWorldTest.pong version v1.4 message HelloWorldTest.
    ```
    
1.  Clean Pingpong:
    ```
    $ kubectl delete -f deploy/crds/pingpong.example.com_v1alpha1_pingpong_cr_message.yaml -n $NAMESPACE
    ```  


<img src="website/static/operator_logo_sdk_color.svg" height="125px"></img>
[![Build Status](https://travis-ci.org/operator-framework/operator-sdk.svg?branch=master)](https://travis-ci.org/operator-framework/operator-sdk)
[![License](http://img.shields.io/:license-apache-blue.svg)](http://www.apache.org/licenses/LICENSE-2.0.html)
[![Go Report Card](https://goreportcard.com/badge/github.com/operator-framework/operator-sdk)](https://goreportcard.com/report/github.com/operator-framework/operator-sdk)

## Documentation

Docs can be found on the [Operator SDK website][sdk-docs].

## Overview

This project is a component of the [Operator Framework][of-home], an
open source toolkit to manage Kubernetes native applications, called
Operators, in an effective, automated, and scalable way. Read more in
the [introduction blog post][of-blog].

[Operators][operator-link] make it easy to manage complex stateful
applications on top of Kubernetes. However writing an operator today can
be difficult because of challenges such as using low level APIs, writing
boilerplate, and a lack of modularity which leads to duplication.

The Operator SDK is a framework that uses the
[controller-runtime][controller-runtime] library to make writing
operators easier by providing:

- High level APIs and abstractions to write the operational logic more intuitively
- Tools for scaffolding and code generation to bootstrap a new project fast
- Extensions to cover common operator use cases

## License

Operator SDK is under Apache 2.0 license. See the [LICENSE][license_file] file for details.

[controller-runtime]: https://github.com/kubernetes-sigs/controller-runtime
[license_file]:./LICENSE
[of-home]: https://github.com/operator-framework
[of-blog]: https://coreos.com/blog/introducing-operator-framework
[operator-link]: https://coreos.com/operators/
[sdk-docs]: https://sdk.operatorframework.io
