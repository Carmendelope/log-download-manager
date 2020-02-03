# Log-Download Manager
​
The Log-Download manager is responsible for operations related to downloading the log entries.

## Getting Started

### Prerequisites
​
Before installing this component, we need to have the following deployed:​

* [`application-manager`](https://github.com/nalej/application-manager): Log-Download manager accesses `application-manager` to send and receive operations and data from the Unified Logging component.
​
### Build and compile
​
In order to build and compile this repository use the provided Makefile:
​
```
make all
```
​
This operation generates the binaries for this repo, downloads the required dependencies, runs existing tests and generates ready-to-deploy Kubernetes files.
​
### Run tests
​
Tests are executed using Ginkgo. To run all the available tests:
​
```
make test
```

### Update dependencies
​
Dependencies are managed using Godep. For an automatic dependencies download use:
​
```
make dep
```
​
In order to have all dependencies up-to-date run:
​
```
dep ensure -update -v
```

​
## Contributing
​
Please read [contributing.md](contributing.md) for details on our code of conduct, and the process for submitting pull requests to us.
​
​
## Versioning
​
We use [SemVer](http://semver.org/) for versioning. For the available versions, see the [tags on this repository](https://github.com/nalej/log-download-manager/tags). 
​
## Authors
​
See also the list of [contributors](https://github.com/nalej/log-download-manager/contributors) who participated in this project.
​
## License
This project is licensed under the Apache 2.0 License - see the [LICENSE-2.0.txt](LICENSE-2.0.txt) file for details.



