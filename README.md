# Investigation for airgap deployments


<img align="right" alt="zarf logo" src=".images/Graviton Logo.png"  height="150" width="150" />

Edge graviton is thought to create specific configurations of enterprises for edge deployments. It utilizes many open source components as the basis, the case of Zarf. The idea is to have zero-touch provisioning of edge airgap environment.

The idea is to deploy out-of-the-box Kubernetes on airgap edge.

**Zero-touch -  SBOM - CIS and FIPS conformance**

Automatic provisioning with customized configurations with reliable out-of-the-box and well-tested solutions.

Are you ready for the computational shift to the edge: “Around 10% of enterprise-generated data is created and processed outside a traditional centralized data center or cloud. By 2025, Gartner predicts this figure will reach 75%”




## Why Use Zarf

- 💸 **Free and Open-Source.** Zarf will always be free to use and maintained by the open-source community.
- ⭐️ **Zero Dependencies.** As a statically compiled binary, the Zarf CLI has zero dependencies to run on any machine.
- 🔓 **No Vendor Lock.** There is no proprietary software that locks you into using Zarf. If you want to remove it, you can still use your Helm charts to deploy your software manually.
- 💻 **OS Agnostic.** Zarf supports numerous operating systems. A full matrix of supported OSes, architectures, and feature sets is coming soon.
- 📦 **Highly Distributable.** Integrate and deploy software from multiple secure development environments, including edge, embedded systems, secure cloud, data centers, and even local environments.
- 🚀 **Develop Connected, Deploy Disconnected.** Teams can build and configure individual applications or entire DevSecOps environments while connected to the internet. Once created, they can be packaged and shipped to a disconnected environment to be deployed.
- 💿 **Single File Deployments.** Zarf allows you to package the parts of the internet your app needs into a single compressed file to be installed without connectivity.
- ♻️ **Declarative Deployments.** Zarf packages define the precise state for your application, enabling it to be deployed the same way every time.
- 🦖 **Inherit Legacy Code.** Zarf packages can wrap legacy code and projects - allowing them to be deployed to modern DevSecOps environments.

## 📦 Out of the Box Features

