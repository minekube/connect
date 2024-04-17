# Kubernetes Management Levels

At Minekube, we understand that every user has different needs and preferences when it comes to managing their Kubernetes cluster. That's why we offer four different levels of management, ranging from fully managed to highly customizable.

No matter which level you choose, you'll always have admin access to your final Kubernetes cluster, and our goal is to provide a managed experience through the web UI. However, the more custom your setup, the less likely you'll be able to manage lower-level operations through the web UI.

[[TOC]]

## Comparison Table

Here's a detailed comparison table for the different Kubernetes management levels offered by Minekube:

| Feature/Level                   | Fully Managed on Your Cloud                 | Bootstrapping a k3s Cluster by Command         | Admin Access to Your Cluster                | Raw Control with Minimum Management         |
|---------------------------------|---------------------------------------------|------------------------------------------------|---------------------------------------------|---------------------------------------------|
| **Go to Guide**                 | [<VPBadge>Level 4 -></VPBadge>](level-4.md) | [<VPBadge>Level 3 -></VPBadge>](level-3.md)    | [<VPBadge>Level 2 -></VPBadge>](level-2.md) | [<VPBadge>Level 1 -></VPBadge>](level-1.md) |
| **Minekube Support**            | Full                                        | Full                                           | [Plus only](/plans)                         | [Plus only](/plans)                         |
| **Cluster Setup**               | Minekube-managed                            | User-managed, Minekube-assisted setup          | User-managed                                | User-managed                                |
| **Infrastructure Management**   | Minekube-managed                            | Minekube-managed                               | User-managed                                | User-managed                                |
| **Minekube Controllers**        | Minekube-managed                            | Minekube-managed                               | Minekube-managed                            | User-managed                                |
| **Web UI Integration**          | Automatic                                   | Automatic                                      | Automatic                                   | Manual setup                                |
| **Node Management**             | Minekube-managed                            | Minekube-managed                               | User-managed                                | User-managed                                |
| **Node Upgrades**               | Minekube-managed                            | Minekube-managed                               | User-managed                                | User-managed                                |
| **Cloud Provider Billing**      | Direct with provider                        | N/A (self-hosted k3s)                          | N/A (user-provided cluster)                 | N/A (user-provided cluster)                 |
| **Kubernetes API Admin Access** | Yes                                         | Yes                                            | Yes                                         | Yes                                         |
| **Customization Level**         | Low (fully managed)                         | Moderate                                       | High                                        | Very High                                   |
| **User Involvement**            | Minimal                                     | Moderate                                       | High                                        | Very High                                   |
| **Ease of Onboarding**          | Very Easy                                   | Easy                                           | Moderate                                    | Requires technical expertise                |
| **Best Suited For**             | Users seeking hands-off management          | Users comfortable with command-line operations | Experienced Kubernetes users                | Kubernetes experts                          |

> This table summarizes the key differences between each management level to help users decide which option best fits their needs. Whether you're looking for a fully managed service or prefer to have granular control over every aspect of your Kubernetes cluster, Minekube has a solution to support your game server hosting requirements.


## Level 4: Fully Managed on Your Cloud

[Go to Level 4 Guide ->](level-4.md)

- **Description**: You only need to provide us with an API key from your preferred cloud provider, and we'll take care of the rest. We create machines, set up Kubernetes, and manage everything on your behalf. You don't even need an existing machine—all aspects are handled by us. Simply pay your cloud provider directly for the resources used.
- **Features**:
  - Fully managed Kubernetes cluster on your cloud provider
  - Machine creation, Kubernetes setup, and management handled by us
  - Automatic connection to the Minekube web UI to manage everything
  - No need to have an existing machine or run any commands
  - Pay your cloud provider directly for resources used


## Level 3: Bootstrapping a k3s Cluster by Command

[Go to Level 3 Guide ->](level-3.md)

- **Description**: At this level, we provide documentation explaining how to bootstrap a k3s cluster with a single command and automatically connect to the Minekube web UI. We control everything, from installing Kubernetes properly to node upgrades and deploying our controllers. This is the most managed level, offering simplicity and ease of use.
- **Features**:
  - Single command to set up a Kubernetes cluster using k3s
  - Automatic connection to the Minekube web UI to manage everything
  - Handled by us: installation, node upgrades, deployment of controllers, and more

## Level 2: Admin Access to Your Cluster

[Go to Level 2 Guide ->](level-2.md)

- **Description**: For users who want to choose how their Kubernetes cluster is set up. You maintain admin access to your Kubernetes API, allowing you to use our Minekube web UI seamlessly. While we handle the deployment and management of our controllers, as well as software updates and UI integration, you're responsible for node upgrades. We don't handle lower-level Kubernetes operations at this level.
- **Features**:
  - More control over your Kubernetes cluster
  - Convenience of using our Minekube web UI to manage everything except kubernetes nodes
  - Handled by us: deployment and management of controllers, software updates, and integration with our UI
  - Responsible for node upgrades

## Level 1: Raw Control with Minimum Management

[Go to Level 1 Guide ->](level-1.md)

- **Description**: This level is for users who desire the most raw control. You're responsible for ensuring our Minekube controllers are installed and kept updated on your cluster. We only handle core tasks necessary for you to manage game servers through our Minekube web UI. There's minimal management from our side—no controller updates, no node updates, just essential support for UI integration.
- **Features**:
  - Most control over your Kubernetes cluster
  - Convenience of using our Minekube web UI to manage game servers only
  - Responsible for installing and updating our controllers, node updates, and other lower-level Kubernetes operations
  - Handled by us: core tasks necessary to allow you to manage game servers through our Minekube web UI

## Getting Started

To get started, choose the management level that's right for you. We recommend starting with Level 4 if you're new to Kubernetes or want a simple and easy way to set up a cluster. If you want more control over your cluster, consider Level 2 or Level 3.

Once you've chosen your management level, follow the instructions on the corresponding level guide page to get started. If you have any questions or need help, don't hesitate to contact us.
