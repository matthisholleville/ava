<div align="center">
<br />
<p align="center">
  <a href="https://github.com/matthisholleville/ava">
    <img src="docs/images/ava-ai.webp" alt="Logo" width="80" height="80">
  </a>

<h1 align="center">A.V.A</h1>
</p>
</div>

![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/matthisholleville/ava)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/matthisholleville/ava)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go version](https://img.shields.io/github/go-mod/go-version/matthisholleville/ava.svg)](https://github.com/matthisholleville/ava)


`Ava` is an AI assistant designed to help you during on-call rotations.

Give it your runbooks, use existing executors (functions), or create your own functions. `Ava` will use these to handle your alerts, following your instructions. The goal is to understand problems and try to resolve them based on your runbooks.

**This system is experimental and should not be used in production without fully understanding the risks.**

https://github.com/user-attachments/assets/0caf1e90-f8d0-4885-96c8-0db33173dfe4

## Motivations Behind `Ava`

The motivation for creating `Ava` comes from my personal experience during on-call rotations. I wanted to make these moments—especially the ones that happen at night—easier for the SRE teams by automating as many repetitive tasks as possible using AI and a knowledge base (runbook). 

When alerts arrive in our team chat, we often find ourselves repeating the same actions:
- Checking monitoring dashboards
- Deleting or modifying resources
- Calling a URL to verify if a service is operational
- etc..

The idea behind `Ava` is to handle these repetitive steps automatically. My goal is that PagerDuty or others contacts the OnCall team only after these routine tasks are completed, allowing the team to focus on solving the core complexity of the issue. If `Ava` can mitigate or even fix the problem on its own, that’s a huge bonus.

To achieve this, I’ve built `Ava` using an REAct AI model powered by OpenAI Assistant. `Ava` understands alerts, investigates the issues, and takes actions (or provides mitigations) based on the runbooks I’ve uploaded into OpenAI’s vector store.

## Features

- Speeds up resolution by following your runbooks and provides useful assistance to resolve the problem.
- **Optionally** Automates fixing one or more alerts using your runbooks and executors (functions).
- Works with Alert Manager webhooks.
- REST API.
- Compatible with OpenAI.
- React with Slack event
- Allow importing knowledge bases from local path & Github.
- More features coming soon... check the roadmap below.

## Demo 

- [Watch the Ava CLI demo on YouTube](https://youtu.be/VDAJqaBEv-s)
- [Watch the Ava Slack Bot demo on Youtube](https://youtu.be/ya9ySnGCTak)

## Installation

1. Clone the repository:
    ```bash
    git clone https://github.com/matthisholleville/ava.git
    ```
2. Go to the project directory:
    ```bash
    cd ava
    ```
3. Install dependencies:
    ```bash
    make build
    ```

## Configuration

Ava uses a configuration file to set up its usage. By default, this file is located at `$XDG_CONFIG_HOME/ava.yaml`. You can change the configuration folder using the `--config` flag.

| OS | Path |
| ------- | ------------------------------------------------ |
| MacOS | ~/Library/Application Support/ava/ava.yaml |
| Linux | ~/.config/ava/ava.yaml |
| Windows | %LOCALAPPDATA%/ava/ava.yaml |

This file is essential, and Ava cannot start without it. It should not contain sensitive data. For such cases, you can use environment variables.

When starting up, Ava parses the configuration file. If it detects a value matching the pattern `${.*}`, it will extract the content and check if an environment variable exists with that name. If not, the program will fail and return an error log.


## Usage

- Get an OpenAI API Key from [OpenAI](https://openai.com).
    - Rename `.envrc-sample` to `envrc` and add your OpenAI API Key.
    - Run `direnv allow`.
- Upload your runbooks using:
    ```bash
    ava knowledge add
    ```
- Start chatting:
    ```bash
    ./ava chat -m "Pod $CHANGE_ME in namespace default Crashlooping."
    ```

## Knowledge

The runbook knowledge base is one of the most important parts of Ava. It allows Ava to understand problems, guide you to mitigate or fix them, and even fix them automatically using executors. 

Your runbooks are uploaded into an OpenAI vector store and utilized by Ava during chats with the model. [Learn more here](https://platform.openai.com/docs/assistants/tools/file-search).

You can currently import this knowledge base from:
- `local`
- `git` (Tested with private GitHub repositories using an authentication token)

More sources will be added soon (see the roadmap section below).

### Examples

<details>

<summary>Using Github private repository & auth token</summary>

```bash
ava knowledge add -s git -r https://github.com/MyPrivateOrg/my-private-repository.git -t "ghp_dflkjcIO..."
```

</details>

## Analysis Mode

In its default version, Ava greatly accelerates the resolution phase by providing assistance for a problem based on the runbooks you’ve uploaded into its knowledge base. This ensures that, regardless of the user handling the issue, they don’t need to locate the right runbook, connect to the correct tool, etc. Additionally, if the problem isn’t covered by a runbook, Ava can leverage the AI model’s knowledge base to guide the user.

## Automatic Fix Mode

**Warning: This mode is experimental and takes actions (see the Executors section below) on the environment where Ava is deployed.**

This optional mode allows Ava to attempt to fix the problem automatically using its list of available Executors. This mode enables rapid mitigation of issues, reducing stress and mental load for the operator in charge of fixing the problem. It allows the operator to focus on problem-solving while minimizing the chances of human error and avoiding repetitive tasks.

### How to enable automatic fix mode

<details>

<summary>CLI</summary>

Use the flag `--enable-executors=true`

Example :

```bash
go run main.go chat -m "Pod web-server-5b866987d8-4nhsg in namespace default Crashlooping." --enable-executors=true
```

</details>

<details>

<summary>API</summary>

You are free to enable or disable executors in the body of the requests.

Example : 

```bash
curl -X POST https://your-url/chat \
-H "Content-Type: application/json" \
-d '{
  "backend": "openai",
  "enableExecutors": true,
  "language": "en",
  "message": "Pod web-server-5b866987d8-sxmtj in namespace default Crashlooping."
}'
```

----------

For reacting to AlertManager webhooks, you simply need to specify the environment variable `ENABLE_EXECUTORS_ON_WEBHOOK`.

Example :

```bash
export ENABLE_EXECUTORS_ON_WEBHOOK="true"
```

</details>

### Executors

Executors are functions Ava can use to act on your system. OpenAI does not perform actions directly; Ava executes them locally and sends the results back for better context. [Learn more about OpenAI assistant functions](https://platform.openai.com/docs/assistants/tools/function-calling).

#### Built-in Executors

##### Kubernetes


##### Read-only

- `describeService`: Describe details of a service.
- `getCronJobs`:  List all CronJobs in a namespace.
- `getDeployment`: Retrieve details of a deployment.
- `getHPA`: Retrieve the status of Horizontal Pod Autoscalers in a namespace.
- `getNode`: Get the details of a node.
- `getPod`: Gets pod details.
- `listNamespaces`: List all namespaces in the cluster.
- `listPods`: Lists pods in a namespace.
- `podLogs`: Shows the last 100 lines of pod logs.
- `topPods`: Retrieve CPU and memory usage of all pods in a namespace.

##### Write

**By default, these executors can impact environments and are not enabled. To activate them, simply enable them in the configuration.**

- `deletePod`: Delete a pod.
- `rolloutDeployment`: Perform a rollout restart for a deployment.

##### Common

- `wait`: Waits before the next action.

##### Web

- `getUrl`: Makes a `GET` request to a URL and returns status and timing.

## Serve Mode

Ava provides an REST API. The CLI mode and SERVER mode offer the same features, with one key difference: the API mode requires a PostgreSQL database to function.

The API stores chat results, threads, and other information in the database, which can be particularly useful for auditing purposes.

### Running the Server Mode

<details>

<summary>Steps to Launch Server Mode</summary>

**Before starting the server mode, you need a PostgreSQL database up and running.** You can use Docker or a cloud solution like CloudSQL.

1. Rename `.envrc-sample` to `.envrc`.
2. Update the environment variable `export DATABASE_URL="CHANGE_ME"` with a valid connection string.
3. Export the variables: `direnv allow`.
4. Initialize the PostgreSQL schema: 
    ```bash
    go run github.com/steebchen/prisma-client-go db push
    ```
5. Generate the Swagger documentation: 
    ```bash
    make swagger
    ```
6. Start the server mode: 
    ```bash
    go run main.go serve
    ```

Once launched, you can access the Swagger documentation at the following URL: [http://localhost:8080/swagger/index.html#](http://localhost:8080/swagger/index.html#).

</details>


## Roadmap

- Create new executors for Kubernetes, databases (e.g., killing a PID), Prometheus, and Grafana (e.g., getting dashboard screenshots).
- Allow importing knowledge bases from other sources (e.g., Backstage, Notion).
- Support other AI backends (e.g., Llama, Gemini).
- Automatic PostMortem Generation.

## Examples

<details>

<summary>Connecting Slack</summary>

### Installation

Before connecting Slack, you must deploy Ava. To connect Ava to Slack, follow the steps below:

1. Create a [new Slack App](https://api.slack.com/apps) named `Ava Bot`. **The name is crucial and must not be changed!**
2. Copy the `Verification Token` from the `Basic Information` page and the `Bot User OAuth Token` from the `OAuth & Permissions` page of your Slack application, then paste them into Ava's configuration file (e.g., `ava config edit` or `./charts/ava/values.yaml` if you uses Kubernetes):
    ```yaml
    # ava config
    events:
        type: slack
        slack:
            validationToken: ${SLACK_VALIDATION_TOKEN}
            botToken: ${SLACK_BOT_TOKEN}
    ```
    - **If you use Kubernetes dont forget to create the K8s secret for Slack** `kubectl create secret generic slack-secret --from-literal=validation-token=$(echo $SLACK_VALIDATION_TOKEN) --from-literal=bot-token=$(echo $SLACK_BOT_TOKEN)`
3. Start Ava in server mode with a publicly accessible URL so Slack can send events.
4. On the `Event Subscriptions` page, enable events and add your URL in the `Request URL` section (`$MY_URL/event/slack`).
5. Further down on the same page, open the `Subscribe to bot events` section and add the permissions `message:channels` and `message:groups`.
6. On the `OAuth & Permissions` page, add the scopes `chat:write`, `users.profile:read`, and `users:read`.
7. Install your app into the desired channel(s).

You can test the setup by sending a direct message to your bot: `@Ava Bot Hello how are you today?`. You should see Ava react to the event in its logs and send a Slack message containing `👀`.

### Interacting with Ava

To start interacting with Ava, you can:
- Begin your message by mentioning it: `@Ava Bot my pod example is crashlooping. Do we have any runbook to understand & fix the problem?`
- Configure AlertManager to send a Slack message. Ava will respond not only when mentioned but also to messages from other bots.

</details>

<details>

<summary>API Mode with AlertManager Webhook</summary>

This section shows how to set up a local environment to demonstrate Ava with AlertManager webhooks. It installs:
- AlertManager
- Prometheus
- Ava (server mode)
- Example webserver-chaos

**1. Install Prometheus**

Run these commands:

```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm install prometheus prometheus-community/prometheus --namespace monitoring --values ./docs/examples/custom-values.yaml --create-namespace
helm install prometheus-operator-crds prometheus-community/prometheus-operator-crds --namespace monitoring
```

**2. Create a secret with your OpenAI API Key**

Run this command:

```bash
kubectl create secret generic ai-backend-secret --from-literal=openai-api-key=$(echo $OPENAI_API_KEY) -n monitoring
```

**3. Deploy Ava using Helm**

Run this command:

```bash
cd charts/ava/
helm install ava . -n monitoring
```

**4. Deploy demo web-server**

Run this command:

```bash
kubectl apply -f ./docs/examples/crashloop.yaml -n default
```

The demo webserver handles two routes:
1. `/`: Returns `Hello from Ava :)`.
2. `/chaos`: Runs `sys.exit(1)`.

---------

When everything is installed, run this command:

```bash
curl http://$(kubectl get svc web-server-service -o jsonpath='{.status.loadBalancer.ingress[*].ip}' -n default)/chaos
```

This will trigger chaos in the webserver. Ava should detect the issue and fix it. After a few minutes, your pod should be healthy.

</details>

## Contributing

Please read our [contributing guide](./CONTRIBUTING.md).

## License

This project is licensed under the Apache-2.0 License. See the [LICENSE](LICENSE) file for details.

## Contact

For questions or feedback, open an issue or contact us at matthish29@gmail.com.
