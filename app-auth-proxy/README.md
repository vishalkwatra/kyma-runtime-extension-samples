## Overview

This sample provides a reverse proxy feature which dispatches requests to other microservices running in Kyma. It includes a middleware to handle authentication which is based on Open ID Connect and can be configured using XSUAA or SAP IAS. This example implementation is storing sessions using an in memory store which is meant for testing only. See [store-implementations](https://github.com/gorilla/sessions#store-implementations) for other options.

This sample demonstrates how to:

- Create a development Namespace in the Kyma runtime.
- Consume the SCP service XSUAA
- Deploy the following Kubernetes resources:
  - API deployment written in GO
  - API Rule
  - Service
  - Configmap
  - ServiceBinding
  - ServiceBindingUsage

## Prerequisites

- SAP Cloud Platform, Kyma runtime instance
- [Docker](https://www.docker.com/)
- [Go](https://golang.org/doc/install)
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) configured to use the `KUBECONFIG` file downloaded from the Kyma runtime

## Steps

### Create XSUAA Service Instance

1. Create a new `dev` Namespace:

```shell script
kubectl create namespace dev
```

2. Within the Kyma console open the namespace `dev`
3. Choose `Service Management` -> `Catalog`.
4. Choose the service `Authorization & Trust Management`
5. Choose `Add`
6. Choose the Plan `application`
7. Choose `Add parameters` and provide the object after adjusting it to your needs.

```json
{
  "oauth2-configuration": {
    "redirect-uris": [
      "https://app-auth-proxy.<cluster domain>/oauth/callback",
      "http://localhost:8000/oauth/callback"
    ]
  },
  "xsappname": "app-auth-proxy"
}
```

<sup> For a complete list of parameters visit [Application Security Descriptor Configuration Syntax](https://help.sap.com/viewer/4505d0bdaf4948449b7f7379d24d0f0d/2.0.04/en-US/6d3ed64092f748cbac691abc5fe52985.html) </sup>

1. Once the instance is provisioned choose the option `Create Credentials`
2. Under the `Credentials` tab choose the `Secret` which should display the instance secret in a dialog. Choose `Decode` to view the values. These will be needed if running the sample locally.

### Run the API locally

1. Set the environment variables required to connect with the XSUAA instance which can be found in the `Secret` generated with the service instance:

```shell script
export IDP_clientid='<instance clientid>'
export IDP_clientsecret=<instance clientsecret>
export IDP_url=<instance url>
```

2. Adjust the config.json which contains the following properties. The provided config.json is configured to use the examples
   - [React frontend MS SQL](../frontend-react-mssql/README.md)
     - Requires the configmap API_URL to point to `https://app-auth-proxy.<cluster domain>`
   - [Golang MS SQL database API](../api-mssql-go/README.md)

| Property                   | Description                                                                         | Remarks                                                             |
| -------------------------- | ----------------------------------------------------------------------------------- | ------------------------------------------------------------------- |
| routes                     | An array of routes to be proxied                                                    |                                                                     |
| routes.path                | The incoming path                                                                   |                                                                     |
| routes.priority            | The priority of the path with 1 be the highest                                      |                                                                     |
| routes.protected           | If the auth middleware is required on the path                                      |                                                                     |
| routes.remove_from_path    | If assigned, this value will be removed from routes.path before the call is proxied |                                                                     |
| routes.target              | The target of the proxied route which can be a service url                          |                                                                     |
| token_endpoint_auth_method | The htttp method used to during authentication                                      | For XSUAA use client_secret_post, for SAPIAS us client_secret_basic |
| redirect_uri               | The registered redirect_uri to be called                                            |                                                                     |
| debug                      | Toggle debug on or off                                                              |                                                                     |
| cookie.session_name        | The name of the session cookie                                                      |                                                                     |
| cookie.max_age_seconds     | The max age of the session cookie                                                   |                                                                     |
| cookie.key                 | The key used to encrypt the session cookie                                          |                                                                     |
| cookie.httponly            | If the cookie can be accessed with Javascript or only http                          |                                                                     |
| cookie.secure              | If the cookie should only be used with https                                        |                                                                     |
| cookie.samesite            | If the cookie shouldn't be sent with cross-origin requests                          |                                                                     |

1. Run the application:

```shell script
go run ./cmd/proxy
```

5. Accessible endpoints include
   - http://localhost:8000/
   -

### Build the Docker image

1. Build and push the image to your Docker repository:

```shell script
docker build -t {your-docker-account}/app-auth-proxy -f docker/Dockerfile .
docker push {your-docker-account}/app-auth-proxy
```

1. To run the image locally adjust the config.json and either set the env variables individually, or copy them from your environment:

```shell script
  docker run -p 8000:8000 --env-file ./env.list --mount type=bind,source=$(pwd)/config/config.json,target=/app/config/config.json -d jcawley5/app-auth-proxy:latest
  OR
  docker run -p 8000:8000 --env-file <(env | grep IDP) --mount type=bind,source=$(pwd)/config/config.json,target=/app/config/config.json -d jcawley5/app-auth-proxy:latest
```

### Deploy the APP

1. Create a new `dev` Namespace:

```shell script
kubectl create namespace dev
```

2. Within `./k8s/configmap.yaml` adjust the values and then apply the ConfigMap:

```shell script
kubectl -n dev apply -f ./k8s/configmap.yaml
```

3. Get the name of the ServiceInstance:

```shell script
kubectl -n dev get serviceinstances
```

For example:

| NAME                   | CLASS                     | PLAN        | STATUS | AGE |
| ---------------------- | ------------------------- | ----------- | ------ | --- |
| **_xsuaa-showy-yard_** | ClusterServiceClass/xsuaa | application | Ready  | 63m |

4. Within `./k8s/deployment.yaml` adjust the value of `<Service Instance Name>` to the XSUAA service instance name and the apply the Deployment:

```shell script
kubectl -n dev apply -f ./k8s/deployment.yaml
```

5. Apply the APIRule:

```shell script
kubectl -n dev apply -f ./k8s/apirule.yaml
```

6. Verify that the Deployment is up and running:

```shell script
kubectl -n dev get deployment app-auth-proxy
```

7. Use the APIRule:

- `https://app-auth-proxy.{cluster-domain}`