- Automate Kubernetes deployments in disconnected environments
- Automate [Software Bill of Materials (SBOM)](https://docs.zarf.dev/docs/create-a-zarf-package/package-sboms) generation
- Build and [publish packages as OCI image artifacts](https://docs.zarf.dev/docs/zarf-tutorials/publish-and-deploy)
- Provide a [web dashboard](https://docs.zarf.dev/docs/deploy-a-zarf-package/view-sboms) for viewing SBOM output
- Create and verify package signatures with [cosign](https://github.com/sigstore/cosign)
- [Publish](https://docs.zarf.dev/docs/the-zarf-cli/cli-commands/zarf_package_publish), [pull](https://docs.zarf.dev/docs/the-zarf-cli/cli-commands/zarf_package_pull), and [deploy](https://docs.zarf.dev/docs/the-zarf-cli/cli-commands/zarf_package_deploy) packages from an [OCI registry](https://opencontainers.org/)
- Powerful component lifecycle [actions](https://docs.zarf.dev/docs/create-a-zarf-package/component-actions)
- Deploy a new cluster while fully disconnected with [K3s](https://k3s.io/) or into any existing cluster using a [kube config](https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/)
- Builtin logging stack with [Loki](https://grafana.com/oss/loki/)
- Built-in Git server with [Gitea](https://gitea.io/en-us/)
- Built-in Docker registry
- Builtin [K9s Dashboard](https://k9scli.io/) for managing a cluster from the terminal
- [Mutating Webhook](adr/0005-mutating-webhook.md) to automatically update Kubernetes pod's image path and pull secrets as well as [Flux Git Repository](https://fluxcd.io/docs/components/source/gitrepositories/) URLs and secret references
- Builtin [command to find images](https://docs.zarf.dev/docs/the-zarf-cli/cli-commands/zarf_dev_find-images) and resources from a Helm chart
- Tunneling capability to [connect to Kubernetes resources](https://docs.zarf.dev/docs/the-zarf-cli/cli-commands/zarf_connect) without network routing, DNS, TLS or Ingress configuration required

## 🛠️ Configurable Features

- Customizable [variables and package templates](https://docs.zarf.dev/examples/variables/) with defaults and user prompting
- [Composable packages](https://docs.zarf.dev/docs/create-a-zarf-package/zarf-components#composing-package-components) to include multiple sub-packages/components
- Component-level OS/architecture filtering

## Demo

[![preview](.images/zarf-v0.21-preview.gif)](https://www.youtube.com/watch?v=WnOYlFVVKDE)

_<https://www.youtube.com/watch?v=WnOYlFVVKDE>_

## ✅ Getting Started

To try Zarf out for yourself, visit the ["Try It Now"](https://zarf.dev/install) section on our website.

To learn more about Zarf and its use cases, visit [docs.zarf.dev](https://docs.zarf.dev/docs/zarf-overview). From the docs, you can learn more about:
- [installation](https://docs.zarf.dev/docs/getting-started/#installing-zarf)
- [using the CLI](https://docs.zarf.dev/docs/the-zarf-cli/),
- [making packages](https://docs.zarf.dev/docs/create-a-zarf-package/zarf-packages/),
- [Zarf package schema](https://docs.zarf.dev/docs/create-a-zarf-package/zarf-schema).

Using Zarf in GitHub workflows? Check out the [setup-zarf](https://github.com/defenseunicorns/setup-zarf) action. Install any version of Zarf and its `init` package with zero added dependencies.

## 🫶 Our Community

Join our community and developers on the [#Zarf slack](https://zarf.dev/slack) hosted on K8s slack. Our active community of developers, users, and contributors are available to answer questions, share examples, and find new ways use Zarf together!

We are so grateful to our Zarf community for contributing bug fixes and collaborating on new features:

<a href="https://github.com/defenseunicorns/zarf/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=defenseunicorns/zarf" />
</a>

Made with [contrib.rocks](https://contrib.rocks).

## 💻 Contributing

Want to contribute to Zarf?
Check out our [Contributor Guide](https://docs.zarf.dev/docs/contribute-to-zarf/contributor-guide) to learn more about how to set up your development environment and begin contributing.
We also recommend checking out our architectural diagram.

To dive deeper into the tech, you can read the [Nerd Notes](https://docs.zarf.dev/docs/contribute-to-zarf/nerd-notes) in our Docs.

![Architecture Diagram](./docs/.images/architecture.drawio.svg)

[Source DrawIO](docs/.images/architecture.drawio.svg)


## ⭐️ Special Thanks

> Early Zarf research and prototypes were developed jointly with [United States Naval Postgraduate School](https://nps.edu/) research you can read [here](https://calhoun.nps.edu/handle/10945/68688).

We would also like to thank the following awesome libraries and projects without which Zarf would not be possible!

[![pterm/pterm](https://img.shields.io/badge/pterm%2Fpterm-007d9c?logo=go&logoColor=white)](https://github.com/pterm/pterm)
[![mholt/archiver](https://img.shields.io/badge/mholt%2Farchiver-007d9c?logo=go&logoColor=white)](https://github.com/mholt/archiver)
[![spf13/cobra](https://img.shields.io/badge/spf13%2Fcobra-007d9c?logo=go&logoColor=white)](https://github.com/spf13/cobra)
[![go-git/go-git](https://img.shields.io/badge/go--git%2Fgo--git-007d9c?logo=go&logoColor=white)](https://github.com/go-git/go-git)
[![sigstore/cosign](https://img.shields.io/badge/sigstore%2Fcosign-2a1e71?logo=linuxfoundation&logoColor=white)](https://github.com/sigstore/cosign)
[![helm.sh/helm](https://img.shields.io/badge/helm.sh%2Fhelm-0f1689?logo=helm&logoColor=white)](https://github.com/helm/helm)
[![kubernetes](https://img.shields.io/badge/kubernetes-316ce6?logo=kubernetes&logoColor=white)](https://github.com/kubernetes)
