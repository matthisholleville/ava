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

## Motivations Behind Ava

The motivation for creating Ava comes from my personal experience during on-call rotations. I wanted to make these moments—especially the ones that happen at night—easier for the SRE teams by automating as many repetitive tasks as possible using AI and a knowledge base (runbook). 

When alerts arrive in our team chat, we often find ourselves repeating the same actions:
- Checking monitoring dashboards
- Deleting or modifying resources
- Calling a URL to verify if a service is operational
- etc..

The idea behind Ava is to handle these repetitive steps automatically. My goal is that PagerDuty or others contacts the OnCall team only after these routine tasks are completed, allowing the team to focus on solving the core complexity of the issue. If Ava can mitigate or even fix the problem on its own, that’s a huge bonus.

To achieve this, I’ve built Ava using an REAct AI model powered by OpenAI Assistant. Ava understands alerts, investigates the issues, and takes actions (or provides mitigations) based on the runbooks I’ve uploaded into OpenAI’s vector store.

## Features

- Speeds up resolution by following your runbooks.
- Automates fixing one or more alerts using your runbooks and executors (functions).
- Works with Alert Manager webhooks.
- Compatible with OpenAI.
- Add an interactive mode to let Ava recheck if the issue is fixed after some time.
- More features coming soon... check the roadmap below.

[Watch the Ava demo on YouTube](https://youtu.be/VDAJqaBEv-s)

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

## Executors

Executors are functions Ava can use to act on your system. OpenAI does not perform actions directly; Ava executes them locally and sends the results back for better context. [Learn more about OpenAI assistant functions](https://platform.openai.com/docs/assistants/tools/function-calling).

### Built-in Executors

#### Kubernetes

- `deletePod`: Deletes a pod.
- `getPod`: Gets pod details.
- `listPod`: Lists pods in a namespace.
- `logsPod`: Shows the last 100 lines of pod logs.

#### Common

- `wait`: Waits before the next action.

#### Web

- `getUrl`: Makes a `GET` request to a URL and returns status and timing.

## Examples

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

## Roadmap

- Connect Ava to Slack (or other platforms) to update alert statuses in a channel.
- Create new executors for Kubernetes, databases (e.g., killing a PID), Prometheus, and Grafana (e.g., getting dashboard screenshots).
- Allow importing knowledge bases from other sources (e.g., Backstage, Notion).
- Support other AI backends (e.g., Llama, Gemini).

## Contributing

Coming soon...

## License

This project is licensed under the Apache-2.0 License. See the [LICENSE](LICENSE) file for details.

## Contact

For questions or feedback, open an issue or contact us at matthish29@gmail.com.
