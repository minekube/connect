# Hosting Options for Minekube Games

Minekube Games is designed with versatility in mind to ensure that you can host your game servers in a way that suits your needs and expertise. Whether you prefer to manage your own infrastructure, rely on containerization, or opt for a managed provider, Minekube Games offers the flexibility to accommodate your choice.

[[TOC]]

## Compare hosting solutions

Consider factors like ease of use, existing infrastructure, and developer needs:

| Feature            | <VPBadge>[Minekube Providers ->](provider.md)</VPBadge> | <VPBadge>[Kubernetes ->](kubernetes/)</VPBadge> | <VPBadge>[Docker ->](container.md)</VPBadge> |
|--------------------|---------------------------------------------------------|-------------------------------------------------|----------------------------------------------|
| **Setup**          | Easiest                                                 | Easy                                            | Moderate                                     |
| **Infrastructure** | Managed                                                 | Anywhere                                        | Anywhere                                     |
| **User Interface** | Minekube Web UI                                         | Minekube Web UI                                 | No                                           |
| **Updates**        | Automatic                                               | Automatic                                       | Manual                                       |
| **Pricing**        | Pay-as-you-go / Subscription / Free                     | Flexible / Free                                 | Flexible / Free                              |

**Recommendations:**

- For players and low-friction users who want simplicity and ease of use, Minekube Provider's managed solution is recommended, as it offers a user-friendly web UI and handles updates automatically.
- If you already have one or more machines available or an existing Kubernetes cluster, Minekube offers various Kubernetes management levels of control and customization.
- If you're a developer or want to quickly set up a game server without the need for Minekube's web UI and automated updates, Docker can be a lightweight option for running Game containers.

## Self-hosted Anywhere

For those who like to keep things under their own control, Minekube Games can run on any platform that supports containers or Kubernetes:

- **Kubernetes (e.g. k3s)** <VPBadge>Recommended</VPBadge>: If you're looking for orchestration and scaling capabilities, Kubernetes is the way to go. Whether it's a lightweight k3s setup or a full-fledged Kubernetes cluster, Minekube Games thrives in these environments.
- **Container / Docker** <VPBadge type='danger'>No Minekube Web UI</VPBadge>: Ideal for straightforward deployment and local development, running Minekube Games within Docker containers is a popular choice for its simplicity and portability.

## Managed Provider

- **Minekube Games Providers**: For ease and convenience, choose one of the Minekube Games Providers. You won't need to worry about the underlying infrastructure and can have your game server up and running in no time, with all the benefits of a managed service.

Discover the benefits and learn how to get started with each hosting option in the following sections.