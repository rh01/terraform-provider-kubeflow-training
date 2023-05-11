# terraform-provider-kubeflow-training
Terraform provider for Kubeflow Training Operator


## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) >= 0.12.x
-	[Go](https://golang.org/doc/install) >= 1.16

## Usage


TODO 


## Contributing

Contributions to this repository are very welcome! Found a bug or do you have a suggestion? Please open an issue. Do you know how to fix it? Pull requests are welcome as well! To get you started faster, a Makefile is provided.

Make sure to install [Terraform](https://learn.hashicorp.com/terraform/getting-started/install.html), [Go](https://golang.org/doc/install) (for automated testing) and Make (optional, if you want to use the Makefile) on your computer. Install [tflint](https://github.com/terraform-linters/tflint) to be able to run the linting.

* Format your code: `make fmt`
* Run tests: `make test`
* Run acceptance tests: `make testacc`. This creates resources on your Kubernetes cluster, use with caution. We use [k3s](https://k3s.io/) in the CICD pipelines, to start from a fresh environment each time.

## License

MIT license. Please see [LICENSE](LICENSE) for details.