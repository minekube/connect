# Connect API
<VPBadge>early stage</VPBadge>

_The Connect API is a powerful tool for developers
to integrate Connect into their projects and products.
Interact with players and [endpoints](/guide/#connect-endpoints) on the [Connect Network](/guide/#the-connect-network)
through the Connect API._

## Getting started

- [Java and Kotlin client](#java-and-kotlin-client)
- [Golang and other language clients](#golang-and-other-language-clients)

## Java and Kotlin client

You can access the Connect API through the official Java/Kotlin client
provided by the [Connect Plugin](https://github.com/minekube/connect-java).

Simply add the `connect-java:api` dependency to your project using Gradle or Maven with
[Jitpack](https://jitpack.io/#minekube/connect-java).

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


## Golang and other language clients

You can quickly get started with
[Buf Remote Packages](https://buf.build/minekube/connect/assets/main)
that provide client libraries for many programming languages for the Connect API.

For example, to install the Go Connect API client:

```shell
$ go get github.com/bufbuild/connect-go@latest
$ go get buf.build/gen/go/minekube/connect/bufbuild/connect-go@latest
$ go get buf.build/gen/go/minekube/connect/protocolbuffers/go@latest
```

[![Buf Remote Packages](/images/bufbuild-assets.png)](https://buf.build/minekube/connect/assets/main)
