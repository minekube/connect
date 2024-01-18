# Connect API Clients

_Connect provides code generated API clients!
They save developers time and effort,
while ensuring consistency and maintainability.
With type safety, IDE support, and language compatibility,
these clients streamline the integration process,
abstract away complexities, and empower developers to
focus on building robust and scalable applications._

[[TOC]]

## Provided by Connect Plugin

As a plugin developer you can depend on [Connect Plugin](/guide/connectors/plugin) in your plugin's dependencies.
Make sure that your `plugin.yml` has a `depend: [ connect ]` to ensure that the Connect Plugin is loaded before your
plugin.

The Connect Plugin provides authenticated stubs to the Connect API through the `ConnectApi` global instance.

```java [Main.java <VPBadge>Java</VPBadge>]
com.minekube.connect.api.ConnectApi.getInstance().getClients()...
```

Simply add the `connect-java:api` dependency to your project using Gradle or Maven with the
[Jitpack](https://jitpack.io/#minekube/connect-java) repository.

::: code-group

```kotlin [build.gradle.kts <VPBadge>Gradle Kotlin</VPBadge>]
repositories {
  maven("https://jitpack.io")
}

dependencies {
  api("com.github.minekube.connect-java:api:latest")
}
```

```groovy [build.gradle <VPBadge>Gradle Groovy</VPBadge>]
repositories {
  maven { url 'https://jitpack.io' }
}

dependencies {
    api 'com.github.minekube.connect-java:api:latest'
}
```

```xml [pom.xml <VPBadge>Maven</VPBadge>]

<repositories>
    <repository>
        <id>jitpack.io</id>
        <url>https://jitpack.io</url>
    </repository>
</repositories>

<dependencies>
<dependency>
    <groupId>com.github.minekube.connect-java</groupId>
    <artifactId>api</artifactId>
    <version>latest</version>
    <scope>provided</scope>
</dependency>
</dependencies>
```

:::

Checkout [Code Examples](/guide/api/examples) to see it in action.

## Other languages

You can also use
[Buf Remote Packages](https://buf.build/minekube/connect/assets/main)
that provide client libraries for many programming languages for the Connect API.
Make sure to include the required request header to [Authenticate](/guide/api/authentication) self-built clients with
the Connect API.

Supported languages:

- Java/Kotlin
- Golang
- JavaScript/TypeScript
- [read more...](https://buf.build/docs/bsr/remote-packages/overview/)

For example, to add the following modules to your Go project:

```shell
go get github.com/bufbuild/connect-go@latest
go get buf.build/gen/go/minekube/connect/bufbuild/connect-go@latest
go get buf.build/gen/go/minekube/connect/protocolbuffers/go@latest
```

Checkout [Code Examples](/guide/api/examples) for code examples in different languages.

---

[![Buf Remote Packages](/images/bufbuild-assets.png)](https://buf.build/minekube/connect/assets/main)
