# ImmuDB .NET examples

These samples use the [official ImmuDB .NET client] (`immudb4net`).

[Official ImmuDB .NET client]: https://github.com/codenotary/immudb4net

Besides what exists here, you can find additional code snippets and examples in:

- [SDKs API](https://docs.immudb.io/master/develop/reading.html) documentation page
- [immudb4net Unit Tests](https://github.com/codenotary/immudb4net/tree/main/ImmuDB4Net.Tests)

## Contents

- [1. Description](#1-description)
- [2. Building locally from source](#2-building-locally-from-source)
- [3. Supported versions](#3-supported-versions)

## 1. Description

The code sample covers into a couple of functions the following functionality:

Function                           | Description
-----------------------------------|----------------------------------------------------------------------------------------------------------------------------
```OpenConnectionExample```        | How to create a new ```ImmuClient``` object using ```ImmuClientBuilder``` and open a connection to an ImmuDB server
```OpenConnectionExample```        | Execute ```VerifiedSet``` and ```VerifiedGet``` command
```SyncOpenConnectionExample```    | How to create a new ```ImmuClientSync``` object using ```ImmuClientBuilderSync``` and open a connection to an ImmuDB server
```AnotherOpenConnectionExample``` | How to create a new ```ImmuClient``` object using ```ImmuClientBuilder.Open``` function
```SetAllGetAllExample```          | How to use ```SetAll``` and ```GetAll``` functions
```GetSetScanUsageExample```       | How to use ```Set```, ```Get```, ```History```, ```Scan``` and ```ZScan```
```SqlUsageExample```              | How to use ```SQLExec``` and ```SQLQuery``` from the awaitable ```ImmuClient``` instance
```SyncSqlUsageExample```          | How to use ```SQLExec``` and ```SQLQuery``` from the non-awaitable ```ImmuClientSync``` instance

## 2. Building locally from source

Use ```dotnet build``` to build locally the ImmuDB client assembly. Optionally, one can run ```dotnet restore``` to ensure the ```ImmuDB4Net``` library is fetched from the [NuGet Repo] (`https://nuget.org/packages/ImmuDB4Net/#readme-body-tab`)

## 3. Supported versions

The sample project targets .NET 6
