# Kubernetes Management Levels

At Minekube, we understand that every user has different needs and preferences when it comes to managing their Kubernetes cluster. That's why we offer four different levels of management, ranging from fully managed to highly customizable.

No matter which level you choose, you'll always have admin access to your final Kubernetes cluster, and our goal is to provide a managed experience through the web UI. However, the more custom your setup, the less likely you'll be able to manage lower-level operations through the web UI.

[[TOC]]

Here are the four management levels we offer:

## Level 4: Fully Managed on Your Cloud

- **Description**: You only need to provide us with an API key from your preferred cloud provider, and we'll take care of the rest. We create machines, set up Kubernetes, and manage everything on your behalf. You don't even need an existing machine—all aspects are handled by us. Simply pay your cloud provider directly for the resources used.
- **Features**:
  - Fully managed Kubernetes cluster on your cloud provider
  - Machine creation, Kubernetes setup, and management handled by us
  - Automatic connection to the Minekube web UI to manage everything
  - No need to have an existing machine or run any commands
  - Pay your cloud provider directly for resources used

## Level 3: Bootstrapping a k3s Cluster by Command

- **Description**: At this level, we provide documentation explaining how to bootstrap a k3s cluster with a single command and automatically connect to the Minekube web UI. We control everything, from installing Kubernetes properly to node upgrades and deploying our controllers. This is the most managed level, offering simplicity and ease of use.
- **Features**:
  - Single command to set up a Kubernetes cluster using k3s
  - Automatic connection to the Minekube web UI to manage everything
  - Handled by us: installation, node upgrades, deployment of controllers, and more

## Level 2: Admin Access to Your Cluster

- **Description**: For users who want to choose how their Kubernetes cluster is set up. You maintain admin access to your Kubernetes API, allowing you to use our Minekube web UI seamlessly. While we handle the deployment and management of our controllers, as well as software updates and UI integration, you're responsible for node upgrades. We don't handle lower-level Kubernetes operations at this level.
- **Features**:
  - More control over your Kubernetes cluster
  - Convenience of using our Minekube web UI to manage everything except kubernetes nodes
  - Handled by us: deployment and management of controllers, software updates, and integration with our UI
  - Responsible for node upgrades

## Level 1: Raw Control with Minimum Management

- **Description**: This level is for users who desire the most raw control. You're responsible for ensuring our Minekube controllers are installed and kept updated on your cluster. We only handle core tasks necessary for you to manage game servers through our Minekube web UI. There's minimal management from our side—no controller updates, no node updates, just essential support for UI integration.
- **Features**:
  - Most control over your Kubernetes cluster
  - Convenience of using our Minekube web UI to manage game servers only
  - Responsible for installing and updating our controllers, node updates, and other lower-level Kubernetes operations
  - Handled by us: core tasks necessary to allow you to manage game servers through our Minekube web UI

## Getting Started

To get started, choose the management level that's right for you. We recommend starting with Level 1 if you're new to Kubernetes or want a simple and easy way to set up a cluster. If you want more control over your cluster, consider Level 2 or Level 3. And if you want a fully managed Kubernetes cluster on your cloud provider, choose Level 4.

Once you've chosen your management level, follow the instructions on the corresponding subdoc page to get started. If you have any questions or need help, don't hesitate to contact us.

