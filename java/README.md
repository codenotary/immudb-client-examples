# immudb Java Examples

These samples use the [official immudb Java client] (`immudb4j`).

[Official immudb Java client]: https://github.com/codenotary/immudb4j

Besides what exists here, you can find additional code snippets and examples in:
- [SDKs API](https://docs.immudb.io/master/sdks-api.html) documentation page
- [immudb4j Unit Tests](https://github.com/codenotary/immudb4j/tree/master/src/test/java/io/codenotary/immudb4j)

<br/>

Please note that you need to have access to an `immudb` server in order to be able to use these examples.<br/>
But that's very easy: just take a look at [immudb Quickstart](https://docs.immudb.io/master/getstarted/quickstart.html) to know how to download and run it locally.

`immudb4j` is published on both [Maven Central](https://search.maven.org/artifact/io.codenotary/immudb4j) and GitHub Packages. 
Further details on installing it as a dependency can be found in [immudb4j](https://github.com/codenotary/immudb4j) repo itself. 
The classic Maven setup is already popular, therefore the section below describes the new GitHub Packages alternative, if interested.

<br/>

### Configuring Maven to use GitHub Packages

using `immudb4j` package is published at [GitHub Packages] and it requires authentication to download dependencies.

[GitHub Packages]: https://docs.github.com/en/packages

Please refer to GitHub documentation for a detailed explanation about [Authenticating with a personal access token]. 
Basically, you will need to do two things in order to download a Maven dependency hosted in `GitHub Packages`: authenticate and add the GitHub repository into your Maven settings.

[Authenticating with a personal access token]: https://docs.github.com/en/packages/using-github-packages-with-your-projects-ecosystem/configuring-apache-maven-for-use-with-github-packages

1. Create GitHub personal access token: https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token

2. Configure Maven to use GitHub Packages

Following info needs to be included into your `~/.m2/settings.xml` file. You will need to place your GitHub username and personal token in the `github` server.

```xml
<settings xmlns="http://maven.apache.org/SETTINGS/1.0.0"
  xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
  xsi:schemaLocation="http://maven.apache.org/SETTINGS/1.0.0
                      http://maven.apache.org/xsd/settings-1.0.0.xsd">

  <activeProfiles>
    <activeProfile>github</activeProfile>
  </activeProfiles>

  <profiles>
    <profile>
      <id>github</id>
      <repositories>
        <repository>
          <id>central</id>
          <url>https://repo1.maven.org/maven2</url>
          <releases><enabled>true</enabled></releases>
          <snapshots><enabled>true</enabled></snapshots>
        </repository>
        <repository>
          <id>github</id>
          <name>GitHub codenotary/immudb4j Apache Maven Packages</name>
          <url>https://maven.pkg.github.com/io/codenotary/immudb4j</url>
        </repository>
      </repositories>
    </profile>
  </profiles>

  <servers>
    <server>
      <id>github</id>
      <username>USERNAME</username>
      <password>TOKEN</password>
    </server>
  </servers>
</settings>
```
